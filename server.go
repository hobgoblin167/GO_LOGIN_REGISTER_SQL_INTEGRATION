package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// Глобальная переменная для подключения к БД
var conn *pgx.Conn

// для верстки функция
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "login.html")
		return
	}
	if r.URL.Path == "/registerweb" {
		http.ServeFile(w, r, "register.html")
		return
	}
	if r.URL.Path == "/loginweb" {
		http.ServeFile(w, r, "login.html")
		return
	}
	http.NotFound(w, r)
}

// функция для обработки логина
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var passInput string
		err := conn.QueryRow(context.Background(),
			"SELECT password FROM userBase WHERE name = $1",
			username).Scan(&passInput)

		if err != nil {
			fmt.Fprintf(w, "Ошибка! Не существует такого аккаунта!")
			fmt.Println("Логин: пользователь не найден:", username)
			return
		}

		if passInput == password {
			fmt.Fprintf(w, "Вы вошли в аккаунт %s!", username)
			fmt.Println("Логин: успешный вход для", username)
		} else {
			fmt.Fprintf(w, "Неверный пароль!")
			fmt.Println("Логин: неверный пароль для", username)
		}
	} else {
		// Если не POST, показываем форму
		http.ServeFile(w, r, "login.html")
	}
}

// функция для обработки регистрации
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var existingPassword string
		err := conn.QueryRow(context.Background(),
			"SELECT password FROM userBase WHERE name = $1",
			username).Scan(&existingPassword)

		if err == nil {
			fmt.Fprintf(w, "Ошибка! Аккаунт '%s' уже существует!", username)
			fmt.Println("Регистрация: пользователь уже существует:", username)
			return
		}

		// Создаем нового пользователя
		_, err = conn.Exec(context.Background(),
			"INSERT INTO userBase (name, password) VALUES ($1, $2)",
			username, password)

		if err != nil {
			fmt.Fprintf(w, "Ошибка создания аккаунта: %v", err)
			fmt.Println("Регистрация: ошибка БД:", err)
			return
		}

		fmt.Fprintf(w, "Аккаунт '%s' успешно создан!", username)
		fmt.Println("Регистрация: создан пользователь:", username)

	} else {
		// Если не POST, показываем форму
		http.ServeFile(w, r, "register.html")
	}
}

func main() {
	// 1. Сначала подключаемся к БД
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dburl := os.Getenv("DATABASE_URL")

	conn, err = pgx.Connect(context.Background(),
		dburl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка подключения к базе: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	fmt.Println("Подключились к базе данных")

	// Создаем таблицу userBase если её нет
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS userBase (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания таблицы: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Таблица userBase создана/проверена")

	// админ
	_, err = conn.Exec(context.Background(), `
		INSERT INTO userBase (name, password) 
		VALUES ('admin', 'admin1')
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		fmt.Printf("Ошибка добавления тестового пользователя: %v\n", err)
	} else {
		fmt.Println("Добавлен аккаунт админа")
	}

	// Настраиваем сервер
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", homeHandler)
	serverMux.HandleFunc("/login", loginHandler)
	serverMux.HandleFunc("/register", registerHandler)

	fmt.Println("Сервер запущен: http://localhost:8080")
	fmt.Println("Логин: /login")
	fmt.Println("Регистрация: /register")

	// Запускаем сервер
	err = http.ListenAndServe(":8080", serverMux)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка сервера: %v\n", err)
		os.Exit(1)
	}
}
