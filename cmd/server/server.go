package main

import (
	"encoding/json"
	"fmt"
	"gonews/pkg/api"
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

//работает с RSS, хранит конфигурацию
type rssEater struct {
	RSS            []string
	Request_period int
	postChan       chan (storage.Post)
	errorChan      chan (error)
}

//запускает опрос заданных RSS
func (c *rssEater) Run() {
	for _, link := range c.RSS {
		go rssEat(link, c.Request_period)
	}
}

func rssEat(link string, period int) {

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

	cFile, err := os.Open("./config.json")
	if err != nil {
		log.Fatalf("os.Open error: %s", err)
	}
	cByte, err := ioutil.ReadAll(cFile)
	if err != nil {
		log.Fatalf("ioutil.ReadAll error: %s", err)
	}
	conf := rssEater{}
	err = json.Unmarshal(cByte, &conf)
	if err != nil {
		log.Fatalf("json.Unmarshal error: %s", err)
	}
	conf.Run()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
