package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"motico-api/config"
)

func main() {
	fmt.Println("🔍 Diagnosing Supabase connection...")
	fmt.Println()

	cfg, err := config.Load("config/config.json")
	if err != nil {
		log.Fatalf("❌ Error loading config: %v", err)
	}

	fmt.Println("📊 Configuration loaded:")
	fmt.Printf("   Host: %s\n", cfg.Database.Host)
	fmt.Printf("   Port: %s\n", cfg.Database.Port)
	fmt.Printf("   User: %s\n", cfg.Database.User)
	fmt.Printf("   Database: %s\n", cfg.Database.Name)
	fmt.Printf("   SSL Mode: %s\n", cfg.Database.SSLMode)
	if cfg.Database.PoolMode != "" {
		fmt.Printf("   Pool Mode: %s\n", cfg.Database.PoolMode)
	}

	if cfg.Database.Password == "" {
		fmt.Println("   Password: ❌ NOT SET")
		log.Fatal("❌ DB_PASSWORD environment variable is not set")
	} else {
		fmt.Println("   Password: ✅ SET")
	}
	fmt.Println()

	// Test DNS resolution
	fmt.Println("🔍 Testing DNS resolution...")
	addresses, err := net.LookupHost(cfg.Database.Host)
	if err != nil {
		log.Fatalf("❌ DNS lookup failed: %v", err)
	}
	fmt.Printf("✅ DNS resolved to:\n")
	for _, addr := range addresses {
		fmt.Printf("   - %s\n", addr)
	}
	fmt.Println()

	// Test TCP connection (IPv4 only)
	fmt.Println("🔍 Testing TCP connection (IPv4)...")
	ipv4Addr := ""
	for _, addr := range addresses {
		ip := net.ParseIP(addr)
		if ip != nil && ip.To4() != nil {
			ipv4Addr = addr
			break
		}
	}

	if ipv4Addr == "" {
		fmt.Println("⚠️  Warning: No IPv4 address found, only IPv6")
		fmt.Println("   This might cause connection issues")
		fmt.Println()
	} else {
		fmt.Printf("   Using IPv4 address: %s\n", ipv4Addr)
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ipv4Addr, cfg.Database.Port), 5*time.Second)
		if err != nil {
			fmt.Printf("❌ TCP connection failed: %v\n", err)
			fmt.Println()
			fmt.Println("💡 Possible solutions:")
			fmt.Println("   1. Check if your Supabase project is active")
			fmt.Println("   2. Check firewall/network settings")
			fmt.Println("   3. Try using the connection pooler endpoint")
			fmt.Println("   4. Verify your network allows outbound connections to port 5432")
		} else {
			if err := conn.Close(); err != nil {
				fmt.Printf("⚠️  Warning: Error closing connection: %v\n", err)
			}
			fmt.Println("✅ TCP connection successful!")
		}
		fmt.Println()
	}

	// Check environment variables
	fmt.Println("🔍 Checking environment variables...")
	envVars := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
		"DB_POOL_MODE",
		"JWT_SECRET_KEY",
	}

	allSet := true
	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		if value == "" {
			fmt.Printf("   ❌ %s: NOT SET\n", envVar)
			allSet = false
		} else {
			if envVar == "DB_PASSWORD" || envVar == "JWT_SECRET_KEY" {
				fmt.Printf("   ✅ %s: SET (hidden)\n", envVar)
			} else {
				fmt.Printf("   ✅ %s: %s\n", envVar, value)
			}
		}
	}
	fmt.Println()

	if !allSet {
		fmt.Println("⚠️  Some environment variables are missing")
		fmt.Println("   Make sure your .env file exists and contains all required variables")
		fmt.Println()
	}

	// Connection string preview (without password)
	fmt.Println("📋 Connection string preview:")
	queryParams := fmt.Sprintf("sslmode=%s", cfg.Database.SSLMode)
	if cfg.Database.PoolMode != "" {
		queryParams += fmt.Sprintf("&pool_mode=%s", cfg.Database.PoolMode)
	}
	connString := fmt.Sprintf(
		"postgres://%s:***@%s:%s/%s?%s",
		cfg.Database.User,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		queryParams,
	)
	fmt.Printf("   %s\n", connString)
	fmt.Println()

	fmt.Println("💡 Next steps:")
	fmt.Println("   1. Verify your Supabase project is active (not paused)")
	fmt.Println("   2. Check the connection string in Supabase Dashboard")
	fmt.Println("   3. Try using the connection pooler (port 6543) if available")
	fmt.Println("   4. Check your network/firewall settings")
}
