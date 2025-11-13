package main
import (
    "log"
    "github.com/Samriddha9619/Go_assignment1/backend/auth"
    "github.com/Samriddha9619/Go_assignment1/backend/config"
    "github.com/Samriddha9619/Go_assignment1/backend/database"
    "github.com/Samriddha9619/Go_assignment1/backend/handlers"
    "github.com/Samriddha9619/Go_assignment1/backend/scraper"
    "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
func main() {
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
    cfg := config.Load()
    if err := database.ConnectDB(cfg); err != nil {
        log.Fatalf("Could not connect to database: %v", err)
    }
    defer database.Close()
	handlers.CreateAdminUser()
    log.Println("Starting initial scrape...")
    scraper.InitialScrape()
    scraper.StartScheduler()
    router := gin.Default()
    router.Use(corsMiddleware())
    api := router.Group("/api")
    {
        api.POST("/auth/register", handlers.Register)
        api.POST("/auth/login", handlers.Login)
        
        api.GET("/cities", handlers.GetCities)
        api.GET("/hotels", handlers.GetAllHotels)
        api.GET("/hotels/city/:city", handlers.GetHotels)
        api.GET("/hotels/:id/history", handlers.GetPriceHistory)
    }
    protected := api.Group("")
    protected.Use(auth.AuthRequired())
    {
        protected.GET("/profile", handlers.GetProfile)
        protected.POST("/scrape/trigger", handlers.TriggerScrape)
        protected.GET("/logs", handlers.GetScrapingLogs)
    }
    admin := protected.Group("")
    admin.Use(auth.AdminRequired())
    {
    }
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    log.Printf("Server running on port %s", cfg.ServerPort)
    router.Run(":" + cfg.ServerPort)
}
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
