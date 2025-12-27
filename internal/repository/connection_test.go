package repository

import (
	"context"
	"motico-api/config"
	"testing"
)

func TestConnectionPool(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	cfg, err := config.Load("../../config/config.json")
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	ctx := context.Background()
	pool, err := NewConnectionPool(ctx, cfg)
	if err != nil {
		t.Fatalf("Error creating connection pool: %v", err)
	}
	defer pool.Close()

	var result int
	err = pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Error executing query: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

