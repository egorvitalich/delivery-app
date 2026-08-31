package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

func main() {
    // Загружаем .env файл
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ .env файл не найден, используем переменные окружения")
    }

    // Берем настройки из .env
    connStr := "user=" + os.Getenv("DB_USER") +
        " password=" + os.Getenv("DB_PASSWORD") +
        " dbname=" + os.Getenv("DB_NAME") +
        " host=" + os.Getenv("DB_HOST") +
        " port=" + os.Getenv("DB_PORT") +
        " sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Ошибка подключения:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("БД не отвечает:", err)
    }
    log.Println("✅ Подключено к БД")

    http.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
        var version string
        err := db.QueryRow("SELECT version FROM application_info LIMIT 1").Scan(&version)
        if err != nil {
            http.Error(w, `{"error":"Версия не найдена"}`, http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"version": version})
    })

    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("pong"))
    })

    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🚀 Сервер запущен на http://localhost:%s", port)
    log.Println("  GET /api/version - получить версию")
    log.Println("  GET /ping - проверка")

    log.Fatal(http.ListenAndServe(":"+port, nil))
}
