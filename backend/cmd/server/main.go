package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"the-illumiteers-station/backend/internal/handlers"
	"the-illumiteers-station/backend/internal/sheets"
	"the-illumiteers-station/backend/internal/youtube"
)

func main() {
	// Try loading .env from workspace locations
	loadEnvFromLocations()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}

	fmt.Printf("Configured PORT: %s\n", port)
	fmt.Printf("Configured ALLOWED_ORIGIN: %s\n", allowedOrigin)

	// Initialize Clients
	sheetsCli, err := sheets.NewClient()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize Google Sheets Client: %v\n", err)
	} else {
		fmt.Println("Google Sheets Client initialized successfully")
	}

	youtubeSrv, err := youtube.NewService()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize YouTube Service: %v\n", err)
	} else {
		fmt.Println("YouTube Service initialized successfully")
	}

	// Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware(allowedOrigin))

	// Setup Handlers
	h := handlers.NewHandler(youtubeSrv, sheetsCli)
	h.RegisterRoutes(r)

	// Health Check / Root Index
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"app": "The Illumiteers API", "status": "ok"}`))
	})

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Starting server on %s...\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("Server failed: %v\n", err)
		os.Exit(1)
	}
}

func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func loadEnvFromLocations() {
	// Check standard file paths relative to expected run CWDs
	paths := []string{
		".env",
		"../.env",
		"../../.env",
		"/home/david/Projects/lorcana-community-platform/.env",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			fmt.Printf("Loading environment from %s\n", p)
			loadEnvFile(p)
			return
		}
	}
	fmt.Println("No .env file found. Reading system environment variables instead.")
}

func loadEnvFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return
	}
	
	lines := strings.Split(string(bytes), "\n")
	var currentKey string
	var currentVal strings.Builder
	inMultiLine := false
	var quoteChar rune

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if inMultiLine {
			currentVal.WriteString("\n")
			// Look for the closing quote
			if strings.HasSuffix(trimmed, string(quoteChar)) {
				// Find index of quoteChar from end
				valPart := line
				lastIndex := strings.LastIndex(valPart, string(quoteChar))
				if lastIndex != -1 {
					currentVal.WriteString(valPart[:lastIndex])
				}
				
				value := currentVal.String()
				if os.Getenv(currentKey) == "" {
					os.Setenv(currentKey, value)
				}
				inMultiLine = false
				currentKey = ""
				currentVal.Reset()
			} else {
				currentVal.WriteString(line)
			}
			continue
		}

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Check if value starts with a quote and is multiline
		if len(val) > 0 && (val[0] == '"' || val[0] == '\'') {
			quoteChar = rune(val[0])
			// If it starts and ends with the same quote (and not empty quote pair or length 1)
			if len(val) > 1 && val[len(val)-1] == val[0] {
				actualVal := val[1 : len(val)-1]
				if os.Getenv(key) == "" {
					os.Setenv(key, actualVal)
				}
			} else {
				inMultiLine = true
				currentKey = key
				// Write the part of the value after the opening quote
				origVal := parts[1]
				firstIndex := strings.Index(origVal, string(quoteChar))
				if firstIndex != -1 && firstIndex < len(origVal)-1 {
					currentVal.WriteString(origVal[firstIndex+1:])
				}
			}
		} else {
			// Normal key-value pair
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
	}
}
