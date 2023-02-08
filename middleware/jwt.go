package middleware

import (
	"app/lib"
	"app/lib/config"
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func isMethodAllowed(method string, methods []string) bool {
	for i := 0; i < len(methods); i++ {
		if method == methods[i] {
			return true
		}
	}
	return false
}

func matchRules(rules map[string]string, target string, method string) (bool, error) {
	for rule, m := range rules {
		matched, err := regexp.MatchString(rule, target)
		if err != nil {
			return false, err
		}
		if matched {
			return isMethodAllowed(strings.ToLower(method), strings.Split(strings.ToLower(m), "|")), nil
		}
	}
	return false, nil
}

func JWT(unless map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		matched, err := matchRules(unless, c.Request.URL.String(), c.Request.Method)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		if matched {
			c.Next()
			return
		}
		headerStr := c.Request.Header.Get("Authorization")
		if headerStr == "" {
			_ = c.Error(errors.New("授权头信息为空"))
			c.Abort()
			return
		}
		sp := strings.Split(headerStr, "Bearer ")
		if len(sp) <= 1 {
			_ = c.Error(errors.New("授权头信息不合法"))
			c.Abort()
			return
		}
		tokenStr := sp[1]
		tokenStr = strings.TrimSpace(tokenStr)
		token, err := lib.DecodeJWTToken(tokenStr, config.App.JWTSecret)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.Set("auth", token["auth"])
		c.Next()
	}
}
