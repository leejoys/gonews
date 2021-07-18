package rss

import (
	"encoding/json"
	"fmt"
	"gonews/pkg/storage/mongodb"
	"io/ioutil"
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
	cFile, err := os.Open("./aconfig.json")
	if err != nil {
		t.Fatalf("main os.Open error: %s", err)
	}
	cByte, err := ioutil.ReadAll(cFile)
	if err != nil {
		t.Fatalf("main ioutil.ReadAll error: %s", err)
	}
	parser := RSSParser{}
	err = json.Unmarshal(cByte, &parser)
	if err != nil {
		t.Fatalf("main json.Unmarshal error: %s", err)
	}
	go parser.Run(db)
	time.Sleep(time.Second * 4)
}
