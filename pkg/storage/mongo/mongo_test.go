package mongo

import (
	"fmt"
	"gonews/pkg/storage"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMongo(t *testing.T) {
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf(
		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/gonews?retryWrites=true&w=majority",
		pwd)
	db, err := New("gonews", connstr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	posts := []storage.Post{
		{1, "Вышел Microsoft Linux",
			"Как сообщают непроверенные источники, новая ОС будет бесплатной.",
			time.Now().Unix(), "https://github.com/microsoft/CBL-Mariner"},
		{2, "Инженеры Google не желают возвращаться в офисы",
			"Инженеры Google не желают возвращаться в офисы, заявляя, что они не менее продуктивны на удалёнке.",
			time.Now().Unix(), "https://habr.com/ru/news/t/568128/"}}
	for _, p := range posts {
		err = db.AddPost(p)
		if err != nil {
			t.Fatalf("AddPost error: %s", err)
		}
	}

	received, err := db.Posts()
	if err != nil {
		t.Fatalf("Posts error: %s", err)
	}
	if !reflect.DeepEqual(posts, received) {
		t.Errorf("received %v, wanted %v", received, posts)
	}

	for _, p := range posts {
		err := db.DeletePost(p)
		if err != nil {
			t.Fatalf("DeletePost error: %s", err)
		}
	}
	received, err = db.Posts()
	if err != nil {
		t.Fatalf("Posts error: %s", err)
	}
	posts = []storage.Post{}
	if !reflect.DeepEqual(posts, received) {
		t.Errorf("received %v, wanted %v", received, posts)
	}
}
