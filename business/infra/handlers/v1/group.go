package v1

import (
	_ "myc-devices-simulator/business/infra/handlers/v1/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// CreateGroupV1 init handlers group v1.
// @title Swagger MYC-DEVICES-SIMULATOR API V1
// @version 1.0
// @description Devices simulator documentation API V1.
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization.
func CreateGroupV1(app *echo.Group) *echo.Group {
	v1 := app.Group("v1")

	// Initialize swagger.
	v1.GET("/docs/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("v1")))

	return v1
}
