package main

import (
	"api-elektrodukasi/database"
	"api-elektrodukasi/handlers"
	"api-elektrodukasi/repositories"
	"api-elektrodukasi/services"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func loadConfig() *Config {
	// 1️⃣ load .env ke environment (DEV only)
	_ = godotenv.Load()

	// 2️⃣ viper baca dari ENV
	viper.AutomaticEnv()

	cfg := &Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// 3️⃣ validasi
	if cfg.Port == "" {
		log.Fatal("PORT is required")
	}
	if cfg.DBConn == "" {
		log.Fatal("DB_CONN is required")
	}

	return cfg
}

func main() {
	cfg := loadConfig()

	addr := "0.0.0.0:" + cfg.Port

	// Database connection
	db, err := database.InitDB(cfg.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize Database: ", err)
	}
	defer db.Close()

	// Dependency Injection
	CategoryRepo := repositories.NewCategoryRepository(db)
	CategoryService := services.NewCategoryService(CategoryRepo)
	CategoryHandler := handlers.NewCategoryHandler(CategoryService)

	// Route Handler
	http.HandleFunc("/categories", CategoryHandler.HandleCategory)
	http.HandleFunc("/categories/", CategoryHandler.HandleCategoryByID)

	// some text
	fmt.Println("Listening on ", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Server Run Failed: ", err)
	}
}
