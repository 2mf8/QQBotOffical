package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/2mf8/QQBotOffical/status"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("@#$%^&2mf8kequ._AFJK")
var jwtRefreshKey = []byte("jku838%$_.djkjghjd")
var setUsername = ""

type JwtClaims struct {
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

func GetKeys() [2][]uint8 {
	var bs [2][]uint8
	bs[0] = jwtKey
	bs[1] = jwtRefreshKey
	return bs
}

func GenJwtClaims(username string, timeout int) JwtClaims {
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(timeout)))
	claims := JwtClaims{
		Username: username,
		Role:     1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
		},
		//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(timeout))),
	}
	return claims
}

func ReGenJwtClaims(username string, timeout int) JwtClaims {
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(timeout)))
	claims := JwtClaims{
		Username: username,
		Role:     1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
		},
		//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(timeout))),
	}
	return claims
}

func (j *JwtClaims) GenToken(key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, j)
	return token.SignedString(key)
}

func GenTokens(username string, timeout int) ([2]string, [2]error) {
	var tokens [2]string
	var errs [2]error
	tg := GenJwtClaims(username, timeout)
	t, e1 := tg.GenToken(jwtKey)
	tokens[0] = t
	errs[0] = e1
	rtg := GenJwtClaims(username, timeout*3)
	rt, e2 := rtg.GenToken(jwtRefreshKey)
	tokens[1] = rt
	errs[1] = e2
	return tokens, errs
}

func RefreshTokens(refreshTokens [2]string, timeout int) ([2]string, string, [2]error) {
	tokens := [2]string{}
	errs := [2]error{}
	bs := GetKeys()
	if refreshTokens[0] == "" {
		status := "token为空，请正常登录"
		return tokens, status, errs
	}
	_r, err := ParseToken(refreshTokens[0], bs[0])
	fmt.Println(err)
	if err != nil {
		rj, e := ParseToken(refreshTokens[1], bs[1])
		fmt.Println(rj, e)
		if e != nil {
			status := "token已过期，请正常登录"
			return tokens, status, errs
		}
		tg := ReGenJwtClaims(rj.Username, timeout)
		t, e1 := tg.GenToken(bs[0])
		tokens[0] = t
		errs[0] = e1
		setUsername = rj.Username
		rtg := ReGenJwtClaims(rj.Username, timeout*3)
		rt, e2 := rtg.GenToken(bs[1])
		tokens[1] = rt
		errs[1] = e2
		status := "token刷新成功"
		fmt.Printf(`curl -H "Authorization: Bearer %s" -H "Refresh: Bearer %s" http://localhost:8080/price/四`, tokens[0], tokens[1])
		return tokens, status, errs
	}
	setUsername = _r.Username
	status := "登录状态正常"
	tokens = refreshTokens
	return tokens, status, errs
}

func ParseToken(tokenString string, key []byte) (*JwtClaims, error) {
	var j = new(JwtClaims)
	token, err := jwt.ParseWithClaims(tokenString, j, func(t *jwt.Token) (interface{}, error) { return key, nil })
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return j, nil
	}
	return nil, errors.New("invaild token")
}

func JWTAuthMiddlewareOrRefreshToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		tokens := [2]string{}
		authHeader := c.Request.Header.Get("Authorization")
		fmt.Println(authHeader)
		if c.Request.Method == "POST" {
			if authHeader == "" {
				c.JSON(int(status.Unauthorized), gin.H{
					"code": status.TokenNull,
					"msg":  "访问失败，token为空，请登录",
				})
				c.Abort()
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				c.JSON(int(status.Unauthorized), gin.H{
					"code": status.TokenInvalid,
					"msg":  "访问失败，无效的token，请登录。",
				})
				c.Abort()
				return
			}
			tokens[0] = parts[1]
			refreshHeader := c.Request.Header.Get("Refresh")
			fmt.Println(refreshHeader)
			if refreshHeader != "" {
				rparts := strings.SplitN(refreshHeader, " ", 2)
				if len(rparts) == 2 && rparts[0] == "Bearer" {
					tokens[1] = rparts[1]
				}
			}
			_, _status, _ := RefreshTokens(tokens, 20)
			if !strings.Contains(_status, "状态正常") {
				c.JSON(int(status.ExpectationFailed), gin.H{
					"code": 200,
					"msg":  _status,
				})
				c.Abort()
				return
			}
			c.Set("Username", setUsername)
			c.Next()
		} else {
			c.JSON(int(status.MethodNotAllowed), gin.H{
				"code": status.MethodNotAllowed,
				"msg":  "请使用 POST 方法提交",
			})
			c.Abort()
			return
		}
	}
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		fmt.Println(authHeader)
		if authHeader == "" {
			c.JSON(int(status.Unauthorized), gin.H{
				"code": status.TokenNull,
				"msg":  "访问失败，token为空，请登录",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(int(status.Unauthorized), gin.H{
				"code": status.TokenInvalid,
				"msg":  "访问失败，无效的token，请登录。",
			})
			c.Abort()
			return
		}
		ms, err := ParseToken(parts[1], jwtKey)
		if err != nil {
			refreshHeader := c.Request.Header.Get("Refresh")
			fmt.Println(refreshHeader)
			if refreshHeader == "" {
				c.JSON(int(status.Unauthorized), gin.H{
					"code": status.RefreshTokenNull,
					"msg":  "访问失败，refreshToken为空，请登录",
				})
				c.Abort()
				return
			}
			rparts := strings.SplitN(refreshHeader, " ", 2)
			if !(len(rparts) == 2 && rparts[0] == "Bearer") {
				c.JSON(int(status.Unauthorized), gin.H{
					"code": status.RefreshTokenInvalid,
					"msg":  "访问失败，无效的refreshToken，请登录。",
				})
				c.Abort()
				return
			}
			rms, err := ParseToken(rparts[1], jwtRefreshKey)
			if err != nil {
				c.JSON(int(status.Unauthorized), gin.H{
					"code": status.RefreshTokenInvalid,
					"msg":  "访问失败，无效的refreshToken，请登录。",
				})
				c.Abort()
				return
			}
			c.Set("Username", rms.Username)
			c.Next()
		} else {
			c.Set("Username", ms.Username)
			c.Next()
		}
	}
}
