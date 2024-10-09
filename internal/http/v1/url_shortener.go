package v1

import (
	managesV1 "github.com/alisher2605/url-shortener/internal/manager/v1"
	"github.com/alisher2605/url-shortener/internal/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type urlShortener struct {
	urlShortenerManager managesV1.UrlShortener
}

func NewUrlShortener(urlShortenerManager managesV1.UrlShortener) *urlShortener {
	return &urlShortener{
		urlShortenerManager: urlShortenerManager,
	}
}

func (u *urlShortener) Init(router *gin.RouterGroup) {
	router.GET("/:urlHash", u.redirectToUrl)
	router.POST("/data/shorten", u.shortenUrl)
}

func (u *urlShortener) shortenUrl(ctx *gin.Context) {
	requestBody := new(model.UrlRequest)

	if err := ctx.ShouldBindJSON(requestBody); err != nil {
		zap.S().Error(err)
		ctx.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: model.ErrInvalidRequestBody.Error(),
		})

		return
	}

	response, err := u.urlShortenerManager.AddUrl(ctx, requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, &model.ShortenedUrlResponse{
		Response: model.Response{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		},
		UrlRequest: response,
	})
}

func (u *urlShortener) redirectToUrl(ctx *gin.Context) {

}
