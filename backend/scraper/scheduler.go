package scraper
import (
    "log"
    "time"
    "github.com/Samriddha9619/Go_assignment1/backend/database"
    "github.com/Samriddha9619/Go_assignment1/backend/models"
    "github.com/robfig/cron/v3"
    "gorm.io/gorm/clause"
)

func StartScheduler() {
    c := cron.New()
    
    c.AddFunc("0 */2 * * *", func() {
        log.Println("Scheduled price update started...")
        updateAllCities()
        log.Println("Scheduled price update completed!")
    })
    c.Start()
    log.Println("Scheduler started - updates every 2 hours")
}
func updateAllCities() {
    cities := []string{"Delhi", "Mumbai", "Bangalore"}
    for _, city := range cities {
        scrapingLog := &models.ScrapingLog{
            City:      city,
            Status:    "in_progress",
            StartedAt: time.Now(),
        }
        if err := database.DB.Create(scrapingLog).Error; err != nil {
            log.Printf("Error creating scraping log: %v", err)
            continue
        }
        hotels := ScrapeCity(city)
        
        for i := range hotels {
            hotels[i].Price *= (1.0 + (0.1 * (0.5 - float64(time.Now().UnixNano()%1000)/1000)))
        }
        errorMsg := ""
        successCount := 0
        for _, hotel := range hotels {
            if err := upsertHotel(hotel); err != nil {
                log.Printf("Error upserting %s: %v", hotel.Name, err)
                errorMsg = err.Error()
            } else {
                successCount++
            }
        }
        status := "success"
        if errorMsg != "" && successCount > 0 {
            status = "partial_failure"
        } else if successCount == 0 {
            status = "failed"
            errorMsg = "No hotels scraped successfully"
        }
        now := time.Now()
        database.DB.Model(scrapingLog).Updates(map[string]interface{}{
            "status":        status,
            "hotels_count":  successCount,
            "error_message": errorMsg,
            "completed_at":  &now,
        })
        log.Printf("Updated %d hotels for %s", successCount, city)
        time.Sleep(5 * time.Second)
    }
}
func InitialScrape() {
    log.Println("Starting initial scrape...")
    cities := []string{"Delhi", "Mumbai", "Bangalore"}
    for _, city := range cities {
        hotels := ScrapeCity(city)
        for _, hotel := range hotels {
            if err := database.DB.Clauses(clause.OnConflict{
                Columns:   []clause.Column{{Name: "name"}, {Name: "city"}},
                DoNothing: true,
            }).Create(&hotel).Error; err != nil {
                log.Printf("Error inserting %s: %v", hotel.Name, err)
            }
        }    
        log.Printf("Initial scrape completed for %s: %d hotels", city, len(hotels))
        
        time.Sleep(5 * time.Second)
    }
    log.Println("Initial scrape finished!")
}
func upsertHotel(hotel models.Hotel) error {
    result := database.DB.Clauses(clause.OnConflict{
        Columns: []clause.Column{{Name: "name"}, {Name: "city"}},
        DoUpdates: clause.AssignmentColumns([]string{
            "price", "rating", "location", "image_url", "source", "updated_at",
        }),
    }).Create(&hotel)
    if result.Error != nil {
        return result.Error
    }
    priceHistory := models.PriceHistory{
        HotelID:   hotel.ID,
        Price:     hotel.Price,
        Timestamp: time.Now(),
    }
    return database.DB.Create(&priceHistory).Error
}
