package v2

import (
	"net/http"

	"github.com/IceWhaleTech/CasaOS/codegen"
	"github.com/labstack/echo/v4"
)

// Path: route/v2/file.go

func (s *CasaOS) GetFileTest(ctx echo.Context) error {

	//http.ServeFile(w, r, r.URL.Path[1:])
	http.ServeFile(ctx.Response().Writer, ctx.Request(), "/DATA/test.img")

	return ctx.String(200, "pong")
}

func (c *CasaOS) CheckUploadChunk(ctx echo.Context, params codegen.CheckUploadChunkParams) error {
	return c.fileUploadService.TestChunk(ctx)
}

func (c *CasaOS) PostUploadFile(ctx echo.Context) error {
	return c.fileUploadService.UploadFile(ctx)
}
