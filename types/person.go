package types

const PERSONADDFRIEND = "add_user"
const PERSONAGREEFRIEND = "agree_user"
const PERSONDOWNLOAD = "file_data"
const PERSONSUMMARY = "summary"
const PERSONGETIP = "get_ip"
const PERSONCONNECTION = "connection"
const PERSONDIRECTORY = "directory"
const PERSONHELLO = "hello"
const PERSONSHAREID = "share_id"
const PERSONUPLOAD = "upload"
const PERSONUPLOADDATA = "upload_data"
const PERSONINTERNALINSPECTION = "internal_inspection"
const PERSONPING = "ping"
const PERSONIMAGETHUMBNAIL = "image_thumbnail"

const PERSONCANCEL = "cancel" // Cancel Download

const (
	PERSONFILEDOWNLOAD = iota //default state
	PERSONFILEUPLOAD
	PERSONFILERECEIVEUPLOAD //receive upload file
)
