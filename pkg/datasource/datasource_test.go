package datasource

import (
	"fmt"
	"gonews/pkg/storage/mongodb"
	"os"
	"testing"
	"time"
)

func TestSource_Run(t *testing.T) {
	// Создаём объект базы данных MongoDB.
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf(
		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/dbtest?retryWrites=true&w=majority",
		pwd)
	db, err := mongodb.New("dbtest", connstr)
	if err != nil {
		t.Fatalf("mongo.New error: %s", err)
	}
	defer db.Close()
	defer func() {
		if err = db.DropDB(); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second * 10)
}
