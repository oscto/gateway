package router

import (
	"context"
	"fmt"
	"net/http"

	"gateway.oscto.icu/handler/image"
	"gateway.oscto.icu/metadata"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
)

var ImageServiceName = "image.oscto.icu"

func Image(r *gin.Engine) {

	r.GET("/image", func(ctx *gin.Context) {

	})
	r.POST("/image/resize", func(ctx *gin.Context) {

		var metaImage metadata.ImageResizeRequest
		if err := ctx.BindJSON(&metaImage); err != nil {
			fmt.Println("bind error", err)
			return
		}
		client := NewImageClient()
		resize, err := client.Resize(context.Background(), &image.CallRequest{
			Url:    metaImage.Url,
			Width:  metaImage.Width,
			Height: metaImage.Height,
		})
		if err != nil {
			fmt.Println("resize, error", err)
			return
		}

		ctx.JSON(http.StatusOK, resize)
	})
	r.POST("/image/to-webp", func(ctx *gin.Context) {

		var metaImage metadata.ImageToWebPRequest
		if err := ctx.BindJSON(&metaImage); err != nil {
			return
		}
		client := NewImageClient()
		webp, err := client.ToWebP(context.Background(), &image.ToWebPRequest{Url: metaImage.Url})
		if err != nil {
			fmt.Println("webp error", err)
			return
		}
		ctx.JSON(http.StatusOK, webp)
	})
	r.POST("/image/draw", func(ctx *gin.Context) {
		var metaImage metadata.ImageDrawRequest
		if err := ctx.BindJSON(&metaImage); err != nil {
			fmt.Println("bind error", err)
			return
		}
		client := NewImageClient()
		draw, err := client.Draw(context.Background(), &image.DrawRequest{
			Url: metaImage.Url,
			X0:  metaImage.X0,
			X1:  metaImage.X1,
			Y0:  metaImage.Y0,
			Y1:  metaImage.Y1,
		})
		if err != nil {
			fmt.Println("draw err", err)
			return
		}
		ctx.JSON(http.StatusOK, draw)
	})
}

func NewImageClient() image.ImageService {

	srv := micro.NewService()
	return image.NewImageService(ImageServiceName, srv.Client())
}
