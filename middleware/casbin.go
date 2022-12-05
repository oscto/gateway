package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gin-gonic/gin"
	"net/http"
)

var modelString = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = priority(p.eft) || deny

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`

func Casbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		m, err := model.NewModelFromString(modelString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "system error!"})
		}
		a := fileadapter.NewAdapter("casbin/policy.csv")
		e, err := casbin.NewEnforcer(m, a)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			return
		}
		ok, err := e.Enforce("bob", "data1", "read1")
		if !ok || err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "access deny"})
			return
		}
		// 请求前
		c.Next()
		// 请求后
		//....
	}
}
