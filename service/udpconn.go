package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	path2 "path"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/quic_helper"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/lucas-clemente/quic-go"
	uuid "github.com/satori/go.uuid"
)

var UDPConn *net.UDPConn
var PeopleMap map[string]quic.Stream
var Message chan model.MessageModel
var UDPAddressMap map[string]string

func Dial(msg model.MessageModel, server bool) (m model.MessageModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Message = make(chan model.MessageModel)
	_, port, err := net.SplitHostPort(UDPConn.LocalAddr().String())
	if config.ServerInfo.UDPPort != port {
		config.ServerInfo.UDPPort = port
		config.Cfg.Section("server").Key("UDPPort").SetValue(port)
		config.Cfg.SaveTo(config.SystemConfigInfo.ConfigPath)
	}
	p, err := strconv.Atoi(port)
	srcAddr := &net.UDPAddr{
		IP: net.IPv4zero, Port: p} //注意端口必须固定
	addr := UDPAddressMap[msg.To]
	ticket := msg.To
	if server {
		addr = config.ServerInfo.Handshake + ":9527"
		ticket = "bench"
	}
	dstAddr, err := net.ResolveUDPAddr("udp", addr)

	//DialTCP在网络协议net上连接本地地址laddr和远端地址raddr。net必须是"udp"、"udp4"、"udp6"；如果laddr不是nil，将使用它作为本地地址，否则自动选择一个本地地址。
	//(conn)UDPConn代表一个UDP网络连接，实现了Conn和PacketConn接口

	session, err := quic.DialContext(ctx, UDPConn, dstAddr, srcAddr.String(), quic_helper.GetClientTlsConfig(ticket), quic_helper.GetQUICConfig())
	if err != nil {
		if msg.Type == types.PERSONDOWNLOAD {
			task := MyService.Download().GetDownloadById(msg.UUId)
			task.Error = err.Error()
			task.State = types.DOWNLOADERROR
			MyService.Download().SetDownloadError(task)
		}
		if config.SystemConfigInfo.Analyse != "False" {
			go MyService.Casa().PushConnectionStatus(msg.UUId, err.Error(), msg.From, msg.To, msg.Type)
		}

		return m, err
	}

	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		if msg.Type == types.PERSONDOWNLOAD {
			task := MyService.Download().GetDownloadById(msg.UUId)
			task.Error = err.Error()
			task.State = types.DOWNLOADERROR
			MyService.Download().SetDownloadError(task)
		}
		if config.SystemConfigInfo.Analyse != "False" {
			go MyService.Casa().PushConnectionStatus(msg.UUId, err.Error(), msg.From, msg.To, msg.Type)
		}
		session.CloseWithError(1, err.Error())
		return m, err
	}

	SayHello(stream, msg.To)

	SendData(stream, msg)

	go ReadContent(stream)
	result := <-Message
	stream.Close()
	if config.SystemConfigInfo.Analyse != "False" {
		go MyService.Casa().PushConnectionStatus(msg.UUId, "OK", msg.From, msg.To, msg.Type)
	}
	return result, nil
}

func SayHello(stream quic.Stream, to string) {
	msg := model.MessageModel{}
	msg.Type = types.PERSONHELLO
	msg.Data = "hello"
	msg.To = to
	msg.From = config.ServerInfo.Token
	msg.UUId = uuid.NewV4().String()
	SendData(stream, msg)
}

//发送数据
func SendData(stream quic.Stream, m model.MessageModel) {
	b, _ := json.Marshal(m)
	prefixLength := file.PrefixLength(len(b))
	data := append(prefixLength, b...)
	stream.Write(data)
}

