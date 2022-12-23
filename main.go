// @title Albums API
// @version 1.0
// @description Swagger API for Golang Course Project.
// @termsOfService http://swagger.io/terms/

// @contact.name eltsova_ad
// @contact.email eltsova.ad@gmail.com

// @license.name MIT

// Album model info
// @Description Info about albums
// @Description with its title, artist and my review
package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"lab1/controller/backend"

	"log"

	"github.com/gin-gonic/gin"
)

// @BasePath docs/v1
func main() {
	db, err := gorm.Open(sqlite.Open("AlbumsTop.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	err = backend.SetupDatabase(db)
	if err != nil {
		log.Fatalf("Database setup error: %s", err)
	}
	r := gin.Default()
	backend.SetupRouter(r, db)
	err = r.Run("localhost:8080")
	if err != nil {
		log.Fatalf("gin Run error: %s", err)
	}
}
