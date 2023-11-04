package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var (
	hostIP string
	cache  sync.Map
)

func main() {
	hostIP := "localhost"
	// hostIP := "db-container"

	// Подключение к базе данных
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     hostIP,
		Port:     5432,
		Database: "L0_db",
		User:     "some_user",
		Password: "zxczxc",
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	// Подписка на сообщения NATS в фоновом режиме
	go SubscribeToNATS(conn)

	// Настройка обработчика запросов кэша
	http.HandleFunc("/", cacheHandler)

	// Запуск HTTP-сервера
	fmt.Println("Starting server on port :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр запроса 'orderUID'
	orderUID := r.URL.Query().Get("orderUID")

	// Проверяем, есть ли такой ключ в кэше
	value, ok := cache.Load(orderUID)
	if !ok {
		// Если в кэше нет значения, отправляем ошибку 404
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Приводим значение к ожидаемому типу, в данном случае к типу Order
	order, ok := value.(Order)
	if !ok {
		// Если не удается привести к типу Order, отправляем ошибку сервера
		http.Error(w, "Error converting cached data to type Order", http.StatusInternalServerError)
		return
	}

	// Загружаем HTML шаблон
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Выполняем рендеринг шаблона с данными
	err = tmpl.Execute(w, order)
	if err != nil {
		fmt.Println("Error executing template:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
