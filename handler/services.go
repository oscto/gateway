package handler

import (
	"context"
	"git.vonechain.com/vone-bfs/gateway/proto/user"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/client"
	"net/http"
)

type Services struct {
	user user.UserService
}

func NewServices(c client.Client) *Services {
	return &Services{
		user: user.NewUserService("user", c),
	}
}

func (s *Services) Register(c *gin.Context) {
	req := &user.CallRequest{}
	rsq, err := s.user.Call(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": rsq})
}
