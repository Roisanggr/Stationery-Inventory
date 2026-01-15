package main

import (
    "fmt"
    "log"
    "database/sql"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"

    "ATK-Backend/routes"
    "ATK-Backend/models"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    _ "github.com/go-sql-driver/mysql"
)

// responseWriter wrapper to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
    return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// Logging middleware - Express.js style
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        wrapped := newResponseWriter(w)
        
        // Call next handler
        next.ServeHTTP(wrapped, r)
        
        // Log after request is handled
        duration := time.Since(start)
        statusCode := wrapped.statusCode
        
        // Determine status text
        statusText := "success"
        if statusCode >= 400 && statusCode < 500 {
            statusText = "client error"
        } else if statusCode >= 500 {
            statusText = "server error"
        }
        
        // Color codes for terminal output (optional)
        statusColor := "\033[32m" // green
        if statusCode >= 400 && statusCode < 500 {
            statusColor = "\033[33m" // yellow
        } else if statusCode >= 500 {
            statusColor = "\033[31m" // red
        }
        resetColor := "\033[0m"
        
        // Log format: METHOD /path STATUS status_text [duration]
        log.Printf("%s %s %s%d%s %s [%s]",
            r.Method,
            r.RequestURI,
            statusColor,
            statusCode,
            resetColor,
            statusText,
            duration.Round(time.Millisecond),
        )
    })
}

func main() {
    // load .env if present
    _ = godotenv.Load()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // init database if DATABASE env present
    dsn := strings.TrimSpace(os.Getenv("DATABASE"))
    // remove surrounding quotes if present
    if len(dsn) >= 2 {
        if (dsn[0] == '"' && dsn[len(dsn)-1] == '"') || (dsn[0] == '\'' && dsn[len(dsn)-1] == '\'') {
            dsn = dsn[1 : len(dsn)-1]
        }
    }

    if dsn != "" {
        // allow both mysql DSN or mysql:// URL; convert if necessary
        if len(dsn) >= 8 && dsn[:8] == "mysql://" {
            // parse URL -> user:pass@tcp(host:port)/dbname?params
            // simple conversion: strip prefix and replace first / with @tcp(
            // safer: use net/url
            // implement basic conversion
            // example: mysql://user:pass@host:3306/dbname?params
            // -> user:pass@tcp(host:3306)/dbname?params
            u, err := url.Parse(dsn)
            if err == nil {
                user := u.User.Username()
                pass, _ := u.User.Password()
                host := u.Host
                path := u.Path
                if len(path) > 0 && path[0] == '/' {
                    path = path[1:]
                }
                q := u.RawQuery
                dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, path)
                if q != "" {
                    dsn = dsn + "?" + q
                }
            }
        }

        db, err := sql.Open("mysql", dsn)
        if err != nil {
            log.Printf("failed to open db (DSN sanitized): %v\n", err)
        } else {
            if err := db.Ping(); err != nil {
                log.Printf("failed to ping db (check credentials/host/port): %v\n", err)
            } else {
                if err := models.InitDB(db); err != nil {
                    log.Printf("failed to init models db (creating mst_atk): %v\n", err)
                } else {
                    log.Printf("Connected to DB and initialized table mst_atk")
                }
            }
        }
    }

    r := mux.NewRouter()
    
    // Apply logging middleware to all routes
    r.Use(loggingMiddleware)
    
    routes.RegisterRoutes(r)

    fmt.Printf("Server berjalan di port %s...\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}