package service

import (
	"bufio"
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

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/lucas-clemente/quic-go"
	uuid "github.com/satori/go.uuid"
)

var UDPconn *net.UDPConn
var PeopleMap map[string]quic.Stream

func Dial(addr string, token string) error {
	quicConfig := &quic.Config{
		ConnectionIDLength: 4,
		KeepAlive:          true,
	}
	tlsConf := &tls.Config{
		InsecureSkipVerify:     true,
		NextProtos:             []string{"bench"},
		SessionTicketsDisabled: true,
	}
	session, err := quic.DialAddr(addr, tlsConf, quicConfig)
	defer session.CloseWithError(0, "")
	if err != nil {
		return err
	}
	// stream, err := session.OpenStreamSync(context.Background())
	// if err != nil {
	// 	return err
	// }

	return nil
}

func SayHello(stream quic.Stream, to string) {
	msg := model.MessageModel{}
	msg.Type = "hello"
	msg.Data = "hello"
	msg.To = to
	msg.From = config.ServerInfo.Token
	msg.UUId = uuid.NewV4().String()
	b, _ := json.Marshal(msg)
	prefixLength := file.PrefixLength(len(b))

	data := append(prefixLength, b...)
	stream.Write(data)
}

var pathsss string

//文件分片发送
func SendFileData(stream quic.Stream, filePath, to, uuid string) error {

	fStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	blockSize, length := file.GetBlockInfo(fStat.Size())

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("读取失败", err)
		return err
	}
	bufferedReader := bufio.NewReader(f)
	buf := make([]byte, blockSize)
	for i := 0; i < length; i++ {

		tran := model.TranFileModel{}

		_, err = bufferedReader.Read(buf)

		if err == io.EOF {
			fmt.Println("读取完毕", err)
		}

		tran.Hash = file.GetHashByContent(buf)
		tran.Index = i

		msg := model.MessageModel{}
		msg.Type = "file_data"
		msg.Data = tran
		msg.From = config.ServerInfo.Token
		msg.To = to
		msg.UUId = uuid
		b, _ := json.Marshal(msg)
		stream.Write(b)
	}
	defer stream.Close()
	return nil
}

//发送数据
func SendData(stream quic.Stream, m model.MessageModel) {
	b, _ := json.Marshal(m)
	stream.Write(b)
}

//读取数据
func ReadContent(stream quic.Stream) (model.MessageModel, error) {
	path := ""
	for {
		prefixByte := make([]byte, 4)
		c1, err := io.ReadFull(stream, prefixByte)
		fmt.Println(c1, err)
		prefixLength, err := strconv.Atoi(string(prefixByte))

		messageByte := make([]byte, prefixLength)
		t, err := io.ReadFull(stream, messageByte)
		fmt.Println(t, err)
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
				file.SpliceFiles("", path, dataModel.Length)
				break
			}
		} else {
			return m, nil
		}
	}
	return model.MessageModel{}, nil
}
