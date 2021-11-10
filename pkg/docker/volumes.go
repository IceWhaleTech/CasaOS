package docker

import "strings"

func GetDir(id, envName string) string {
	var path string

	if len(id) == 0 {
		id = "$AppID"
	}

	switch {
	case strings.Contains(strings.ToLower(envName), "config") || strings.Contains(strings.ToLower(envName), "photoprism/storage") || strings.Contains(strings.ToLower(envName), "config"):
		path = "/DATA/AppData/" + id + "/"
	case strings.Contains(strings.ToLower(envName), "media"):
		path = "/DATA/Media/"
	case strings.Contains(strings.ToLower(envName), "movie"):
		path = "/DATA/Media/Movies/"
	case strings.Contains(strings.ToLower(envName), "music"):
		path = "/DATA/Media/Music/"
	case strings.Contains(strings.ToLower(envName), "photoprism/originals"):
		path = "/DATA/Gallery"
	case strings.Contains(strings.ToLower(envName), "download"):
		path = "/DATA/Downloads/"
	case strings.Contains(strings.ToLower(envName), "photo") || strings.Contains(strings.ToLower(envName), "pictures"):
		path = "/DATA/Downloads/"
	case strings.ToLower(envName) == "/srv":
		path = "/DATA/"
	default:
		//path = "/media"
	}
	return path
}
