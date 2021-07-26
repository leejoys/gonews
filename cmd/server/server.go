package main

import (
	"encoding/json"
	"fmt"
	"gonews/pkg/api"
	"gonews/pkg/datasource"
	"gonews/pkg/storage"
	"gonews/pkg/storage/mongodb"
	"log"
	"net/http"
	"os"
	"time"
)

// Сервер GoNews.
type server struct {
	ds        *datasource.Source
	db        storage.Interface
	api       *api.API
	postChan  chan storage.Post
	errorChan chan error
}

func main() {
	// Создаём объект сервера
	var srv server

	// Создаем источник данных
	cByte, err := os.ReadFile("./aconfig.json")
	if err != nil {
		log.Fatalf("main ioutil.ReadFile error: %s", err)
	}
	srv.ds = &datasource.Source{}
	err = json.Unmarshal(cByte, srv.ds)
	if err != nil {
		log.Fatalf("main json.Unmarshal error: %s", err)
	}
	srv.ds.PostChan = make(chan storage.Post)
	srv.ds.ErrorChan = make(chan error)

	// Создаём объект базы данных MongoDB.
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf(
		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/gonews?retryWrites=true&w=majority",
		pwd)
	db, err := mongodb.New("gonews", connstr)
	if err != nil {
		log.Fatalf("mongo.New error: %s", err)
	}

	// Инициализируем хранилище сервера БД
	srv.db = db

	// Освобождаем ресурс
	defer srv.db.Close()

	go srv.poster()
	go srv.logger()
	go srv.ds.Run()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	log.Fatal(http.ListenAndServe(":8080", srv.api.Router()))
}

//обрабатывает ответы из каналов с постами
func (s *server) poster() {
	for post := range s.postChan {

		t, err := time.Parse(time.RFC1123, post.PubDate)
		if err != nil {
			s.errorChan <- fmt.Errorf("poster time.Parse error: %s", err)
		}
		post.PubTime = t.Unix()
		err = s.db.AddPost(post)
		if err != nil {
			s.errorChan <- fmt.Errorf("poster storage.AddPost error: %s", err)
		}
	}
}

//обрабатывает ответы из каналов с ошибками
func (s *server) logger() {
	for err := range s.errorChan {
		log.Println(err)
	}
}
