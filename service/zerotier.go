package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/IceWhaleTech/CasaOS/pkg/config"
	command2 "github.com/IceWhaleTech/CasaOS/pkg/utils/command"
	httper2 "github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	"github.com/IceWhaleTech/CasaOS/pkg/zerotier"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ZeroTierService interface {
	GetToken(username, pwd string) string
	ZeroTierRegister(email, lastName, firstName, password string) string
	ZeroTierNetworkList(token string) (interface{}, []string)
	ZeroTierJoinNetwork(networkId string)
	ZeroTierLeaveNetwork(networkId string)
	ZeroTierGetInfo(token, id string) (interface{}, []string)
	ZeroTierGetStatus(token string) interface{}
	EditNetwork(token string, data string, id string) interface{}
	CreateNetwork(token string) interface{}
	MemberList(token string, id string) interface{}
	EditNetworkMember(token string, data string, id, mId string) interface{}
	DeleteMember(token string, id, mId string) interface{}
	DeleteNetwork(token, id string) interface{}
	GetJoinNetworks() string
}
type zerotierstruct struct {
}

var client http.Client

func (c *zerotierstruct) ZeroTierJoinNetwork(networkId string) {
	command2.OnlyExec(`zerotier-cli join ` + networkId)
}
func (c *zerotierstruct) ZeroTierLeaveNetwork(networkId string) {
	command2.OnlyExec(`zerotier-cli leave ` + networkId)
}

//登录并获取token
func (c *zerotierstruct) GetToken(username, pwd string) string {
	if len(config.ZeroTierInfo.Token) > 0 {
		return config.ZeroTierInfo.Token
	} else {
		return LoginGetToken(username, pwd)
	}
}

func (c *zerotierstruct) ZeroTierRegister(email, lastName, firstName, password string) string {

	url := "https://accounts.zerotier.com/auth/realms/zerotier/protocol/openid-connect/registrations?client_id=zt-central&redirect_uri=https%3A%2F%2Fmy.zerotier.com%2Fapi%2F_auth%2Foidc%2Fcallback&response_type=code&scope=openid+profile+email+offline_access&state=state"

	action, cookies, _ := ZeroTierGet(url, nil, 4)
	var buff bytes.Buffer
	buff.WriteString("email=")
	buff.WriteString(email)
	buff.WriteString("&password=")
	buff.WriteString(password)
	buff.WriteString("&password-confirm=")
	buff.WriteString(password)
	buff.WriteString("&user.attributes.marketingOptIn=true")
	buff.WriteString("&firstName")
	buff.WriteString(firstName)
	buff.WriteString("&lastName")
	buff.WriteString(lastName)

	action, errInfo, _ := ZeroTierPost(buff, action, cookies, false)
	if len(errInfo) > 0 {
		return errInfo
	}
	action, _, _ = ZeroTierGet(action, cookies, 5)
	return ""
}

