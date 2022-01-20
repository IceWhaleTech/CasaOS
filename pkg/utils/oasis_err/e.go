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
	PORT_IS_OCCUPIED    = 20004

	//zerotier
	GET_TOKEN_ERROR = 30001

	//disk
	NAME_NOT_AVAILABLE       = 40001
	DISK_NEEDS_FORMAT        = 40002
	DISK_BUSYING             = 40003
	REMOVE_MOUNT_POINT_ERROR = 40004
	FORMAT_ERROR             = 40005

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
	INVALID_PARAMS:   "Parameters Error",
	ERROR_AUTH_TOKEN: "Error auth token",

	//user
	PWD_INVALID:     "Invalid password",
	PWD_IS_EMPTY:    "Password is empty",
	PWD_INVALID_OLD: "Invalid old password",
	ACCOUNT_LOCK:    "Account is locked",

	//system
	DIR_ALREADY_EXISTS:  "Folder already exists",
	FILE_ALREADY_EXISTS: "File already exists",
	FILE_OR_DIR_EXISTS:  "File or folder already exists",
	PORT_IS_OCCUPIED:    "Port is occupied",

	//zerotier
	GET_TOKEN_ERROR: "Get token error,Please log in to zerotier's official website to confirm whether the account is available",

	//app
	UNINSTALL_APP_ERROR: "Error uninstalling app",
	PULL_IMAGE_ERROR:    "Error pulling image",
	DEVICE_NOT_EXIST:    "Device does not exist",

	//disk
	NAME_NOT_AVAILABLE:       "Name not available",
	DISK_NEEDS_FORMAT:        "Drive needs to be formatted",
	REMOVE_MOUNT_POINT_ERROR: "Failed to remove mount point",
	DISK_BUSYING:             "Drive is busy",
	FORMAT_ERROR:             "Formatting failed, please check if the directory is occupied",

	//
	FILE_DOES_NOT_EXIST: "File does not exist",

	FILE_READ_ERROR:     "File read error",
	SHORTCUTS_URL_ERROR: "URL error",
}

//获取错误信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
