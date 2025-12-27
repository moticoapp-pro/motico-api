package main

import (
	"context"
	"fmt"
	"motico-api/config"
	"motico-api/internal/repository"
	"os"
)

func main() {
	fmt.Println("🔌 Testing Supabase connection...")
	fmt.Println()

	cfg, err := config.Load("config/config.json")
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cfg.Database.Password == "" {
		fmt.Println("❌ DB_PASSWORD environment variable is not set. Please create a .env file with your Supabase credentials.")
		os.Exit(1)
	}

	if cfg.Database.Host == "localhost" {
		fmt.Println("⚠️  Warning: DB_HOST is set to 'localhost'. Make sure you're using your Supabase host.")
	}

	fmt.Printf("📊 Connection details:\n")
	fmt.Printf("   Host: %s\n", cfg.Database.Host)
	fmt.Printf("   Port: %s\n", cfg.Database.Port)
	fmt.Printf("   User: %s\n", cfg.Database.User)
	fmt.Printf("   Database: %s\n", cfg.Database.Name)
	fmt.Printf("   SSL Mode: %s\n", cfg.Database.SSLMode)
	fmt.Println()

	ctx := context.Background()

	fmt.Println("🔄 Creating connection pool...")
	pool, err := repository.NewConnectionPool(ctx, cfg)
	if err != nil {
		fmt.Printf("❌ Error creating connection pool: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Connection pool created successfully!")
	fmt.Println()

	fmt.Println("🔄 Testing database connection...")
	var result int
	err = pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		pool.Close()
		fmt.Printf("❌ Error executing test query: %v\n", err)
		os.Exit(1)
	}

	if result != 1 {
		pool.Close()
		fmt.Printf("❌ Unexpected query result: expected 1, got %d\n", result)
		os.Exit(1)
	}

	fmt.Println("✅ Database connection test successful!")
	fmt.Println()

	fmt.Println("🔄 Testing PostgreSQL version...")
	var version string
	err = pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		pool.Close()
		fmt.Printf("❌ Error getting PostgreSQL version: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ PostgreSQL version: %s\n", version)
	fmt.Println()

	fmt.Println("🎉 All connection tests passed! Your Supabase connection is working correctly.")
	pool.Close()
}
