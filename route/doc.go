//go:build doc
// +build doc

package route

import (
	_ "github.com/IceWhaleTech/CasaOS/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

func init() {
	//	swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
	swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
}
