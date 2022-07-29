package common_err

const (
	SUCCESS          = 200
	SERVICE_ERROR    = 500
	CLIENT_ERROR     = 400
	ERROR_AUTH_TOKEN = 401

	INVALID_PARAMS = 4000
	//user
	PWD_INVALID              = 10001
	PWD_IS_EMPTY             = 10002
	PWD_INVALID_OLD          = 10003
	ACCOUNT_LOCK             = 10004
	PWD_IS_TOO_SIMPLE        = 10005
	USER_NOT_EXIST           = 10006
	USER_EXIST               = 10007
	KEY_NOT_EXIST            = 10008
	NOT_IMAGE                = 10009
	IMAGE_TOO_LARGE          = 10010
	INSUFFICIENT_PERMISSIONS = 10011

	//system
	DIR_ALREADY_EXISTS              = 20001
	FILE_ALREADY_EXISTS             = 20002
	FILE_OR_DIR_EXISTS              = 20003
	PORT_IS_OCCUPIED                = 20004
	COMMAND_ERROR_INVALID_OPERATION = 20005
	VERIFICATION_FAILURE            = 20006
	Record_NOT_EXIST                = 20007
	Record_ALREADY_EXIST            = 20008
	SERVICE_NOT_RUNNING             = 20009

	//disk
	NAME_NOT_AVAILABLE       = 40001
	DISK_NEEDS_FORMAT        = 40002
	DISK_BUSYING             = 40003
	REMOVE_MOUNT_POINT_ERROR = 40004
	FORMAT_ERROR             = 40005

	//app
	UNINSTALL_APP_ERROR  = 50001
	PULL_IMAGE_ERROR     = 50002
	DEVICE_NOT_EXIST     = 50003
	ERROR_APP_NAME_EXIST = 50004

	//file
	FILE_DOES_NOT_EXIST = 60001
	FILE_READ_ERROR     = 60002
	FILE_DELETE_ERROR   = 60003
	DIR_NOT_EXISTS      = 60004
	SOURCE_DES_SAME     = 60005

	//share
	SHARE_ALREADY_EXISTS      = 70001
	SHARE_NAME_ALREADY_EXISTS = 70002
)

var MsgFlags = map[int]string{
	SUCCESS:          "ok",
	SERVICE_ERROR:    "Fail",
	CLIENT_ERROR:     "Fail",
	INVALID_PARAMS:   "Parameters Error",
	ERROR_AUTH_TOKEN: "Error auth token",

	//user
	PWD_INVALID:              "Invalid password",
	PWD_IS_EMPTY:             "Password is empty",
	PWD_INVALID_OLD:          "Invalid old password",
	ACCOUNT_LOCK:             "Account is locked",
	PWD_IS_TOO_SIMPLE:        "Password is too simple",
	USER_NOT_EXIST:           "User does not exist",
	USER_EXIST:               "User already exists",
	KEY_NOT_EXIST:            "Key does not exist",
	IMAGE_TOO_LARGE:          "Image is too large",
	NOT_IMAGE:                "Not an image",
	INSUFFICIENT_PERMISSIONS: "Insufficient permissions",

	//system
	DIR_ALREADY_EXISTS:   "Folder already exists",
	FILE_ALREADY_EXISTS:  "File already exists",
	FILE_OR_DIR_EXISTS:   "File or folder already exists",
	PORT_IS_OCCUPIED:     "Port is occupied",
	VERIFICATION_FAILURE: "Verification failure",
	Record_ALREADY_EXIST: "Record already exists",
	Record_NOT_EXIST:     "Record does not exist",
	SERVICE_NOT_RUNNING:  "Service is not running",

	//app
	UNINSTALL_APP_ERROR:  "Error uninstalling app",
	PULL_IMAGE_ERROR:     "Error pulling image",
	DEVICE_NOT_EXIST:     "Device does not exist",
	ERROR_APP_NAME_EXIST: "App name already exists",

	//disk
	NAME_NOT_AVAILABLE:       "Name not available",
	DISK_NEEDS_FORMAT:        "Drive needs to be formatted",
	REMOVE_MOUNT_POINT_ERROR: "Failed to remove mount point",
	DISK_BUSYING:             "Drive is busy",
	FORMAT_ERROR:             "Formatting failed, please check if the directory is occupied",
	//share
	SHARE_ALREADY_EXISTS:      "Share already exists",
	SHARE_NAME_ALREADY_EXISTS: "Share name already exists",
	//
	SOURCE_DES_SAME:     "Source and destination cannot be the same.",
	FILE_DOES_NOT_EXIST: "File does not exist",

	DIR_NOT_EXISTS: "Directory does not exist",

	FILE_READ_ERROR:   "File read error",
	FILE_DELETE_ERROR: "Delete error",

	COMMAND_ERROR_INVALID_OPERATION: "invalid operation",
}

//获取错误信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[SERVICE_ERROR]
}
