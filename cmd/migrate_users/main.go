package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

// LegacyUser matches the structure of the old users.json
type LegacyUser struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"` // JSON key might have been "PasswordHash" or "password_hash" depending on old code
	BestScore    int       `json:"best_score"`
	TotalGames   int       `json:"total_games"`
	TotalWins    int       `json:"total_wins"`
	CreatedAt    time.Time `json:"created_at"`
}

func main() {
	// 1. Check if users.json exists
	if _, err := os.Stat("users.json"); os.IsNotExist(err) {
		log.Fatal("users.json not found in current directory. Please place your backup users.json here.")
	}

	// 2. Open SQLite DB
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite", "data/game.db")
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer db.Close()

	// Ensure table exists (in case the server hasn't run yet)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password_hash TEXT NOT NULL,
			best_score INTEGER DEFAULT 0,
			total_games INTEGER DEFAULT 0,
			total_wins INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// 3. Read and Parse JSON
	fileContent, err := os.ReadFile("users.json")
	if err != nil {
		log.Fatal("Failed to read users.json:", err)
	}

	var users map[string]LegacyUser
	// Try parsing as map first (common structure)
	if err := json.Unmarshal(fileContent, &users); err != nil {
		// If map failed, try array
		var userList []LegacyUser
		if err2 := json.Unmarshal(fileContent, &userList); err2 != nil {
			log.Fatalf("Failed to parse users.json: %v", err)
		}
		// Convert list to map for easier handling
		users = make(map[string]LegacyUser)
		for _, u := range userList {
			users[u.Username] = u
		}
	}

	// 4. Migrate
	log.Printf("Found %d users to migrate...", len(users))
	count := 0
	for _, u := range users {
		// Use INSERT OR IGNORE to skip duplicates
		_, err := db.Exec(
			`INSERT OR IGNORE INTO users (username, password_hash, best_score, total_games, total_wins, created_at) 
			 VALUES (?, ?, ?, ?, ?, ?)`,
			u.Username, u.PasswordHash, u.BestScore, u.TotalGames, u.TotalWins, u.CreatedAt,
		)
		if err != nil {
			log.Printf("Error migrating user %s: %v\n", u.Username, err)
		} else {
			count++
		}
	}

	fmt.Printf("✅ Migration complete! Successfully imported %d users into data/game.db\n", count)
}
