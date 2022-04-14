package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
}

func Connect() Mongo {
	config := LoadConfig()

	var dbURL string
	if config.DBUser != "" && config.DBPass != "" {
		dbURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", config.DBUser, config.DBPass, config.DBHost, config.DBPort)
	} else {
		dbURL = fmt.Sprintf("mongodb://%s:%s", config.DBHost, config.DBPort)
	}

	clientOptions := options.Client().ApplyURI(dbURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalln("MongoDB连接失败: " + err.Error())
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("MongoDB连通性测试失败: " + err.Error())
	}
	return Mongo{
		Client: client,
	}
}

func (m Mongo) Database() *mongo.Database {
	config := LoadConfig()
	return m.Client.Database(config.DBName)
}

func (m Mongo) Disconnect() {
	err := m.Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (m Mongo) Query(collection string, any interface{}, dateStart int64, dateEnd int64) error {
	coll := m.Database().Collection(collection)
	/* filter := bson.D{
		{
			Key: "timestamp",
			Value: bson.D{
				{Key: "$gt", Value: time.Now().Unix() - (second * 60)},
			},
		},
	} */

	filter := bson.M{
		"timestamp": bson.M{
			"$gt": dateStart,
			"$lt": dateEnd,
		},
	}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println("查询错误: " + err.Error())
	}
	if err = cursor.All(context.TODO(), any); err != nil {
		log.Println("查询错误: " + err.Error())
	}

	return err
}

/* func (m Mongo) Query(collection string, any interface{}, second int64) error {
	coll := m.Database().Collection(collection)
	 filter := bson.D{
		{
			Key: "timestamp",
			Value: bson.D{
				{Key: "$gt", Value: time.Now().Unix() - (second * 60)},
			},
		},
	}

	filter := bson.M{
		"timestamp": bson.M{
			"$gt": time.Now().Unix() - (second * 60),
		},
	}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println("查询错误: " + err.Error())
	}
	if err = cursor.All(context.TODO(), any); err != nil {
		log.Println("查询错误: " + err.Error())
	}

	return err
}
*/
