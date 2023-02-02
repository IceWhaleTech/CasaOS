package dropbox

import (
	"time"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/model"
	"go.uber.org/zap"
)

type UserInfo struct {
	AccountID string `json:"account_id"`
	Name      struct {
		GivenName       string `json:"given_name"`
		Surname         string `json:"surname"`
		FamiliarName    string `json:"familiar_name"`
		DisplayName     string `json:"display_name"`
		AbbreviatedName string `json:"abbreviated_name"`
	} `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Disabled      bool   `json:"disabled"`
	Country       string `json:"country"`
	Locale        string `json:"locale"`
	ReferralLink  string `json:"referral_link"`
	IsPaired      bool   `json:"is_paired"`
	AccountType   struct {
		Tag string `json:".tag"`
	} `json:"account_type"`
	RootInfo struct {
		Tag             string `json:".tag"`
		RootNamespaceID string `json:"root_namespace_id"`
		HomeNamespaceID string `json:"home_namespace_id"`
	} `json:"root_info"`
}
type TokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
type File struct {
	Tag            string    `json:".tag"`
	Name           string    `json:"name"`
	PathLower      string    `json:"path_lower"`
	PathDisplay    string    `json:"path_display"`
	ID             string    `json:"id"`
	ClientModified time.Time `json:"client_modified,omitempty"`
	ServerModified time.Time `json:"server_modified,omitempty"`
	Rev            string    `json:"rev,omitempty"`
	Size           int       `json:"size,omitempty"`
	IsDownloadable bool      `json:"is_downloadable,omitempty"`
	ContentHash    string    `json:"content_hash,omitempty"`
}

type Files struct {
	Files   []File `json:"entries"`
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
}

type Error struct {
	Error struct {
		Errors []struct {
			Domain       string `json:"domain"`
			Reason       string `json:"reason"`
			Message      string `json:"message"`
			LocationType string `json:"location_type"`
			Location     string `json:"location"`
		}
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func fileToObj(f File) *model.ObjThumb {
	logger.Info("dropbox file", zap.Any("file", f))
	obj := &model.ObjThumb{
		Object: model.Object{
			ID:       f.ID,
			Name:     f.Name,
			Size:     int64(f.Size),
			Modified: f.ClientModified,
			IsFolder: f.Tag == "folder",
			Path:     f.PathDisplay,
		},
		Thumbnail: model.Thumbnail{},
	}
	return obj
}
