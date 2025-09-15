// Package main provides the main entry point for the Gopher Social API
//
//	@title						Gopher Social API
//	@version					1.1.0
//	@description				A social media API service built with Go featuring user management, posts, comments, and following functionality
//	@termsOfService				http://swagger.io/terms/
//
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host						localhost:8080
//	@BasePath					/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your bearer token in the format **Bearer &lt;token&gt;**
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-KEY
//	@description				API Key for authorization
package main

import (
	httpserver "github.com/orangeMangoDimz/go-social/internal/server/http"
	"github.com/orangeMangoDimz/go-social/internal/service/domain"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres"
)

const VERSION = "1.1.0"

func main() {

	db, app := httpserver.NewApp()
	repositories := postgres.NewStore(db)
	services := domain.NewService(repositories, app.Logger, app.Config)
	app.Services = *services

	mux := app.Mount(VERSION)
	err := app.Run(mux, VERSION)
	if err != nil {
		app.Logger.Errorw("Failed to start server", "error", err)
	}
}
