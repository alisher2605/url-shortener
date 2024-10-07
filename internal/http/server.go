package http

import (
	_ "github.com/alisher2605/url-shortener/api/swagger"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

const (
	v1Prefix = "/v1"
)

type server struct {
	maxAge  int
	appPort string
	router  *gin.Engine
}

func NewServer(appPort string, maxAge int) *server {
	return &server{
		appPort: appPort,
		maxAge:  maxAge,
		router:  gin.New(),
	}
}

func (srv *server) setupRouter() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Println(err.Error())
	}

	srv.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Origin"},
		AllowCredentials: true,
		MaxAge:           time.Duration(srv.maxAge),
	}))

	srv.router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{TimeFormat: time.RFC3339, UTC: false, SkipPaths: []string{"/", "/health"}}))
	srv.router.GET("/healthz", srv.healthz)

	//v1 := srv.router.Group(v1Prefix)

	srv.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (srv *server) Run() {
	srv.setupRouter()

	if err := srv.router.Run(":" + srv.appPort); err != nil {
		zap.S().Fatal("Couldn't run HTTP server")
	}
}
