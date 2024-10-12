package api

import (
	"cdn/src/api/http/middlewares"
	"cdn/src/api/routes"
	"cdn/src/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
)

var (
	configs      = config.GetInstance()
	isProduction = configs.Get("APP_ENV") == "production"
	g            errgroup.Group
)

func Init() (err error) {

	g.Go(func() error {
		return initServer()
	})

	if err = g.Wait(); err != nil {
		log.Fatalln(err)
		return err
	}

	return err
}

func getNewRouter() *gin.Engine {

	// set gin to release mode.
	gin.SetMode(gin.ReleaseMode)

	// Initialize new app.
	router := gin.New()

	// Attach CORS middleware.
	router.Use(middlewares.Cors())

	// Attach logger middleware.
	router.Use(gin.Logger())

	// Attach recovery middleware.
	router.Use(gin.Recovery())

	// Attach request id middleware.
	router.Use(middlewares.RequestID)

	router.Use(middlewares.NewIPMiddleware().Middleware())

	router.Use(middlewares.NewApiKeyMiddleware().Middleware())

	if isProduction {
		// Trusted proxies.
		_ = router.SetTrustedProxies([]string{"https://" + configs.Get("APP_HOST")})
	}

	return router
}

func initServer() error {
	router := getNewRouter()

	v1 := router.Group("api/v1")
	{
		routes.ObjectRoutes(v1)
		routes.BucketRoutes(v1)
		routes.SwaggerRoutes(v1)
	}

	// Run App.
	if err := router.RunTLS(
		fmt.Sprintf(":%s", configs.Get("APP_PORT")),
		configs.Get("SSL_CERT_PATH"),
		configs.Get("SSL_KEY_PATH"),
	); err != nil {
		return err
	}

	return nil
}
