package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/2mf8/QQBotOffical/middleware"
	"github.com/2mf8/QQBotOffical/status"
	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

// curl -X POST -H "Content-Type:application/json" -d "{\"username\":\"admin\", \"code\": \"123456\"}" http://localhost:8080/login
func TokenGetApi(c *gin.Context) {
	l := Login{}
	len := c.Request.ContentLength
	body := make([]byte, len)
	c.Request.Body.Read(body)
	fmt.Println(len, body, string(body))
	e := json.Unmarshal(body, &l)
	fmt.Println(e, l)
	if l.Username == "admin" && l.Code == "123456" {
		ts, es := middleware.GenTokens(l.Username, 20)
		if es[0] != nil {
			c.JSON(int(status.ExpectationFailed), gin.H{
				"code": status.GetTokenError,
				"msg":  "获取Token失败，请重试",
			})
		} else {
			/*c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"msg":     "登录成功",
				"token":   ts[0],
				"refresh": ts[1],
			})*/
			c.String(http.StatusOK, `curl -H "Authorization: Bearer %s" -H "Refresh: Bearer %s" http://localhost:8080/prices/四`, ts[0], ts[1])
		}
	} else {
		c.JSON(int(status.ExpectationFailed), gin.H{
			"code": status.LoginError,
			"msg":  "获取Token失败，用户名错误或密码错误",
		})
	}
}
