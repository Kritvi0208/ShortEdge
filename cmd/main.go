// @title URLify - Branded URL Shortener with Insights
// @version 1.0
// @description API documentation for URLify service
// @contact.name Ritvi K
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"os"

	_ "github.com/Kritvi0208/ShortEdge/docs" // 🔥 auto-generated docs will go here
	"github.com/Kritvi0208/ShortEdge/factory"
	"github.com/Kritvi0208/ShortEdge/handler"
	"github.com/Kritvi0208/ShortEdge/middleware"
	"github.com/Kritvi0208/ShortEdge/service"

	"github.com/joho/godotenv"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
	//httpSwagger "github.com/swaggo/http-swagger"
	"gofr.dev/pkg/gofr"
)

func main() {
	err := godotenv.Load("F:/ShortEdge/.env")
	if err != nil {
		log.Fatalf("❌ Failed to load .env file: %v", err)
	}

	log.Println("✅ DB_URL from .env is:", os.Getenv("DB_URL")) // debug print
	// ✅ Initialize GoFr app

	//os.Setenv("GOFR_DB_URL", os.Getenv("DB_URL"))

	app := gofr.New()

	// Visit Analytics Dependencies
	visitStore := factory.NewVisitStore(app)
	visitService := service.NewVisitService(visitStore)
	visitHandler := handler.NewVisitHandler(visitService)

	// URL Shortener Dependencies
	urlStore := factory.NewURLStore(app)
	urlService := service.New(urlStore)
	urlHandler := handler.NewURLHandler(urlService, visitService)
	//fileserver, router.handle, promhttp, metricshandler
	// Routes
	//app.Server().Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui"))))
	//app.GET("/swagger/*", gofrSwagger.NewHandler())
	app.GET("/all", middleware.RedirectMiddleware(urlHandler.GetAll))
	app.GET("/health", handler.HealthHandler)
	//app.Router.Handle("/metrics", http.HandlerFunc(promhttp.Handler().ServeHTTP))
	//app.GET("/metrics", app.MetricsHandler())
	app.POST("/shorten", middleware.RedirectMiddleware(urlHandler.Shorten))
	app.PUT("/update/{code}", middleware.RedirectMiddleware(urlHandler.Update))
	app.DELETE("/delete/{code}", middleware.RedirectMiddleware(urlHandler.Delete))
	app.GET("/analytics/{code}", middleware.RedirectMiddleware(visitHandler.GetAnalytics))
	app.GET("/{code}", middleware.RedirectMiddleware(urlHandler.Redirect))

	app.Run()
}
