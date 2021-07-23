package main

import (
	"encoding/json"
	"fmt"
	"gonews/pkg/api"
	"gonews/pkg/rss"
	"gonews/pkg/storage"
	"gonews/pkg/storage/mongodb"
	"log"
	"net/http"
	"os"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объект базы данных MongoDB.
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf(
		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/gonews?retryWrites=true&w=majority",
		pwd)
	db, err := mongodb.New("gonews", connstr)
	if err != nil {
		log.Fatalf("mongo.New error: %s", err)
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Освобождаем ресурс
	defer srv.db.Close()

	cByte, err := os.ReadFile("./aconfig.json")
	if err != nil {
		log.Fatalf("main ioutil.ReadFile error: %s", err)
	}
	parser := rss.RSSParser{}
	err = json.Unmarshal(cByte, &parser)
	if err != nil {
		log.Fatalf("main json.Unmarshal error: %s", err)
	}
	go parser.Run(db)

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	log.Fatal(http.ListenAndServe(":8080", srv.api.Router()))
}
