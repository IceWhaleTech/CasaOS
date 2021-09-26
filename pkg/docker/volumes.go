package docker

func GetDir(id, envName string) string {
	var path string
	switch envName {
	case "/config":
		path = "/oasis/app_data/" + id + "/"
	default:
		//path = "/media"
	}
	return path
}
