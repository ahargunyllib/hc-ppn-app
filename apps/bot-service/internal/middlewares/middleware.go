package middlewares

import "github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/jwt"

type Middleware struct {
	jwt jwt.CustomJwtInterface
}

func NewMiddleware(
	jwt jwt.CustomJwtInterface,
) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}
