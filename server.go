package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	hostIP       string
	cache        sync.Map
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func main() {
	//hostIP = "localhost"
	hostIP = "postgres"

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

	loadCacheFromDB(conn, &cache)

	// Подписка на сообщения NATS в фоновом режиме
	go SubscribeToNATS(conn)

	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())

	// Настройка обработчика запросов кэша
	http.HandleFunc("/", cacheHandler)

	// Запуск HTTP-сервера
	fmt.Println("Starting server on port :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	// Загружаем HTML шаблон
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Если запрос - POST, обрабатываем поиск
	if r.Method == "POST" {
		// Парсим форму
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusInternalServerError)
			return
		}

		// Получаем orderUID из формы
		orderUID := r.FormValue("orderUID")

		// Проверяем, есть ли такой ключ в кэше
		value, ok := cache.Load(orderUID)
		if !ok {
			// Если в кэше нет значения, выводим сообщение об ошибке
			tmpl.Execute(w, "Order not found")
			return
		}

		// Приводим значение к ожидаемому типу, в данном случае к типу Order
		order, ok := value.(Order)
		if !ok {
			// Если не удается привести к типу Order, отправляем сообщение об ошибке
			tmpl.Execute(w, "Error converting cached data to type Order")
			return
		}

		// Выполняем рендеринг шаблона с данными заказа
		err = tmpl.Execute(w, order)
		if err != nil {
			fmt.Println("Error executing template:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Если запрос - GET, показываем форму поиска
	tmpl.Execute(w, nil)
}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}
