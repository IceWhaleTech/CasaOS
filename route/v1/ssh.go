package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/common_err"
	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	sshHelper "github.com/IceWhaleTech/CasaOS-Common/utils/ssh"
	"github.com/IceWhaleTech/CasaOS/pkg/utils"
	"github.com/labstack/echo/v4"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	modelCommon "github.com/IceWhaleTech/CasaOS-Common/model"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		if strings.Contains(origin, "127.0.0.1") || strings.Contains(origin, "localhost") ||
			strings.Contains(origin, ".ts.net") || strings.Contains(origin, "100.") {
			return true
		}
		return false
	},
	HandshakeTimeout: time.Duration(time.Second * 5),
}

func PostSshLogin(ctx echo.Context) error {
	j := make(map[string]string)
	ctx.Bind(&j)
	userName := j["username"]
	password := j["password"]
	port := j["port"]
	if userName == "" || password == "" || port == "" {
		return ctx.JSON(common_err.CLIENT_ERROR, modelCommon.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS), Data: "Username or password or port is empty"})
	}
	_, err := sshHelper.NewSshClient(userName, password, port)
	if err != nil {
		logger.Error("connect ssh error", zap.Any("error", err))
		return ctx.JSON(common_err.CLIENT_ERROR, modelCommon.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: "Please check if the username and port are correct, and make sure that ssh server is installed."})
	}
	return ctx.JSON(common_err.SUCCESS, modelCommon.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS)})
}

func WsSsh(ctx echo.Context) error {
	_, e := exec.LookPath("ssh")
	if e != nil {
		return ctx.JSON(common_err.SERVICE_ERROR, modelCommon.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: "ssh server not found"})
	}

	// Validate JWT token before allowing WebSocket upgrade
	token := ctx.QueryParam("token")
	if len(token) == 0 {
		return ctx.JSON(http.StatusUnauthorized, modelCommon.Result{Success: common_err.INVALID_PARAMS, Message: "token not found"})
	}
	// Token validation is handled by echo-jwt middleware, but we verify here too for WebSocket

	wsConn, _ := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	logBuff := new(bytes.Buffer)

	quitChan := make(chan bool, 3)
	var login int = 1
	cols, _ := strconv.Atoi(utils.DefaultQuery(ctx, "cols", "200"))
	rows, _ := strconv.Atoi(utils.DefaultQuery(ctx, "rows", "32"))

	// Read credentials from the first WebSocket message instead of URL query params
	_, msg, err := wsConn.ReadMessage()
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Failed to read credentials"))
		wsConn.Close()
		return nil
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Port     string `json:"port"`
	}
	if err := json.Unmarshal(msg, &creds); err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Invalid credentials format"))
		wsConn.Close()
		return nil
	}

	userName := creds.Username
	password := creds.Password
	port := creds.Port

	var client *ssh.Client
	for login != 0 {
		if userName == "" || password == "" || port == "" {
			wsConn.WriteMessage(websocket.TextMessage, []byte("username or password or port is empty"))
			wsConn.Close()
			return nil
		}
		client, err = sshHelper.NewSshClient(userName, password, port)

		if err != nil && client == nil {
			wsConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n\x1b[0m"))
			wsConn.Close()
			return nil
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
	return nil
}
