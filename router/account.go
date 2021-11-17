package router

import (
	"context"
	account "gateway.oscto.icu/handler"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
	"net/http"
)

var AccountServiceName = "account.oscto.icu"

func Account(r *gin.Engine)  {
	r.GET("/hello", func(ctx *gin.Context) {
		c := NewAccountClient()
		login, err := c.Login(context.Background(), &account.LoginRequest{Email: "email",Password: "password"})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		ctx.JSONP(http.StatusOK, login)
	})
}

func NewAccountClient() account.AccountService {

	service := micro.NewService()
	client := account.NewAccountService(AccountServiceName, service.Client())

	return client
}