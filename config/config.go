package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB        *gorm.DB
	JwtSecret string
	dbMutex   sync.Mutex
	isConnected bool
)

// Load konfigurasi dari file .env
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Gagal memuat file .env, menggunakan environment variables sistem")
	} else {
		log.Println("File .env berhasil dimuat!")
	}
}

// Connect menghubungkan aplikasi dengan database MySQL dan memuat JWT_SECRET
func Connect() {
	// Memuat .env hanya sekali
	LoadEnv()

	// Ambil environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	// Validasi environment variables
	if dbHost == "" || dbPort == "" || dbDatabase == "" || dbUsername == "" || dbPassword == "" {
		log.Fatal("Environment variables database tidak lengkap. Periksa DB_HOST, DB_PORT, DB_DATABASE, DB_USERNAME, DB_PASSWORD")
	}

	// Escape password untuk menangani karakter khusus
	escapedPassword := url.QueryEscape(dbPassword)

	// Buat DSN untuk MySQL dengan parameter tambahan
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s&readTimeout=30s&writeTimeout=30s&allowNativePasswords=true&allowOldPasswords=true&tls=false",
		dbUsername,
		escapedPassword,
		dbHost,
		dbPort,
		dbDatabase,
	)

	// Log DSN tanpa password untuk debugging (opsional)
	safeDsn := fmt.Sprintf("%s:***@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUsername, dbHost, dbPort, dbDatabase)
	log.Printf("Mencoba koneksi ke: %s", safeDsn)

	// Konfigurasi GORM dengan logger dan pengaturan koneksi
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Coba koneksi ke database
	d, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		log.Printf("Gagal terhubung ke database: %v", err)
		log.Println("Kemungkinan penyebab:")
		log.Println("1. Password mengandung karakter khusus")
		log.Println("2. IP lokal tidak diizinkan oleh hosting")
		log.Println("3. User tidak memiliki permission yang tepat")
		log.Println("4. Database server tidak dapat diakses")
		log.Fatal("Koneksi database gagal - aplikasi tidak dapat berjalan tanpa database")
	}

	// Konfigurasi connection pool
	sqlDB, err := d.DB()
	if err != nil {
		log.Fatal("Gagal mendapatkan database instance:", err)
	}

	// Pengaturan connection pool
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test koneksi
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Gagal ping database:", err)
	}

	// Set DB global dengan mutex untuk thread safety
	dbMutex.Lock()
	DB = d
	isConnected = true
	dbMutex.Unlock()
	
	log.Println("Database MySQL berhasil terhubung!")

	// Muat JWT_SECRET
	LoadJwtSecret()
}

// GetDB returns the database connection with nil check
func GetDB() (*gorm.DB, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is nil, please ensure Connect() was called successfully")
	}
	return DB, nil
}

// IsConnected returns true if the database is connected
func IsConnected() bool {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	return isConnected
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

// TestConnection untuk testing koneksi database
func TestConnection() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	return sqlDB.Ping()
}
