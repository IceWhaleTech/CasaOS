package v1

import (
	"bytes"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/common_err"
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	sshHelper "github.com/IceWhaleTech/CasaOS-Common/utils/ssh"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	modelCommon "github.com/IceWhaleTech/CasaOS-Common/model"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: time.Duration(time.Second * 5),
}

func PostSshLogin(c *gin.Context) {
	j := make(map[string]string)
	c.ShouldBind(&j)
	userName := j["username"]
	password := j["password"]
	port := j["port"]
	if userName == "" || password == "" || port == "" {
		c.JSON(common_err.CLIENT_ERROR, modelCommon.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS), Data: "Username or password or port is empty"})
		return
	}
	_, err := sshHelper.NewSshClient(userName, password, port)
	if err != nil {
		c.JSON(common_err.CLIENT_ERROR, modelCommon.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: "Please check if the username and port are correct, and make sure that ssh server is installed."})
		logger.Error("connect ssh error", zap.Any("error", err))
		return
	}
	c.JSON(common_err.SUCCESS, modelCommon.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

func WsSsh(c *gin.Context) {
	_, e := exec.LookPath("ssh")
	if e != nil {
		c.JSON(common_err.SERVICE_ERROR, modelCommon.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: "ssh server not found"})
		return
	}

	userName := c.Query("username")
	password := c.Query("password")
	port := c.Query("port")
	wsConn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	logBuff := new(bytes.Buffer)

	quitChan := make(chan bool, 3)
	// user := ""
	// password := ""
	var login int = 1
	cols, _ := strconv.Atoi(c.DefaultQuery("cols", "200"))
	rows, _ := strconv.Atoi(c.DefaultQuery("rows", "32"))
	var client *ssh.Client
	for login != 0 {

		var err error
		if userName == "" || password == "" || port == "" {
			wsConn.WriteMessage(websocket.TextMessage, []byte("username or password or port is empty"))
		}
		client, err = sshHelper.NewSshClient(userName, password, port)

		if err != nil && client == nil {
			wsConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n\x1b[0m"))
		} else {
			login = 0
		}

	}
	if client != nil {
		defer client.Close()
	}

	ssConn, _ := sshHelper.NewSshConn(cols, rows, client)
	defer ssConn.Close()

	go ssConn.ReceiveWsMsg(wsConn, logBuff, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan
}
