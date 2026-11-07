package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"storage-manager-api/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah instance global koneksi database GORM (opsional diakses selain via dependency injection).
var DB *gorm.DB

// BuildDSN membentuk Data Source Name (DSN) MySQL dari konfigurasi environment.
func BuildDSN(cfg *Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
}

// ensureDatabaseExists membuat database otomatis jika belum tersedia di server MySQL.
func ensureDatabaseExists(cfg *Config) error {
	rootDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
	)

	sqlDB, err := sql.Open("mysql", rootDSN)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.DBName)
	_, err = sqlDB.Exec(query)
	return err
}

// ConnectDatabase membuka koneksi ke MySQL menggunakan GORM dan menjalankan AutoMigrate.
func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	if err := ensureDatabaseExists(cfg); err != nil {
		log.Printf("Peringatan: gagal memastikan database '%s' tersedia secara otomatis: %v", cfg.DBName, err)
	}

	dsn := BuildDSN(cfg)

	logLevel := logger.Info
	if cfg.GinMode == "release" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan instance sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("gagal melakukan ping ke database MySQL: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.File{}); err != nil {
		return nil, fmt.Errorf("gagal menjalankan migrasi skema database: %w", err)
	}

	DB = db
	log.Println("Database MySQL berhasil terhubung dan skema (users, files) telah dimigrasi.")
	return db, nil
}

