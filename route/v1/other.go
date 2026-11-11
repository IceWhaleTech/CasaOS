package v1

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/common_err"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/labstack/echo/v4"
)

// allowedSearchDomains defines which domains AgentSearch is allowed to access
var allowedSearchDomains = []string{
	"www.google.com",
	"www.bing.com",
	"api.duckduckgo.com",
	"duckduckgo.com",
}

func GetSearchResult(ctx echo.Context) error {
	json := make(map[string]string)
	ctx.Bind(&json)
	rawURL := json["url"]

	if rawURL == "" {
		return ctx.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: common_err.GetMsg(common_err.INVALID_PARAMS), Data: "key is empty"})
	}

	// Validate URL against allowed domains
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ctx.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: "invalid URL"})
	}

	hostname := parsedURL.Hostname()
	allowed := false
	for _, domain := range allowedSearchDomains {
		if hostname == domain || strings.HasSuffix(hostname, "."+domain) {
			allowed = true
			break
		}
	}
	if !allowed {
		return ctx.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: "URL domain not allowed"})
	}

	if file.IsPrivateOrReservedIP(hostname) {
		return ctx.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.INVALID_PARAMS, Message: "access to internal resources is forbidden"})
	}

	data, err := service.MyService.Other().AgentSearch(rawURL)
	if err != nil {
		fmt.Println(err)
		return ctx.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR)})
	}

	return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: data})
}
