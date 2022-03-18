package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/quic_helper"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	httper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/lucas-clemente/quic-go"
	"gorm.io/gorm"
)

type PersonService interface {
	GetPersionInfo(token string) (m model.PersionModel, err error)
}

type personService struct {
	db *gorm.DB
}

var IpInfo model.PersionModel

func PushIpInfo(token string) {

	m := model.PersionModel{}
	m.Ips = GetDeviceAllIP()
	m.Token = token
	b, _ := json.Marshal(m)

	if reflect.DeepEqual(IpInfo, m) {
		return
	}
	head := make(map[string]string)
	infoS := httper2.Post(config.ServerInfo.Handshake+"/v1/update", b, "application/json", head)
	fmt.Println(infoS)
}
func (p *personService) GetPersionInfo(token string) (m model.PersionModel, err error) {
	infoS := httper2.Get(config.ServerInfo.Handshake+"/v1/ips/"+token, nil)
	err = json.Unmarshal([]byte(infoS), &m)
	return
}

func NewPersonService(db *gorm.DB) PersonService {
	return &personService{db: db}
}

//=======================================================================================================================================================================

var StreamList map[string]quic.Stream
var ServiceMessage chan model.MessageModel

func UDPService() {

	srcAddr := &net.UDPAddr{
		IP: net.IPv4zero, Port: 9904}
	var err error
	UDPConn, err = net.ListenUDP("udp", srcAddr)
	if err != nil {
		fmt.Println(err)
	}
	listener, err := quic.Listen(UDPConn, quic_helper.GetGenerateTLSConfig(config.ServerInfo.Token), quic_helper.GetQUICConfig())
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()
	ctx := context.Background()
	acceptFailures := 0
	const maxAcceptFailures = 10
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		default:
		}

		session, err := listener.Accept(ctx)
		if err != nil {
			fmt.Println("Listen (BEP/quic): Accepting connection:", err)

			acceptFailures++
			if acceptFailures > maxAcceptFailures {
				// Return to restart the listener, because something
				// seems permanently damaged.
				fmt.Println(err)
				return
			}

			// Slightly increased delay for each failure.
			time.Sleep(time.Duration(acceptFailures) * time.Second)

			continue
		}

		acceptFailures = 0

		streamCtx, cancel := context.WithTimeout(ctx, time.Second*10)
		stream, err := session.AcceptStream(streamCtx)
		cancel()
		if err != nil {
			fmt.Println("failed to accept stream from %s: %v", session.RemoteAddr(), err)
			_ = session.CloseWithError(1, err.Error())
			continue
		}

		// prefixByte := make([]byte, 4)
		// c1, err := io.ReadFull(stream, prefixByte)
		// fmt.Println(c1, err)
		// prefixLength, err := strconv.Atoi(string(prefixByte))
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// messageByte := make([]byte, prefixLength)
		// t, err := io.ReadFull(stream, messageByte)
		// fmt.Println(t, err)
		// m := model.MessageModel{}
		// err = json.Unmarshal(messageByte, &m)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		go ProcessingContent(stream)
	}
}

