package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	analyticsDB *sql.DB
	analyticsOnce sync.Once
	analyticsMutex sync.RWMutex
)

// InitAnalyticsDB initializes the SQLite database for analytics tracking
func InitAnalyticsDB() error {
	var initErr error
	analyticsOnce.Do(func() {
		// Create data directory if it doesn't exist
		dataDir := "data"
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			initErr = fmt.Errorf("failed to create data directory: %w", err)
			return
		}

		// Open SQLite database
		dbPath := filepath.Join(dataDir, "analytics.db")
		db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=1")
		if err != nil {
			initErr = fmt.Errorf("failed to open analytics database: %w", err)
			return
		}

		// Test connection
		if err := db.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping analytics database: %w", err)
			return
		}

		// Create tables
		if err := createAnalyticsTables(db); err != nil {
			initErr = fmt.Errorf("failed to create analytics tables: %w", err)
			return
		}

		analyticsDB = db
		log.Println("Analytics database initialized successfully")
	})

	return initErr
}

// createAnalyticsTables creates the necessary tables for analytics
func createAnalyticsTables(db *sql.DB) error {
	// Create stats table (simpler approach - single row with counters)
	statsTable := `
	CREATE TABLE IF NOT EXISTS stats (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		total_generations INTEGER DEFAULT 0,
		total_downloads INTEGER DEFAULT 0,
		last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(statsTable); err != nil {
		return fmt.Errorf("failed to create stats table: %w", err)
	}

	// Insert initial row if it doesn't exist
	insertInitial := `
	INSERT OR IGNORE INTO stats (id, total_generations, total_downloads, last_updated)
	VALUES (1, 0, 0, CURRENT_TIMESTAMP);
	`

	if _, err := db.Exec(insertInitial); err != nil {
		return fmt.Errorf("failed to insert initial stats row: %w", err)
	}

	return nil
}

// Stats represents the analytics statistics
type Stats struct {
	TotalGenerations int64 `json:"total_generations"`
	TotalDownloads   int64 `json:"total_downloads"`
}

// GetStats retrieves the current analytics statistics
func GetStats() (*Stats, error) {
	analyticsMutex.RLock()
	defer analyticsMutex.RUnlock()

	if analyticsDB == nil {
		// Return zero stats if database not initialized
		return &Stats{TotalGenerations: 0, TotalDownloads: 0}, nil
	}

	var stats Stats
	query := `SELECT total_generations, total_downloads FROM stats WHERE id = 1`
	err := analyticsDB.QueryRow(query).Scan(&stats.TotalGenerations, &stats.TotalDownloads)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Stats{TotalGenerations: 0, TotalDownloads: 0}, nil
		}
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &stats, nil
}

// IncrementGeneration increments the generation counter
func IncrementGeneration() error {
	analyticsMutex.Lock()
	defer analyticsMutex.Unlock()

	if analyticsDB == nil {
		// Silently fail if database not initialized
		return nil
	}

	query := `
	UPDATE stats 
	SET total_generations = total_generations + 1,
	    last_updated = CURRENT_TIMESTAMP
	WHERE id = 1
	`

	_, err := analyticsDB.Exec(query)
	if err != nil {
		log.Printf("Failed to increment generation counter: %v", err)
		return fmt.Errorf("failed to increment generation: %w", err)
	}

	return nil
}

// IncrementDownload increments the download counter
func IncrementDownload() error {
	analyticsMutex.Lock()
	defer analyticsMutex.Unlock()

	if analyticsDB == nil {
		// Silently fail if database not initialized
		return nil
	}

	query := `
	UPDATE stats 
	SET total_downloads = total_downloads + 1,
	    last_updated = CURRENT_TIMESTAMP
	WHERE id = 1
	`

	_, err := analyticsDB.Exec(query)
	if err != nil {
		log.Printf("Failed to increment download counter: %v", err)
		return fmt.Errorf("failed to increment download: %w", err)
	}

	return nil
}

// CloseAnalyticsDB closes the analytics database connection
func CloseAnalyticsDB() error {
	analyticsMutex.Lock()
	defer analyticsMutex.Unlock()

	if analyticsDB != nil {
		return analyticsDB.Close()
	}
	return nil
}

