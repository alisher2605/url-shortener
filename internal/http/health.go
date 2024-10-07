package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (srv *server) healthz(ctx *gin.Context) {
	ctx.String(http.StatusOK, http.StatusText(http.StatusOK))
}