//处理内容
func ProcessingContent(stream quic.Stream) {
	for {
		prefixByte := make([]byte, 6)
		_, err := io.ReadFull(stream, prefixByte)
		if err != nil {
			return
		}
		prefixLength, err := strconv.Atoi(string(prefixByte))
		if err != nil {
			fmt.Println(err)
		}
		messageByte := make([]byte, prefixLength)
		_, err = io.ReadFull(stream, messageByte)
		if err != nil {
			return
		}
		m := model.MessageModel{}
		err = json.Unmarshal(messageByte, &m)
		if err != nil {
			fmt.Println(err)
		}
		if m.Type == types.PERSONHELLO {
			//nothing
			continue
		} else if m.Type == types.PERSONDIRECTORY {
			var list []model.Path
			if m.Data.(string) == "" || m.Data.(string) == "/" {
				for _, v := range config.FileSettingInfo.ShareDir {
					tempList := MyService.ZiMa().GetDirPath(v)
					list = append(list, tempList...)
				}
			} else {
				list = MyService.ZiMa().GetDirPath(m.Data.(string))
			}
			m.To = m.From
			m.Data = list
			m.From = config.ServerInfo.Token
			SendData(stream, m)
			break
		} else if m.Type == types.PERSONDOWNLOAD {

			SendFileData(stream, m.Data.(string), m.From, m.UUId)
			break
		} else if m.Type == types.PERSONADDFRIEND {
			friend := model2.FriendModel{}
			dataModelByte, _ := json.Marshal(m.Data)
			err := json.Unmarshal(dataModelByte, &friend)
			if err != nil {
				fmt.Println(err)
				continue
			}
			go MyService.Friend().UpdateOrCreate(friend)
			mi := model2.FriendModel{}
			mi.Avatar = config.UserInfo.Avatar
			mi.Profile = config.UserInfo.Description
			mi.Name = config.UserInfo.NickName
			m.To = m.From
			m.Data = mi
			m.Type = types.PERSONADDFRIEND
			m.From = config.ServerInfo.Token

			SendData(stream, m)
			break
		} else if m.Type == types.PERSONCONNECTION {
			if len(m.Data.(string)) > 0 {
				UDPAddressMap[m.From] = m.Data.(string)
			} else {
				delete(UDPAddressMap, m.From)
			}
			mi := model2.FriendModel{}
			mi.Avatar = config.UserInfo.Avatar
			mi.Profile = config.UserInfo.Description
			mi.Name = config.UserInfo.NickName
			mi.Token = config.ServerInfo.Token
			msg := model.MessageModel{}
			msg.Type = types.PERSONADDFRIEND
			msg.Data = mi
			msg.To = m.From
			msg.From = config.ServerInfo.Token
			msg.UUId = m.UUId
			Dial(msg, false)

			break
		} else {
			//不应有不做返回的数据
			//ServiceMessage <- m
			break
		}
	}
	stream.Close()

}

//文件分片发送
func SendFileData(stream quic.Stream, filePath, to, uuid string) error {
	summary := model.FileSummaryModel{}

	msg := model.MessageModel{}
	msg.Type = types.PERSONSUMMARY
	msg.From = config.ServerInfo.Token
	msg.To = to
	msg.UUId = uuid

	fStat, err := os.Stat(filePath)
	if err != nil {

		summary.Message = err.Error()

		msg.Data = summary

		summaryByte, _ := json.Marshal(msg)
		summaryPrefixLength := file.PrefixLength(len(summaryByte))
		summaryData := append(summaryPrefixLength, summaryByte...)
		stream.Write(summaryData)
		return err
	}

	blockSize, length := file.GetBlockInfo(fStat.Size())

	f, err := os.Open(filePath)
	if err != nil {

		summary.Message = err.Error()
		msg.Data = summary

		summaryByte, _ := json.Marshal(msg)
		summaryPrefixLength := file.PrefixLength(len(summaryByte))
		summaryData := append(summaryPrefixLength, summaryByte...)
		stream.Write(summaryData)
		return err
	}

	//send file summary first
	summary.BlockSize = blockSize
	summary.Hash = file.GetHashByPath(filePath)
	summary.Length = length
	summary.Name = fStat.Name()
	summary.Size = fStat.Size()

	msg.Data = summary

	summaryByte, _ := json.Marshal(msg)
	summaryPrefixLength := file.PrefixLength(len(summaryByte))
	summaryData := append(summaryPrefixLength, summaryByte...)
	stream.Write(summaryData)

	bufferedReader := bufio.NewReader(f)
	buf := make([]byte, blockSize)
	for i := 0; i < length; i++ {

		tran := model.TranFileModel{}

		n, err := bufferedReader.Read(buf)

		if err == io.EOF {
			fmt.Println("读取完毕", err)
		}

		tran.Hash = file.GetHashByContent(buf[:n])
		tran.Index = i
		tran.Length = length

		fileMsg := model.MessageModel{}
		fileMsg.Type = types.PERSONDOWNLOAD
		fileMsg.Data = tran
		fileMsg.From = config.ServerInfo.Token
		fileMsg.To = to
		fileMsg.UUId = uuid
		b, _ := json.Marshal(fileMsg)
		prefixLength := file.PrefixLength(len(b))
		dataLength := file.DataLength(len(buf[:n]))
		data := append(append(append(prefixLength, b...), dataLength...), buf[:n]...)
		stream.Write(data)
	}
	defer stream.Close()
	return nil
}
