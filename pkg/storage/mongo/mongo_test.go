package mongo

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMongo(t *testing.T) {
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf(
		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/gonews?retryWrites=true&w=majority",
		pwd)
	db, err := New("gonews", connstr)
	if err != nil {
		log.Fatal(err)
	}

}
