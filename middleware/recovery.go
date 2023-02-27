package middleware

import (
	"errors"
	"fmt"
	"github/inoth/ino-gateway/res"
	"github/inoth/ino-gateway/util"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if util.GetRunEnv() == util.EnvDev {
					fmt.Println(string(debug.Stack()))
				}
				switch e := err.(type) {
				case error:
					res.ResultErr(c, res.InternalErrorCode, e)
				default:
					res.ResultErr(c, res.InternalErrorCode, errors.New("internal server error"))
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
