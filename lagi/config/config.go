package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	JwtSecret string
)

// Load konfigurasi dari file .env
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env. Pastikan file ada di direktori yang benar.")
	}
}

// Koneksi ke database
func Connect() {
	// Pastikan .env sudah dimuat sebelum mengambil variabel
	LoadEnv()

	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") +
		")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(" Gagal terhubung ke database:", err)
	}

	DB = d
	log.Println("Database berhasil terhubung!")
}

func LoadJwtSecret() {
	JwtSecret = os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		log.Fatal(" JWT_SECRET tidak ditemukan dalam environment variables. Pastikan ada di file .env.")
	}
	log.Println(" JWT_SECRET berhasil dimuat!")
}

func GetJwtSecret() string {
	return JwtSecret
}
