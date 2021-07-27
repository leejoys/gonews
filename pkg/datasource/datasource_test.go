package datasource

import (
	"encoding/json"
	"gonews/pkg/datasource/rssparser"
	"gonews/pkg/storage"
	"os"
	"testing"
)

//Интеграционный тест пакета datasource
func TestSource_Run(t *testing.T) {
	// Создаем источник данных
	cByte, err := os.ReadFile("./aconfig.json")
	if err != nil {
		t.Fatalf("main ioutil.ReadFile error: %s", err)
	}
	ds := &Source{}
	err = json.Unmarshal(cByte, ds)
	if err != nil {
		t.Fatalf("main json.Unmarshal error: %s", err)
	}
	ds.PostChan = make(chan storage.Post)
	ds.ErrorChan = make(chan error)
	p := rssparser.New()
	ds.Parser = p
	go ds.Run()

	<-ds.PostChan
}
