/*
 * @Author: fyfishie
 * @Date: 2023-03-21:08
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-08:16
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package mongoProxy

import (
	"context"
	"log"
	"time"
	"vp/lib"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Proxy struct {
	url            string
	client         *mongo.Client
	dbName         string
	db             *mongo.Database
	collectionName string
	collection     *mongo.Collection
	//每个proxy只有一个用途，如果需要其他功能请新建一个
	used      bool
	connected bool
}

// mongodb://localhost:27017
func NewProxy(addr lib.MongoCollectionAddr) *Proxy {
	url := "mongodb://"
	if addr.Username == "" || addr.Password == "" {
		url = url + addr.IP + ":" + addr.Port
	} else {
		url = url + addr.Username + ":" + addr.Password + "@" + addr.IP + ":" + addr.Port
	}
	m := Proxy{
		url:            url,
		dbName:         addr.DBName,
		collectionName: addr.CollectionName,
		connected:      false,
	}
	return &m
}

func (p *Proxy) Connect() error {
	/*
	   Connect to my cluster
	*/
	client, err := mongo.NewClient(options.Client().ApplyURI(p.url))
	p.client = client
	if err != nil {
		log.Fatal(err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = p.client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//test connection
	_, err = p.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	p.db = client.Database(p.dbName)
	p.collection = p.db.Collection(p.collectionName)
	return nil
}

func (p *Proxy) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := p.client.Disconnect(ctx)
	return err
}
