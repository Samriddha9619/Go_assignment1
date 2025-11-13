package handlers
import (
    "net/http"
    "strconv"
    "github.com/Samriddha9619/Go_assignment1/backend/database"
    "github.com/Samriddha9619/Go_assignment1/backend/models"
    "github.com/Samriddha9619/Go_assignment1/backend/scraper"
    "github.com/gin-gonic/gin"
)

func GetHotels(c *gin.Context) {
    city := c.Param("city")
    sortBy := c.DefaultQuery("sort", "price")
    order := c.DefaultQuery("order", "asc")
    
    var hotels []models.Hotel
    query := database.DB.Where("city = ?", city)
    
    if order == "desc" {
        query = query.Order(sortBy + " desc")
    } else {
        query = query.Order(sortBy)
    }
    
    if err := query.Find(&hotels).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, hotels)
}
func GetAllHotels(c *gin.Context) {
    sortBy := c.DefaultQuery("sort", "city")
    order := c.DefaultQuery("order", "asc")
    
    var hotels []models.Hotel
    query := database.DB
    
    if order == "desc" {
        query = query.Order(sortBy + " desc")
    } else {
        query = query.Order(sortBy)
    }
    
    if err := query.Find(&hotels).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, hotels)
}
func GetCities(c *gin.Context) {
    cities := []string{"Delhi", "Mumbai", "Bangalore"}
    c.JSON(http.StatusOK, gin.H{"cities": cities})
}
func TriggerScrape(c *gin.Context) {
    go scraper.InitialScrape()
    c.JSON(http.StatusOK, gin.H{"message": "Scraping started for all cities"})
}
func GetPriceHistory(c *gin.Context) {
    idStr := c.Param("id")
    hotelID, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
        return
    }
    
    var history []models.PriceHistory
    if err := database.DB.Where("hotel_id = ?", uint(hotelID)).Order("timestamp desc").Limit(100).Find(&history).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, history)
}
func GetScrapingLogs(c *gin.Context) {
    var logs []models.ScrapingLog
    if err := database.DB.Order("started_at desc").Limit(50).Find(&logs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, logs)
}
