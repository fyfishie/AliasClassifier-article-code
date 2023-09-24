/*
 * @Author: fyfishie
 * @Date: 2023-02-20:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-15:10
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package mongoproxy

import (
	"aliasParseMaster/lib"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WriteData struct {
	IpAndArray IpAndArray
	Opration   int //INSERT,APPEND...
}
type ReadReq []int
type ReadRes []IpAndArray
type IpAndArray struct {
	Ip           int   `json:"ip"`
	AntiAliasSet []int `json:"antiAliasSet"`
}

const INSERT int = 0
const APPEND int = 1
const FLUSH int = 2 //通知mongoProxy将所有未写入的数据写入数据库

type Proxy struct {
	url            string
	client         *mongo.Client
	ctxToDo        context.Context
	dbName         string
	db             *mongo.Database
	collectionName string
	collection     *mongo.Collection
	//每个proxy只有一个用途，如果需要其他功能请新建一个
	used bool
}

// mongodb://localhost:27017
func NewProxy(addr lib.MongoCollectionAddr) Proxy {
	url := "mongodb://"
	if addr.Username == "" || addr.Password == "" {
		url = url + addr.IP + ":" + addr.Port
	} else {
		url = url + addr.Username + ":" + addr.Password + "@" + addr.IP + ":" + addr.Port
	}
	m := Proxy{
		url:            url,
		ctxToDo:        context.TODO(),
		dbName:         addr.DBName,
		collectionName: addr.CollectionName,
	}
	return m
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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = p.client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}

	/*
	   List databases
	*/
	_, err = p.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	p.db = client.Database(p.dbName)
	p.collection = p.db.Collection(p.collectionName)
	return nil
}

func (p *Proxy) Disconnect() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := p.client.Disconnect(ctx)
	return err
}
