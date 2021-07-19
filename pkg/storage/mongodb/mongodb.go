package mongodb

import (
	"context"
	"gonews/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Хранилище данных.
type Store struct {
	c  *mongo.Client
	db *mongo.Database
}

//New - Конструктор объекта хранилища.
func New(name string, connstr string) (*Store, error) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(connstr))
	if err != nil {
		return nil, err
	}
	// проверка связи с БД
	err = client.Ping(context.Background(), nil)
	if err != nil {
		client.Disconnect(context.Background())
		return nil, err
	}

	s := &Store{c: client, db: client.Database(name)}
	t := true
	_, err = s.db.Collection("posts").Indexes().CreateOne(
		context.Background(), mongo.IndexModel{
			Keys:    bson.D{{Key: "title", Value: 1}},
			Options: &options.IndexOptions{Unique: &t}})
	if err != nil {
		s.c.Disconnect(context.Background())
		return nil, err
	}

	return s, nil
}

//Close - освобождение ресурса
func (s *Store) Close() {
	s.c.Disconnect(context.Background())
}

//Posts - получение всех публикаций
func (s *Store) Posts() ([]storage.Post, error) {

	coll := s.db.Collection("posts")
	ctx := context.Background()
	filter := bson.D{}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	posts := []storage.Post{}
	for cur.Next(ctx) {
		var p storage.Post
		err = cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

//PostsN - получение n последних публикаций
func (s *Store) PostsN(n int) ([]storage.Post, error) {

	coll := s.db.Collection("posts")
	ctx := context.Background()
	filter := bson.D{}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	posts := []storage.Post{}
	for cur.Next(ctx) {
		var p storage.Post
		err = cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

//AddPost - создание новой публикации
func (s *Store) AddPost(p storage.Post) error {
	coll := s.db.Collection("posts")
	_, err := coll.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

//UpdatePost - обновление по id значения title, content, author_id и published_at
func (s *Store) UpdatePost(p storage.Post) error {
	coll := s.db.Collection("posts")
	filter := bson.D{{Key: "id", Value: p.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "title", Value: p.Title},
		{Key: "content", Value: p.Content},
		{Key: "pubtime", Value: p.PubTime},
		{Key: "link", Value: p.Link}}}}
	_, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

//DeletePost - удаляет пост по id
func (s *Store) DeletePost(p storage.Post) error {
	coll := s.db.Collection("posts")
	filter := bson.D{{Key: "id", Value: p.ID}}
	_, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