//读取数据
func ReadContent(stream quic.Stream) {
	for {
		prefixByte := make([]byte, 6)
		_, err := io.ReadFull(stream, prefixByte)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 1)
			for k, v := range CancelList {
				tempPath := config.AppInfo.RootPath + "/temp" + "/" + v
				fmt.Println(file.RMDir(tempPath))
				delete(CancelList, k)
			}
			break
		}
		prefixLength, err := strconv.Atoi(string(prefixByte))
		if err != nil {
			fmt.Println(err)
			break
		}
		messageByte := make([]byte, prefixLength)
		_, err = io.ReadFull(stream, messageByte)
		if err != nil {
			fmt.Println(err)
			break
		}
		m := model.MessageModel{}
		err = json.Unmarshal(messageByte, &m)
		fmt.Println("客户端", m)
		if err != nil {
			fmt.Println(err)
			break
		}

		if m.Type == types.PERSONDOWNLOAD {

			dataModelByte, _ := json.Marshal(m.Data)
			dataModel := model.TranFileModel{}
			err := json.Unmarshal(dataModelByte, &dataModel)
			if err != nil {
				fmt.Println(err)
				continue
			}

			dataLengthByte := make([]byte, 8)
			_, err = io.ReadFull(stream, dataLengthByte)
			if err != nil {
				fmt.Println(err)
				continue
			}
			dataLength, err := strconv.Atoi(string(dataLengthByte))
			if err != nil {
				fmt.Println(err)
				continue
			}
			dataByte := make([]byte, dataLength)
			_, err = io.ReadFull(stream, dataByte)
			if err != nil {
				fmt.Println(err)
				continue
			}
			sum := md5.Sum(dataByte)
			hash := hex.EncodeToString(sum[:])
			if dataModel.Hash != hash {
				fmt.Println("hash不匹配", hash, dataModel.Hash)
				continue
			}
			tempPath := config.AppInfo.RootPath + "/temp" + "/" + m.UUId
			file.IsNotExistMkDir(tempPath)
			filepath := tempPath + "/" + strconv.Itoa(dataModel.Index)
			_, err = os.Stat(filepath)

			if os.IsNotExist(err) {
				err = ioutil.WriteFile(filepath, dataByte, 0644)
				task := model2.PersonDownloadDBModel{}
				task.UUID = m.UUId
				if err != nil {
					task.Error = err.Error()
					task.State = types.DOWNLOADERROR
					MyService.Download().SetDownloadError(task)
				}

			} else {
				if file.GetHashByPath(filepath) != dataModel.Hash {
					os.Remove(filepath)
					err = ioutil.WriteFile(filepath, dataByte, 0644)
					task := model2.PersonDownloadDBModel{}
					task.UUID = m.UUId
					if err != nil {
						task.Error = err.Error()
						task.State = types.DOWNLOADERROR
						MyService.Download().SetDownloadError(task)
					}
				}
			}

			files, err := ioutil.ReadDir(tempPath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(files) >= dataModel.Length {
				summary := MyService.Download().GetDownloadById(m.UUId)
				summary.State = types.DOWNLOADFINISH
				MyService.Download().EditDownloadState(summary)
				fullPath := file.GetNoDuplicateFileName(path2.Join(summary.LocalPath, summary.Name))
				file.SpliceFiles(tempPath, fullPath, dataModel.Length, 0)
				if file.GetHashByPath(fullPath) == summary.Hash {
					file.RMDir(tempPath)
					summary.State = types.DOWNLOADFINISHED
					MyService.Download().EditDownloadState(summary)
				} else {
					os.Remove(config.FileSettingInfo.DownloadDir + "/" + summary.Name)

					summary.State = types.DOWNLOADERROR
					summary.Error = "hash mismatch"
					MyService.Download().SetDownloadError(summary)
				}

				break
			}
		} else if m.Type == types.PERSONSUMMARY {

			dataModel := model.FileSummaryModel{}
			dataModelByte, _ := json.Marshal(m.Data)
			err := json.Unmarshal(dataModelByte, &dataModel)
			if err != nil {
				fmt.Println(err)
			}

			task := MyService.Download().GetDownloadById(m.UUId)
			fullPath := path2.Join(task.LocalPath, task.Name)
			task.State = types.DOWNLOADING
			if len(dataModel.Message) > 0 {
				task.State = types.DOWNLOADERROR
				task.Error = dataModel.Message
			}
			if file.Exists(fullPath) && file.GetHashByPath(fullPath) == dataModel.Hash {
				task.State = types.DOWNLOADFINISHED
				go func(from, uuid string) {
					m := model.MessageModel{}
					m.Data = ""
					m.From = config.ServerInfo.Token
					m.To = from
					m.Type = types.PERSONCANCEL
					m.UUId = uuid
					CancelList[uuid] = uuid
					Dial(m, false)
				}(task.From, task.UUID)

			}
			task.UUID = m.UUId
			task.Name = dataModel.Name
			task.Length = dataModel.Length
			task.Size = dataModel.Size
			task.BlockSize = dataModel.BlockSize
			task.Hash = dataModel.Hash
			task.Type = 0
			task.From = m.From
			MyService.Download().SaveDownload(task)

		} else if m.Type == types.PERSONCONNECTION {
			if len(m.Data.(string)) > 0 {
				UDPAddressMap[m.From] = m.Data.(string)
			} else {
				delete(UDPAddressMap, m.From)
			}
			// mi := model2.FriendModel{}
			// mi.Avatar = config.UserInfo.Avatar
			// mi.Profile = config.UserInfo.Description
			// mi.NickName = config.UserInfo.NickName
			// mi.Token = config.ServerInfo.Token
			msg := model.MessageModel{}
			msg.Type = types.PERSONHELLO
			msg.Data = ""
			msg.To = m.From
			msg.From = config.ServerInfo.Token
			msg.UUId = m.UUId
			go Dial(msg, false)
			Message <- m
			break
		} else if m.Type == "get_ip" {
			notify := model2.AppNotify{}
			notify.CustomId = m.From
			if len(m.Data.(string)) == 0 {
				if _, ok := UDPAddressMap[m.From]; ok {
					notify.Type = types.NOTIFY_TYPE_PERSION_FIRNED_LEAVE
					go MyService.Notify().SendText(notify)
				}
				delete(UDPAddressMap, m.From)
				Message <- m
				break
			}
			if _, ok := UDPAddressMap[m.From]; !ok {
				notify.Type = types.NOTIFY_TYPE_PERSION_FIRNED_LIVE
				go MyService.Notify().SendText(notify)
			}
			UDPAddressMap[m.From] = m.Data.(string)
			Message <- m
			break
		} else {
			Message <- m
		}
	}
	Message <- model.MessageModel{}
}

