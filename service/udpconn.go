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
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/lucas-clemente/quic-go"
	uuid "github.com/satori/go.uuid"
)

var UDPconn *net.UDPConn
var PeopleMap map[string]quic.Stream
var Message chan model.MessageModel

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

	session, err := quic.DialAddrContext(ctx, addr, tlsConf, quicConfig)
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

//读取数据
func ReadContent(stream quic.Stream) {
	path := ""
	for {
		prefixByte := make([]byte, 4)
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

			filepath := path + strconv.Itoa(dataModel.Index)

			err = ioutil.WriteFile(filepath, dataByte, 0644)
			if dataModel.Index >= (dataModel.Length - 1) {
				//file.SpliceFiles("", path, dataModel.Length)
				break
			}
		} else {
			Message <- m
		}
	}
	Message <- model.MessageModel{}
}
