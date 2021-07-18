package rss

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

//функциональный тест пакета rss
func Test_rssEater(t *testing.T) {
	cFile, err := os.Open("./bconfig.json")
	if err != nil {
		log.Fatalf("main os.Open error: %s", err)
	}
	cByte, err := ioutil.ReadAll(cFile)
	if err != nil {
		log.Fatalf("main ioutil.ReadAll error: %s", err)
	}
	eater := RSSEater{}
	err = json.Unmarshal(cByte, &eater)
	if err != nil {
		log.Fatalf("main json.Unmarshal error: %s", err)
	}
	eater.Run()
}