//固定请求head
func GetHead() map[string]string {
	var head = make(map[string]string, 4)
	head["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	head["Accept-Language"] = "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3"
	head["Connection"] = "keep-alive"
	head["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"
	return head
}

//登录并获取token，会出现账号密码错误，和邮箱未验证情况，目前未出现其他情况
func LoginGetToken(username, pwd string) string {
	//拿到登录的action
	var loginUrl = "https://accounts.zerotier.com/auth/realms/zerotier/protocol/openid-connect/auth?client_id=zt-central&redirect_uri=https%3A%2F%2Fmy.zerotier.com%2Fapi%2F_auth%2Foidc%2Fcallback&response_type=code&scope=openid+profile+email+offline_access&state=states"
	action, cookies, _ := ZeroTierGet(loginUrl, nil, 1)
	if len(action) == 0 {
		//没有拿到action，页面结构变了
		return ""
	}
	//登录
	var str bytes.Buffer
	str.WriteString("username=")
	str.WriteString(username)
	str.WriteString("&password=")
	str.WriteString(pwd)
	str.WriteString("&credentialId=&login=Log+In")
	url, logingErrInfo, _ := ZeroTierPost(str, action, cookies, true)

	action, cookies, isLoginOk := ZeroTierGet(url, cookies, 2)

	if isLoginOk {
		//登录成功，可以继续调用api
		randomTokenUrl := "https://my.zerotier.com/api/randomToken"
		json, _, _ := ZeroTierGet(randomTokenUrl, cookies, 3)
		//获取一个随机token
		token := gjson.Get(json, "token")

		userInfoUrl := "https://my.zerotier.com/api/status"
		json, _, _ = ZeroTierGet(userInfoUrl, cookies, 3)
		//拿到用户id
		userId := gjson.Get(json, "user.id")

		//设置新token
		addTokenUrl := "https://my.zerotier.com/api/user/" + userId.String() + "/token"
		data := make(map[string]string)
		rand.Seed(time.Now().UnixNano())
		data["tokenName"] = "oasis-token-" + strconv.Itoa(rand.Intn(1000))
		data["token"] = token.String()
		head := make(map[string]string)
		head["Content-Type"] = "application/json"
		_, statusCode := httper2.ZeroTierPost(addTokenUrl, data, head, cookies)
		if statusCode == http.StatusOK {
			config.Cfg.Section("zerotier").Key("Token").SetValue(token.String())
			config.Cfg.SaveTo("conf/conf.ini")
			config.ZeroTierInfo.Token = token.String()
		}
	} else {
		//登录错误信息
		if len(logingErrInfo) > 0 {
			return logingErrInfo
		} else {
			//验证邮箱
			action, _, _ = ZeroTierGet(url, cookies, 5)
			return "You need to verify your email address to activate your account."
		}
	}
	return ""
}

// t 1:获取action，2：登录成功后拿session（可能需要验证有了或登录失败） 3:随机生成token 4:注册页面拿action  5:注册成功后拿验证邮箱的地址
func ZeroTierGet(url string, cookies []*http.Cookie, t uint8) (action string, c []*http.Cookie, isExistSession bool) {
	isExistSession = false
	action = ""
	c = []*http.Cookie{}
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	for k, v := range GetHead() {
		request.Header.Add(k, v)
	}
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	c = resp.Cookies()
	if t == 1 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return
		}
		action, _ = doc.Find("#kc-form-login").Attr("action")
		return
	} else if t == 2 {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "pgx-session" {
				isExistSession = true
				break
			}
		}
		//判断是否登录成功，如果需要验证邮箱，则返回验证邮箱的地址。
		if resp.StatusCode == http.StatusFound && len(resp.Header.Get("Location")) > 0 {
			action = resp.Header.Get("Location")
		}
		return
	} else if t == 3 {
		//返回获取到的字符串
		byteArr, _ := ioutil.ReadAll(resp.Body)
		action = string(byteArr)
	} else if t == 4 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return
		}
		action, _ = doc.Find("#kc-register-form").Attr("action")
		return

	} else if t == 5 {
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		fmt.Println(doc.Html())
		action, _ = doc.Find("#kc-info-wrapper a").Attr("href")
		return
	}

	return
}

