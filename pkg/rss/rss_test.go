package rss

import (
	"encoding/json"
	"fmt"
	"gonews/pkg/storage/mongodb"
	"os"
	"testing"
	"time"
)

//функциональный тест пакета rss
//результат смотреть в базе данных
func Test_RSSParser(t *testing.T) {
	// Создаём объект базы данных MongoDB.
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf(
		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/gonews?retryWrites=true&w=majority",
		pwd)
	db, err := mongodb.New("gonews", connstr)
	if err != nil {
		t.Fatalf("mongo.New error: %s", err)
	}

	//Создаем обработчик RSS
	cByte, err := os.ReadFile("./aconfig.json")
	if err != nil {
		t.Fatalf("ioutil.ReadFile error: %s", err)
	}
	parser := RSSParser{}
	err = json.Unmarshal(cByte, &parser)
	if err != nil {
		t.Fatalf("json.Unmarshal error: %s", err)
	}
	go parser.Run(db)
	time.Sleep(time.Second * 10)
}