func SendIPToServer() {
	msg := model.MessageModel{}
	msg.Type = types.PERSONHELLO
	msg.Data = ""
	msg.From = config.ServerInfo.Token
	msg.To = config.ServerInfo.Token
	msg.UUId = uuid.NewV4().String()

	Dial(msg, true)
}

func LoopFriend() {
	list := MyService.Friend().GetFriendList()
	for i := 0; i < len(list); i++ {
		msg := model.MessageModel{}
		msg.Type = "get_ip"
		msg.Data = ""
		msg.From = config.ServerInfo.Token
		msg.To = list[i].Token
		msg.UUId = uuid.NewV4().String()

		Dial(msg, true)

		msg.Type = types.PERSONHELLO
		msg.Data = ""
		msg.From = config.ServerInfo.Token
		msg.To = list[i].Token
		msg.UUId = uuid.NewV4().String()
		if _, ok := UDPAddressMap[list[i].Token]; ok {
			go Dial(msg, false)
		}
		go func(shareId string) {
			user := MyService.Casa().GetUserInfoByShareId(shareId)
			m := model2.FriendModel{}
			m.Token = shareId
			friend := MyService.Friend().GetFriendById(m)
			if friend.Version != user.Version {
				friend.Avatar = user.Avatar
				friend.NickName = user.NickName
				friend.Profile = user.Desc
				friend.Version = user.Version
				MyService.Friend().UpdateOrCreate(friend)
			}
		}(list[i].Token)

	}
}
