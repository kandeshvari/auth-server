package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kandeshvari/gin-jwt-middleware"
)

func RunAuthServer() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// setup CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "UPDATE", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Auth-Request"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// setup JWT infrastructure
	rts, err := NewRefreshTokenStorage(ctx.conf.Database.Driver, ctx.conf.Database.DbConnect)
	if err != nil {
		return err
	}
	authMW := &jwt.GinMiddleware{
		SecretKey:           ctx.conf.Auth.SecretKey,
		Timeout:             ctx.conf.Auth.Timeout * time.Minute,
		RefreshTimeout:      ctx.conf.Auth.RefreshTimeout * time.Minute,
		Authenticator:       ctx.auth.Authenticator,
		RefreshTokenStorage: rts,
	}

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authMW.LoginHandler)
			authGroup.GET("/refresh", authMW.RefreshHandler)
			authGroup.DELETE("/logout", authMW.LogoutHandler)
		}
	}

	if err := http.ListenAndServe(ctx.conf.Server.Address, r); err != nil {
		return err
	}

	return nil
}
