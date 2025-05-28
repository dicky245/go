package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
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
	log.Println("File .env berhasil dimuat!")
}

// Connect menghubungkan aplikasi dengan database PostgreSQL dan memuat JWT_SECRET
func Connect() {
	// Memuat .env hanya sekali
	LoadEnv()

	// Koneksi database PostgreSQL
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USERNAME") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_DATABASE") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=require TimeZone=Asia/Jakarta" // Gunakan sslmode sesuai kebutuhan

	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}
	DB = d
	log.Println("Database PostgreSQL berhasil terhubung!")

	// Muat JWT_SECRET
	LoadJwtSecret()
}

// LoadJwtSecret memuat nilai JWT_SECRET dari environment variables
func LoadJwtSecret() {
	JwtSecret = os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		log.Fatal("JWT_SECRET tidak ditemukan dalam environment variables. Pastikan ada di file .env.")
	}
	log.Println("JWT_SECRET berhasil dimuat!")
}

// GetJwtSecret mengembalikan nilai JWT_SECRET
func GetJwtSecret() string {
	return JwtSecret
}
