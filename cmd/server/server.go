package main

import (
	"encoding/json"
	"fmt"
	"gonews/pkg/api"
	"gonews/pkg/rss"
	"gonews/pkg/storage"
	"gonews/pkg/storage/mongo"
	"io/ioutil"
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
	db, err := mongo.New("gonews", connstr)
	if err != nil {
		log.Fatalf("mongo.New error: %s", err)
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Освобождаем ресурс
	defer srv.db.Close()

	cFile, err := os.Open("./bconfig.json")
	if err != nil {
		log.Fatalf("main os.Open error: %s", err)
	}
	cByte, err := ioutil.ReadAll(cFile)
	if err != nil {
		log.Fatalf("main ioutil.ReadAll error: %s", err)
	}
	eater := rss.RSSEater{}
	err = json.Unmarshal(cByte, &eater)
	if err != nil {
		log.Fatalf("main json.Unmarshal error: %s", err)
	}
	eater.Run()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
