package data

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	InsertChat(ChatEntry) (string, error)
	GetChat(string) (ChatEntry, error)
	UpdateChat(string, ChatEntry) error
	DeleteChat(string) error
}

type Mongo struct {
	db *mongo.Database
}

const (
	database   = "interview"
	collection = "chat"
)

func NewMongo(uri string) *Mongo {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	return &Mongo{
		db: client.Database(database),
	}
}

func (m *Mongo) InsertChat(data ChatEntry) (string, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	_, err := m.db.Collection(collection).InsertOne(context.Background(), data)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (m *Mongo) GetChat(id string) (ChatEntry, error) {
	var data ChatEntry

	err := m.db.Collection(collection).FindOne(context.Background(), bson.M{"_id": id}).Decode(&data)
	if err != nil {
		return ChatEntry{}, err
	}

	return data, nil
}

func (m *Mongo) UpdateChat(id string, data ChatEntry) error {
	_, err := m.db.Collection(collection).UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": data})
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) DeleteChat(id string) error {
	_, err := m.db.Collection(collection).DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
