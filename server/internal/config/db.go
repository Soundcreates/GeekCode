package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"geekCode/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	// Check for Render's DATABASE_URL first (most common on Render)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Fallback to POSTGRES_URL
		databaseURL = os.Getenv("POSTGRES_URL")
	}
	
	// Debug logging
	log.Printf("Raw DATABASE_URL: %s", maskPassword(databaseURL))
	
	var dsn string

	if databaseURL != "" {
		// Try to parse the database URL
		parsedURL, err := url.Parse(databaseURL)
		if err != nil {
			log.Printf("Warning: Invalid DATABASE_URL format, falling back to individual env vars: %v", err)
		} else {
			// Check if the URL is malformed (host is localhost but we have a different hostname in the string)
			if parsedURL.Hostname() == "localhost" && strings.Contains(databaseURL, "dpg-") {
				log.Printf("Detected malformed URL, attempting to fix...")
				// This is likely a malformed URL where the hostname got mixed up
				// Let's try to extract the correct components manually
				dsn = fixMalformedURL(databaseURL)
			} else {
				// Normal URL parsing
				// Fix Render's format: "postgresql://" → "postgres://"
				if parsedURL.Scheme == "postgresql" {
					parsedURL.Scheme = "postgres"
				}

				// Ensure SSL mode is enforced for production
				query := parsedURL.Query()
				if query.Get("sslmode") == "" {
					query.Set("sslmode", "require")
				}
				parsedURL.RawQuery = query.Encode()

				dsn = parsedURL.String()
			}
			log.Printf("Using database URL from environment variable")
		}
	}

	// If we don't have a valid DSN from URL parsing, try individual env vars
	if dsn == "" {
		// Try to construct DSN from individual environment variables
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")

		// For production, try to extract from DATABASE_URL if individual vars are not set
		if host == "" && databaseURL != "" {
			if parsedURL, err := url.Parse(databaseURL); err == nil {
				host = parsedURL.Hostname()
				if parsedURL.Port() != "" {
					port = parsedURL.Port()
				}
				user = parsedURL.User.Username()
				password, _ = parsedURL.User.Password()
				name = strings.TrimPrefix(parsedURL.Path, "/")
			}
		}

		// Set defaults for local development
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "5432"
		}
		if user == "" {
			user = "postgres"
		}
		if name == "" {
			name = "geekcode"
		}

		if password == "" {
			log.Fatal("Database password is required. Please set DB_PASSWORD environment variable or provide a valid DATABASE_URL")
		}

		// Determine SSL mode based on environment
		sslMode := "disable"
		if os.Getenv("PORT") != "" { // Production environment
			sslMode = "require"
		}

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, name, sslMode,
		)
		log.Printf("Using individual database environment variables")
		log.Printf("Database host: %s, port: %s, user: %s, dbname: %s", host, port, user, name)
	}

	log.Printf("Attempting to connect to database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Database connection string (without password): %s", maskPassword(dsn))
		log.Fatal("Failed to connect to database: ", err)
	}
	
	log.Printf("Successfully connected to database!")

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Room{},
	); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	return db, nil
}

// maskPassword masks the password in the database connection string for logging
func maskPassword(dsn string) string {
	// Replace password=xxx with password=***
	re := regexp.MustCompile(`password=[^&\s]+`)
	return re.ReplaceAllString(dsn, "password=***")
}

// fixMalformedURL attempts to fix malformed database URLs
func fixMalformedURL(databaseURL string) string {
	// Example malformed URL: postgres://localhost:password@dpg-hostname/dbname
	// We need to extract: user, password, hostname, port, dbname
	
	// Use regex to extract components from malformed URL
	// Pattern: postgres://user:password@hostname:port/dbname
	re := regexp.MustCompile(`postgres://([^:]+):([^@]+)@([^:]+):?(\d*)/([^?]+)`)
	matches := re.FindStringSubmatch(databaseURL)
	
	if len(matches) >= 6 {
		user := matches[1]
		password := matches[2]
		hostname := matches[3]
		port := matches[4]
		dbname := matches[5]
		
		// Set default port if not provided
		if port == "" {
			port = "5432"
		}
		
		// Construct proper DSN
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			hostname, port, user, password, dbname)
		
		log.Printf("Fixed malformed URL - host: %s, port: %s, user: %s, dbname: %s", 
			hostname, port, user, dbname)
		
		return dsn
	}
	
	// If regex doesn't match, try a simpler approach
	// Look for the pattern: postgres://user:password@hostname/dbname
	re2 := regexp.MustCompile(`postgres://([^:]+):([^@]+)@([^/]+)/([^?]+)`)
	matches2 := re2.FindStringSubmatch(databaseURL)
	
	if len(matches2) >= 5 {
		user := matches2[1]
		password := matches2[2]
		hostname := matches2[3]
		dbname := matches2[4]
		
		dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=require",
			hostname, user, password, dbname)
		
		log.Printf("Fixed malformed URL (simple) - host: %s, user: %s, dbname: %s", 
			hostname, user, dbname)
		
		return dsn
	}
	
	// Try one more pattern for the specific case we're seeing
	// postgres://localhost:password@dpg-hostname/dbname
	re3 := regexp.MustCompile(`postgres://localhost:([^@]+)@([^/]+)/([^?]+)`)
	matches3 := re3.FindStringSubmatch(databaseURL)
	
	if len(matches3) >= 4 {
		password := matches3[1]
		hostname := matches3[2]
		dbname := matches3[3]
		
		// Extract username from the hostname if it contains user info
		// This handles cases where the URL is really malformed
		user := "postgres" // default user
		if strings.Contains(hostname, ":") {
			parts := strings.Split(hostname, ":")
			if len(parts) > 1 {
				user = parts[0]
				hostname = parts[1]
			}
		}
		
		dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=require",
			hostname, user, password, dbname)
		
		log.Printf("Fixed malformed URL (localhost pattern) - host: %s, user: %s, dbname: %s", 
			hostname, user, dbname)
		
		return dsn
	}
	
	log.Printf("Could not fix malformed URL, returning empty string")
	return ""
}
