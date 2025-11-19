package main

import (
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// Monitor represents a monitored target and its latest state.
type Monitor struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Type      string    `json:"type" gorm:"not null"`
    Status    string    `json:"status" gorm:"not null"`
    LastCheck time.Time `json:"last_check" gorm:"not null"`
}

// monitorCreateRequest captures required data for creating a monitor.
type monitorCreateRequest struct {
    Name   string `json:"name" binding:"required"`
    Type   string `json:"type" binding:"required"`
    Status string `json:"status" binding:"required"`
}

// monitorUpdateRequest captures fields that can be updated for a monitor.
type monitorUpdateRequest struct {
    Name   *string `json:"name"`
    Type   *string `json:"type"`
    Status *string `json:"status"`
}

func main() {
    _ = godotenv.Load()

    readKey := strings.TrimSpace(getEnv("READ_KEY"))
    adminKey := strings.TrimSpace(getEnv("ADMIN_KEY"))

    if readKey == "" || adminKey == "" {
        log.Fatal("READ_KEY and ADMIN_KEY must be provided via environment variables")
    }

    db, err := gorm.Open(sqlite.Open("monitors.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }

    if err := db.AutoMigrate(&Monitor{}); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    router := gin.Default()

    router.GET("/monitor", authorize(readKey, adminKey, true), func(c *gin.Context) {
        var monitors []Monitor
        if err := db.Order("id asc").Find(&monitors).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch monitors"})
            return
        }
        c.JSON(http.StatusOK, monitors)
    })

    router.POST("/monitor", authorize(readKey, adminKey, false), func(c *gin.Context) {
        var req monitorCreateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
            return
        }

        monitor := Monitor{
            Name:      req.Name,
            Type:      req.Type,
            Status:    req.Status,
            LastCheck: time.Now(),
        }

        if err := db.Create(&monitor).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create monitor"})
            return
        }

        c.JSON(http.StatusCreated, monitor)
    })

    router.PUT("/monitor/:id", authorize(readKey, adminKey, false), func(c *gin.Context) {
        var req monitorUpdateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
            return
        }

        var monitor Monitor
        if err := db.First(&monitor, c.Param("id")).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"message": "Monitor not found"})
            return
        }

        if req.Name != nil {
            monitor.Name = *req.Name
        }
        if req.Type != nil {
            monitor.Type = *req.Type
        }
        if req.Status != nil {
            monitor.Status = *req.Status
            monitor.LastCheck = time.Now()
        }

        if err := db.Save(&monitor).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update monitor"})
            return
        }

        c.JSON(http.StatusOK, monitor)
    })

    router.DELETE("/monitor/:id", authorize(readKey, adminKey, false), func(c *gin.Context) {
        if err := db.Delete(&Monitor{}, c.Param("id")).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete monitor"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Monitor deleted"})
    })

    router.GET("/status", authorize(readKey, adminKey, true), func(c *gin.Context) {
        var monitors []Monitor
        if err := db.Find(&monitors).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch status"})
            return
        }

        healthy := 0
        for _, m := range monitors {
            if strings.EqualFold(m.Status, "healthy") {
                healthy++
            }
        }

        status := "Degraded"
        if len(monitors) == 0 || healthy == len(monitors) {
            status = "OK"
        }

        c.JSON(http.StatusOK, gin.H{
            "status":           status,
            "monitors":         len(monitors),
            "healthy_monitors": healthy,
        })
    })

    if err := router.Run(); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}

// authorize returns middleware enforcing key-based access control.
func authorize(readKey, adminKey string, allowRead bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := strings.TrimSpace(c.GetHeader("Authorization"))
        if key == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
            c.Abort()
            return
        }

        if key == adminKey {
            c.Next()
            return
        }

        if allowRead && key == readKey {
            c.Next()
            return
        }

        c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
        c.Abort()
    }
}

// getEnv wraps lookup to simplify testing and defaults.
func getEnv(key string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return ""
}
