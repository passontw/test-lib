package middlewares

/**
 * @Author: M
 * @Date: 2024/7/31 9:42
 * @Desc:
 */
import (
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
)

func AuthMiddleware(ctx *context.Context) {
	token := ctx.Input.Header("Authorization")
	if token == "" || !validateToken(token) {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.Body([]byte("Unauthorized"))
		ctx.Abort(401, "Unauthorized")
	}
}

func validateToken(token string) bool {
	// Implement your token validation logic here
	return token == "valid-token"
}
