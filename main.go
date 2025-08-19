package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

const connStr = "host=host port=5432 user=pg_user password=pg_password dbname=db_name sslmode=disable"

var (
	db   *sql.DB
	once sync.Once
)

type MockResponse struct {
	Path               string
	Method             string
	ResponseBody       string
	ResponseStatusCode int
	Headers            string
}

func initDB() error {
	var err error
	once.Do(func() {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return
		}

		if err = db.Ping(); err != nil {
			db.Close()
			return
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(15 * time.Minute)
		db.SetConnMaxIdleTime(3 * time.Minute)

		fmt.Println("Database connection pool initialized")
	})
	return err
}

func getMockResponse(path string, method string) (*MockResponse, error) {
	var mockResp MockResponse
	query := "SELECT path, method, response_body, headers, response_status_code FROM return.mock_responses WHERE path = $1 AND method = $2"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, query, path, method)
	err := row.Scan(&mockResp.Path, &mockResp.Method, &mockResp.ResponseBody, &mockResp.Headers, &mockResp.ResponseStatusCode)

	if err != nil {
		return nil, err
	}

	return &mockResp, nil
}

func parseHeaders(headerStr string) map[string]string {
	headers := make(map[string]string)
	if headerStr == "" {
		return headers
	}

	pairs := strings.Split(headerStr, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			headers[key] = value
		}
	}

	return headers
}

func buildFullPath(r *http.Request) string {
	urlPath := r.URL.Path
	if r.URL.RawQuery != "" {
		urlPath += "?" + r.URL.RawQuery
	}
	return urlPath
}

func proxyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	urlPath := buildFullPath(r)
	method := r.Method
	fmt.Printf("Method: %s, Path: %s\n", method, urlPath)

	mockResp, err := getMockResponse(urlPath, method)
	if err != nil {
		if err == sql.ErrNoRows {
			// No mock response found, return 404
			http.NotFound(w, r)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	headers := parseHeaders(mockResp.Headers)
	for key, value := range headers {
		w.Header().Set(key, value)
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	statusCode := mockResp.ResponseStatusCode
	if statusCode == 0 {
		statusCode = http.StatusOK // Default to 200 if not specified
	}
	w.WriteHeader(statusCode)
	w.Write([]byte(mockResp.ResponseBody))
}

func registerHandlers(router *httprouter.Router, path string, handler httprouter.Handle) {
	router.GET(path, handler)
	router.POST(path, handler)
	router.PUT(path, handler)
	router.DELETE(path, handler)
	router.PATCH(path, handler)
	router.OPTIONS(path, handler)
	router.HEAD(path, handler)
}

func main() {
	if err := initDB(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer db.Close()

	router := httprouter.New()
	registerHandlers(router, "/*path", proxyHandler)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
