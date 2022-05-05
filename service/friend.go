package service

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"github.com/IceWhaleTech/CasaOS/pkg/quic_helper"
	model2 "github.com/IceWhaleTech/CasaOS/service/model"
	"github.com/IceWhaleTech/CasaOS/types"
	"github.com/lucas-clemente/quic-go"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type FriendService interface {
	AddFriend(m model2.FriendModel)
	DeleteFriend(m model2.FriendModel)
	EditFriendMark(m model2.FriendModel)
	EditFriendWrite(m model2.FriendModel)
	EditFriendBlock(m model2.FriendModel)
	GetFriendById(m model2.FriendModel) model2.FriendModel
	GetFriendList() (list []model2.FriendModel)
	GetFriendListRemote() (list []model2.FriendModel)
	UpdateAddFriendType(m model2.FriendModel)
	AgreeFrined(id string)
	GetFriendByToken(token string) model2.FriendModel
	UpdateOrCreate(m model2.FriendModel)
	InternalInspection(ips []string, token string)
}

type friendService struct {
	db *gorm.DB
}

func (p *friendService) AgreeFrined(id string) {
	var m model2.FriendModel
	p.db.Model(&m).Where("token = ?", id).Update("state", types.FRIENDSTATEDEFAULT)
}
func (p *friendService) AddFriend(m model2.FriendModel) {
	p.db.Create(&m)
}
func (p *friendService) DeleteFriend(m model2.FriendModel) {
	p.db.Where("token = ?", m.Token).Delete(&m)
}
func (p *friendService) EditFriendMark(m model2.FriendModel) {
	p.db.Model(&m).Where("token = ?", m.Token).Update("mark", m.Mark)
}
func (p *friendService) EditFriendWrite(m model2.FriendModel) {
	p.db.Model(&m).Where("token = ?", m.Token).Update("write", m.Write)
}
func (p *friendService) EditFriendBlock(m model2.FriendModel) {
	p.db.Model(&m).Where("token = ?", m.Token).Update("block", m.Block)
}
func (p *friendService) GetFriendById(m model2.FriendModel) model2.FriendModel {
	p.db.Model(m).Where("token = ?", m.Token).First(&m)
	return m
}

func (p *friendService) GetFriendList() (list []model2.FriendModel) {
	p.db.Select("nick_name", "avatar", "profile", "token", "state", "mark", "block", "version").Find(&list)
	return list
}
func (p *friendService) GetFriendListRemote() (list []model2.FriendModel) {
	p.db.Select("nick_name", "avatar", "profile", "token", "state", "mark", "block", "version").Where("internal_ip == '' OR internal_ip is null").Find(&list)
	return list
}
func (p *friendService) GetFriendListInternal() (list []model2.FriendModel) {
	p.db.Select("nick_name", "avatar", "profile", "token", "state", "mark", "block", "version").Where("internal_ip != ''").Find(&list)
	return list
}
func (p *friendService) UpdateOrCreate(m model2.FriendModel) {
	friend := model2.FriendModel{}
	p.db.Where("token = ?", m.Token).First(&friend)
	if reflect.DeepEqual(friend, model2.FriendModel{}) {
		p.db.Create(&m)
	} else {
		p.db.Model(&m).Updates(m)
	}

}

func (p *friendService) UpdateAddFriendType(m model2.FriendModel) {
	p.db.Model(&m).Updates(m)
}

func (p *friendService) GetFriendByToken(token string) model2.FriendModel {
	var m model2.FriendModel
	p.db.Model(&m).Where("token = ?", token).First(&m)
	return m
}

func (p *friendService) InternalInspection(ips []string, token string) {
	for _, v := range ips {
		fmt.Println("开始遍历 ip:", v)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		dstAddr, err := net.ResolveUDPAddr("udp", v)
		if err != nil {
			fmt.Println("1", err.Error())
			continue
		}
		port, err := strconv.Atoi(config.ServerInfo.UDPPort)
		if err != nil {
			fmt.Println("2", err)
			continue
		}
		srcAddr := &net.UDPAddr{
			IP: net.IPv4zero, Port: port}
		ticket := token
		session, err := quic.DialContext(ctx, UDPConn, dstAddr, srcAddr.String(), quic_helper.GetClientTlsConfig(ticket), quic_helper.GetQUICConfig())
		if err != nil {
			fmt.Println("3", err, v)
			continue
		}

		stream, err := session.OpenStreamSync(ctx)
		if err != nil {
			fmt.Println("4", err)
			continue
		}
		uuid := uuid.NewV4().String()
		SayHello(stream, token)
		msg := model.MessageModel{
			Type: types.PERSONPING,
			Data: "",
			From: config.ServerInfo.Token,
			To:   token,
			UUId: uuid,
		}

		SendData(stream, msg)

		go ReadContent(stream)
		result := <-Message
		fmt.Println("ping返回结果:", result, msg)
		stream.Close()
		if !reflect.DeepEqual(result, model.MessageModel{}) && result.Data.(string) == token && result.From == token {
			fmt.Println("获取到正确的ip", v)
			UDPAddressMap[result.From] = v
			p.db.Model(&model2.FriendModel{}).Where("token = ?", token).Update("internal_ip", v)
			return
		}
	}
}

func NewFriendService(db *gorm.DB) FriendService {
	return &friendService{db: db}
}
