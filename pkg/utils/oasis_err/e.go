package oasis_err

const (
	SUCCESS          = 200
	ERROR            = 500
	INVALID_PARAMS   = 400
	ERROR_AUTH_TOKEN = 401

	//user
	PWD_INVALID  = 10001
	PWD_IS_EMPTY = 10002

	PWD_INVALID_OLD = 10003
	ACCOUNT_LOCK    = 10004
	//system
	DIR_ALREADY_EXISTS  = 20001
	FILE_ALREADY_EXISTS = 20002
	FILE_OR_DIR_EXISTS  = 20003

	//zerotier
	GET_TOKEN_ERROR = 30001

	//app
	UNINSTALL_APP_ERROR = 50001
	PULL_IMAGE_ERROR    = 50002
	DEVICE_NOT_EXIST    = 50003

	//file
	FILE_DOES_NOT_EXIST = 60001
	FILE_READ_ERROR     = 60002

	//shortcuts
	SHORTCUTS_URL_ERROR = 70001
)

var MsgFlags = map[int]string{
	SUCCESS:          "ok",
	ERROR:            "fail",
	INVALID_PARAMS:   "Invalid params",
	ERROR_AUTH_TOKEN: "error auth token",

	//user
	PWD_INVALID:     "Password invalid",
	PWD_IS_EMPTY:    "Password is empty",
	PWD_INVALID_OLD: "Old Password invalid",
	ACCOUNT_LOCK:    "Account Lock",

	//system
	DIR_ALREADY_EXISTS:  "Directory already exists",
	FILE_ALREADY_EXISTS: "File already exists",
	FILE_OR_DIR_EXISTS:  "File or directory already exists",

	//zerotier
	GET_TOKEN_ERROR: "Get token error,Please log in to zerotier's official website to confirm whether the account is available",

	//app
	UNINSTALL_APP_ERROR: "uninstall app error",
	PULL_IMAGE_ERROR:    "pull image error",
	DEVICE_NOT_EXIST:    "device not exist",

	//
	FILE_DOES_NOT_EXIST: "file does not exist",

	FILE_READ_ERROR:     "file read error",
	SHORTCUTS_URL_ERROR: "url error",
}

//获取错误信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
