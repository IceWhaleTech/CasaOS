package service

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
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

func Dial(addr string, msg model.MessageModel) (m model.MessageModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Message = make(chan model.MessageModel)
	quicConfig := &quic.Config{
		ConnectionIDLength: 4,
		KeepAlive:          true,
	}
	tlsConf := &tls.Config{
		InsecureSkipVerify:     true,
		NextProtos:             []string{"bench"},
		SessionTicketsDisabled: true,
	}
	srcAddr := &net.UDPAddr{
		IP: net.IPv4zero, Port: 9904} //注意端口必须固定
	//addr
	if len(addr) == 0 {
		addr = config.ServerInfo.Handshake + ":9527"
	}
	dstAddr, err := net.ResolveUDPAddr("udp", addr)

	//DialTCP在网络协议net上连接本地地址laddr和远端地址raddr。net必须是"udp"、"udp4"、"udp6"；如果laddr不是nil，将使用它作为本地地址，否则自动选择一个本地地址。
	//(conn)UDPConn代表一个UDP网络连接，实现了Conn和PacketConn接口

	session, err := quic.DialContext(ctx, UDPConn, dstAddr, srcAddr.String(), tlsConf, quicConfig)
	if err != nil {
		return m, err
	}

	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		session.CloseWithError(1, err.Error())
		return m, err
	}

	SayHello(stream, msg.To)

	SendData(stream, msg)

	go ReadContent(stream)
	result := <-Message
	stream.Close()
	return result, nil
}

func SayHello(stream quic.Stream, to string) {
	msg := model.MessageModel{}
	msg.Type = "hello"
	msg.Data = "hello"
	msg.To = to
	msg.From = config.ServerInfo.Token
	msg.UUId = uuid.NewV4().String()
	SendData(stream, msg)
}

var pathsss string

//发送数据
func SendData(stream quic.Stream, m model.MessageModel) {
	b, _ := json.Marshal(m)
	prefixLength := file.PrefixLength(len(b))
	data := append(prefixLength, b...)
	stream.Write(data)
}

var Summary map[string]model.FileSummaryModel

//读取数据
func ReadContent(stream quic.Stream) {
	for {
		prefixByte := make([]byte, 6)
		c1, err := io.ReadFull(stream, prefixByte)
		fmt.Println(c1, err, string(prefixByte))
		prefixLength, err := strconv.Atoi(string(prefixByte))

		messageByte := make([]byte, prefixLength)
		t, err := io.ReadFull(stream, messageByte)
		fmt.Println(t, err, string(messageByte))
		m := model.MessageModel{}
		err = json.Unmarshal(messageByte, &m)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(m)
		//传输数据需要继续读取
		if m.Type == "file_data" {
			dataModelByte, _ := json.Marshal(m.Data)
			dataModel := model.TranFileModel{}
			err := json.Unmarshal(dataModelByte, &dataModel)
			fmt.Println(err)

			dataLengthByte := make([]byte, 8)
			t, err = io.ReadFull(stream, dataLengthByte)
			dataLength, err := strconv.Atoi(string(dataLengthByte))
			if err != nil {
				fmt.Println(err)
			}
			dataByte := make([]byte, dataLength)
			t, err = io.ReadFull(stream, dataByte)
			if err != nil {
				fmt.Println(err)
			}
			sum := md5.Sum(dataByte)
			hash := hex.EncodeToString(sum[:])
			if dataModel.Hash != hash {
				fmt.Println("hash不匹配", hash, dataModel.Hash)
			}

			tempPath := config.AppInfo.RootPath + "/temp" + "/" + m.UUId
			file.IsNotExistMkDir(tempPath)
			filepath := tempPath + "/" + strconv.Itoa(dataModel.Index)
			tempFile, err := os.Stat(filepath)

			if os.IsNotExist(err) || tempFile.Size() == 0 {
				err = ioutil.WriteFile(filepath, dataByte, 0644)
			} else {
				if file.GetHashByPath(filepath) != dataModel.Hash {
					os.Remove(filepath)
					err = ioutil.WriteFile(filepath, dataByte, 0644)
				}
			}

			files, err := ioutil.ReadDir(tempPath)

			if len(files) >= dataModel.Length {
				summary := Summary[m.UUId]
				file.SpliceFiles(tempPath, config.FileSettingInfo.DownloadDir+"/"+summary.Name, dataModel.Length, 0)
				if file.GetHashByPath(config.FileSettingInfo.DownloadDir+"/"+summary.Name) == summary.Hash {
					file.RMDir(tempPath)
					task := model2.PersionDownloadDBModel{}
					task.UUID = m.UUId
					task.State = types.DOWNLOADFINISH
					MyService.Download().EditDownloadState(task)
				} else {
					os.Remove(config.FileSettingInfo.DownloadDir + "/" + summary.Name)
					task := model2.PersionDownloadDBModel{}
					task.UUID = m.UUId
					task.State = types.DOWNLOADERROR
					MyService.Download().EditDownloadState(task)
				}

				break
			}
		} else if m.Type == "summary" {

			dataModel := model.FileSummaryModel{}
			if m.UUId == m.UUId {
				dataModelByte, _ := json.Marshal(m.Data)
				err := json.Unmarshal(dataModelByte, &dataModel)
				fmt.Println(err)
			}

			task := model2.PersionDownloadDBModel{}
			task.UUID = m.UUId
			task.Name = dataModel.Name
			task.Length = dataModel.Length
			task.Size = dataModel.Size
			task.State = types.DOWNLOADING
			task.BlockSize = dataModel.BlockSize
			task.Hash = dataModel.Hash
			task.Type = 0
			MyService.Download().SaveDownload(task)

			Summary[m.UUId] = dataModel

		} else if m.Type == "connection" {
			UDPAddressMap[m.From] = m.Data.(string)
			fmt.Println("udpconn", m)
			mi := model2.FriendModel{}
			mi.Avatar = config.UserInfo.Avatar
			mi.Profile = config.UserInfo.Description
			mi.Name = config.UserInfo.UserName
			mi.Token = config.ServerInfo.Token
			msg := model.MessageModel{}
			msg.Type = types.PERSONADDFRIEND
			msg.Data = mi
			msg.To = m.From
			msg.From = config.ServerInfo.Token
			msg.UUId = m.UUId
			Dial(m.Data.(string), msg)
			Message <- m
			break
		} else {
			Message <- m
		}
	}
	Message <- model.MessageModel{}
}
