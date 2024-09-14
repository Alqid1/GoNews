package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "proj"   // имя учебной БД
	collectionName = "gonews" // имя коллекции в учебной БД
)

type Storage struct {
	db *mongo.Client
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		log.Fatal(err)
	}
	// проверка связи с БД
	err = db.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

func (s *Storage) Posts() ([]storage.Post, error) {
	db := s.db.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	cur, err := db.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var posts []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, cur.Err()
}

func (s *Storage) AddPost(p storage.Post) error {
	db := s.db.Database(databaseName).Collection(collectionName)
	_, err := db.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdatePost(p storage.Post) error {
	db := s.db.Database(databaseName).Collection(collectionName)

	filter := bson.M{"id": p.ID}
	update := bson.M{
		"$set": bson.M{
			"title":       p.Title,
			"content":     p.Content,
			"createdat":   p.CreatedAt,
			"publishedat": p.PublishedAt,
		},
	}
	_, err := db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeletePost(p storage.Post) error {
	db := s.db.Database(databaseName).Collection(collectionName)
	_, err := db.DeleteOne(context.Background(), p)
	if err != nil {
		return err
	}

	return nil
}
