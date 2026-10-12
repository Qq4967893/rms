package bootstrap

import (
	"backend-rms/http/middleware"
	"backend-rms/http/response"
	"backend-rms/http/route"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	registerGlobalMiddleware(router)
	route.RegisterRoutes(router)
	setup404Handler(router)

	return router
}

func registerGlobalMiddleware(router *gin.Engine) {
	router.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.Sentry(),
		middleware.AuthJWT(),
		middleware.RbacHandler(),
	)
}

func setup404Handler(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		response.Abort(c, response.StatusNotFound, response.CodeNotFound,
			response.WithMessage("The route is undefined"),
		)
	})
}
