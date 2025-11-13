package database

import (
    "fmt"
    "log"
    "time"
    "github.com/Samriddha9619/Go_assignment1/backend/config"
    "github.com/Samriddha9619/Go_assignment1/backend/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) error {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
    )
    
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })

    if err != nil {
        return fmt.Errorf("error opening database: %w", err)
    }

    sqlDB, err := DB.DB()
    if err != nil {
        return fmt.Errorf("error getting database: %w", err)
    }

    sqlDB.SetMaxOpenConns(15)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    if err := sqlDB.Ping(); err != nil {
        return fmt.Errorf("error connecting to database: %w", err)
    }
    log.Println("Database connected successfully")

    if err := automigrate(); err != nil {
        return fmt.Errorf("error during automigration: %w", err)
    }
    return nil
}

func automigrate() error {
    err := DB.AutoMigrate(
        &models.User{},
        &models.Hotel{},
        &models.PriceHistory{},
        &models.ScrapingLog{},
    )

    if err != nil {
        return err
    }
    log.Println("Database migrated successfully")
    return nil
}

func Close() error {
    if DB != nil {
        sqlDB, err := DB.DB()
        if err != nil {
            return err
        }
        return sqlDB.Close()
    }
    return nil
}