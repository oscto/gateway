package router

import (
	"context"
	"gateway.oscto.icu/handler/image"
	"github.com/gin-gonic/gin"
	"github.com/oscto/ky3k"
	"go-micro.dev/v4"
	"net/http"
	"strings"
)

var ImageServiceName = "image.oscto.icu"

func Image(r *gin.Engine)  {

	r.GET("/image", func(ctx *gin.Context) {

	})
	r.POST("/image/resize", func(ctx *gin.Context) {
		client := NewImageClient()
		resize, err := client.Resize(context.Background(), &image.CallRequest{
			Url:    ctx.GetString("url"),
			Width:  ctx.GetInt64("width"),
			Height: ctx.GetInt64("height"),
		})
		if err != nil {
			return 
		}

		ctx.JSON(http.StatusOK,resize)
	})
	r.POST("/image/to-webp", func(ctx *gin.Context) {
		client := NewImageClient()
		webp, err := client.ToWebP(context.Background(), &image.ToWebPRequest{Url: ctx.GetString("url")})
		if err != nil {
			return 
		}
		ctx.JSON(http.StatusOK, webp)
	})
	r.POST("/image/draw", func(ctx *gin.Context) {
		client := NewImageClient()
		position := strings.Split(ctx.GetString("position"), ",")
		var x0,x1,y0,y1 int64= 0,0,0,0
		if len(position) >= 4 {
			x0,_ = ky3k.StringToInt64(position[0])
			y0,_ = ky3k.StringToInt64(position[1])
			x1,_ = ky3k.StringToInt64(position[2])
			y1,_ = ky3k.StringToInt64(position[3])
		} else if len(position) >= 2 {
			x1,_ = ky3k.StringToInt64(position[0])
			y1,_ = ky3k.StringToInt64(position[1])
		}
		draw, err := client.Draw(context.Background(), &image.DrawRequest{
			Url: ctx.GetString("url"), X0:  x0, X1:  x1, Y0:  y0, Y1:  y1,})
		if err != nil {
			return 
		}

		ctx.JSON(http.StatusOK, draw)
	})
}

func NewImageClient() image.ImageService {

	srv := micro.NewService()
	return image.NewImageService(ImageServiceName, srv.Client())
}