//模拟提交表单
func ZeroTierPost(str bytes.Buffer, action string, cookes []*http.Cookie, isLogin bool) (url, errInfo string, err error) {
	req, err := http.NewRequest(http.MethodPost, action, strings.NewReader(str.String()))
	if err != nil {
		return "", "", errors.New("newrequest error")
	}
	for k, v := range GetHead() {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, cookie := range cookes {
		req.AddCookie(cookie)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return "", "", errors.New("request error")
	}
	if !isLogin {
		//注册成功
		if res.StatusCode == http.StatusFound && len(res.Header.Get("Location")) > 0 {
			return res.Header.Get("Location"), "", nil
		} else {
			register, _ := goquery.NewDocumentFromReader(res.Body)
			firstErr := strings.TrimSpace(register.Find("#input-error-firstname").Text())
			lastErr := strings.TrimSpace(register.Find("#input-error-lastname").Text())
			emailErr := strings.TrimSpace(register.Find("#input-error-email").Text())
			pwdErr := strings.TrimSpace(register.Find("#input-error-password").Text())
			var errD strings.Builder
			if len(firstErr) > 0 {
				errD.WriteString(firstErr + ",")
			}
			if len(lastErr) > 0 {
				errD.WriteString(lastErr + ",")
			}
			if len(emailErr) > 0 {
				errD.WriteString(emailErr + ",")
			}
			if len(pwdErr) > 0 {
				errD.WriteString(pwdErr + ",")
			}
			return "", errD.String(), nil
		}

	} else {
		if res.StatusCode == http.StatusFound && len(res.Header.Get("Location")) > 0 {
			return res.Header.Get("Location"), "", nil
		}
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return "", "", errors.New("request error")
		}

		errDesc := doc.Find("#input-error").Text()
		if len(errDesc) > 0 {
			return "", strings.TrimSpace(errDesc), nil
		}

	}

	return "", "", nil
}

//获取zerotile网络列表和本地用户已加入的网络
func (c *zerotierstruct) ZeroTierNetworkList(token string) (interface{}, []string) {
	url := "https://my.zerotier.com/api/network"
	return zerotier.GetData(url, token), command2.ExecResultStrArray(`zerotier-cli listnetworks | awk 'NR>1 {print $3} {line=$0}'`)
}

// get network info
func (c *zerotierstruct) ZeroTierGetInfo(token, id string) (interface{}, []string) {
	url := "https://my.zerotier.com/api/network/" + id
	info := zerotier.GetData(url, token)
	return info, command2.ExecResultStrArray(`zerotier-cli listnetworks | awk 'NR>1 {print $3} {line=$0}'`)
}

//get status
func (c *zerotierstruct) ZeroTierGetStatus(token string) interface{} {
	url := "https://my.zerotier.com/api/v1/status"
	info := zerotier.GetData(url, token)
	return info
}

func (c *zerotierstruct) EditNetwork(token string, data string, id string) interface{} {
	url := "https://my.zerotier.com/api/v1/network/" + id
	info := zerotier.PostData(url, token, data)
	return info
}

func (c *zerotierstruct) EditNetworkMember(token string, data string, id, mId string) interface{} {
	url := "https://my.zerotier.com/api/v1/network/" + id + "/member/" + mId
	info := zerotier.PostData(url, token, data)
	return info
}

func (c *zerotierstruct) MemberList(token string, id string) interface{} {
	url := "https://my.zerotier.com/api/v1/network/" + id + "/member"
	info := zerotier.GetData(url, token)
	return info
}

func (c *zerotierstruct) DeleteMember(token string, id, mId string) interface{} {
	url := "https://my.zerotier.com/api/v1/network/" + id + "/member/" + mId
	info := zerotier.DeleteMember(url, token)
	return info
}

func (c *zerotierstruct) DeleteNetwork(token, id string) interface{} {
	url := "https://my.zerotier.com/api/v1/network/" + id
	info := zerotier.DeleteMember(url, token)
	return info
}

func (c *zerotierstruct) CreateNetwork(token string) interface{} {
	url := "https://my.zerotier.com/api/v1/network"
	info := zerotier.PostData(url, token, "{}")
	return info
}

func (c *zerotierstruct) GetJoinNetworks() string {
	json := command2.ExecResultStr("source " + config.AppInfo.ProjectPath + "/shell/helper.sh ;GetLocalJoinNetworks")
	return json
}

func NewZeroTierService() ZeroTierService {
	//初始化client
	client = http.Client{Timeout: 30 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse //禁止重定向
	},
	}
	return &zerotierstruct{}
}
