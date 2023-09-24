package mongoproxy

import (
	"aliasParseMaster/lib"
	"context"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
cursor, err := coll.Find(
	context.TODO(),
	bson.D{{"status", bson.D{{"$in", bson.A{"A", "D"}}}}})
*/
//[from,to)
func (p *Proxy) QueryMany(ips []int, options *options.FindOptions) []IpAndArray {
	cursor, err := p.collection.Find(
		context.TODO(),
		bson.D{
			{"ip", bson.D{{"$in", ips}}},
		},
		options)
	if err != nil {
		log.Printf("error in get cursor of result, err = %v\n", err.Error())
		return nil
	}
	queryRes := []IpAndArray{}
	if err = cursor.All(context.TODO(), &queryRes); err != nil {
		log.Printf("error in query many, err = %v\n", err.Error())
		return nil
	}
	if len(queryRes) > 0 {
		return queryRes
	}
	return nil
}

func (p *Proxy) BulkWrite(models []mongo.WriteModel) bool {
	opts := options.BulkWrite().SetOrdered(false)
	_, err := p.collection.BulkWrite(context.TODO(), models, opts)
	return err == nil
}

func (p *Proxy) Distinct() {
	result, err := p.collection.Distinct(context.TODO(), "antiAliasSet", bson.D{})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("result: %v\n", result)
	}
}

// func appendArraysDoc(ipAndArrays []IpAndArray) []mongo.WriteModel {
// 	models := []mongo.WriteModel{}
// 	for _, item := range ipAndArrays {
// 		update := bson.D{
// 			{"$push", bson.D{
// 				{"antiAliasSet", bson.D{{"$each", item.AntiAliasSet}}},
// 			}},
// 		}

// 		m := mongo.NewUpdateOneModel().SetFilter(bson.D{{"ip", item.Ip}}).SetUpdate(update).SetUpsert(true)
// 		models = append(models, m)
// 	}
// 	return models
// }

// func insertArraysDoc(ipAndArrays []IpAndArray) []mongo.WriteModel {
// 	models := []mongo.WriteModel{}
// 	for _, item := range ipAndArrays {
// 		m := mongo.NewInsertOneModel().SetDocument(item)
// 		models = append(models, m)
// 	}
// 	return models
// }

// 查询集合中的非别名组key值
// 已经排好序
func (p *Proxy) QueryIps(ipToDoList []int) (map[int]bool, []int) {
	projection := bson.D{
		{"ip", 1},
		{"_id", 0},
	}
	filter := bson.D{{"ip", bson.D{{"$in", ipToDoList}}}}
	cursor, err := p.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetProjection(projection))
	var queryRes []IpAndArray
	if err = cursor.All(context.TODO(), &queryRes); err != nil {
		panic(err)
	}
	if len(queryRes) > 0 {
		ipInDBList := []int{}
		ipInDB := map[int]bool{}
		for _, item := range queryRes {
			ipInDB[item.Ip] = true
			ipInDBList = append(ipInDBList, item.Ip)
		}
		sort.Ints(ipInDBList)
		return ipInDB, ipInDBList
	}

	return nil, nil
}

func (p *Proxy) InsertManyInterfaces(docs []interface{}) error {
	_, err := p.collection.InsertMany(context.TODO(), docs, nil)
	return err
}

// 查询得到IPList和数据库中ip字段组成集合的交集
func (p *Proxy) GetAllIPAmong(IPList []int) (map[int]bool, error) {
	if p.used {
		return nil, errors.New("this proxy is used, make a new one please")
	}
	err := p.Connect()
	if err != nil {
		return nil, err
	}
	projection := bson.D{
		{"ip", 1},
		{"_id", 0},
	}
	filter := bson.D{{"ip", bson.D{{"$in", IPList}}}}
	cursor, err := p.collection.Find(
		context.TODO(),
		filter,
		options.Find().SetProjection(projection))
	ips := []int{}
	err = cursor.All(context.Background(), &ips)
	if err != nil {
		return nil, err
	}
	ipsMap := map[int]bool{}
	for _, ip := range ips {
		ipsMap[ip] = true
	}
	return ipsMap, nil
}

/*
@description:

@param {chan[]lib.IP} ipGroupChan

@return {*}

@block: false
*/
func (p *Proxy) QueryIPGroupByLocID(ipGroupChan chan []lib.IntIP) (groupNum int, err error) {
	//query two times and I have no more time to make it more graceful
	cur, err := p.collection.Find(
		context.TODO(),
		bson.D{},
		options.Find().SetProjection(bson.D{}))
	allLocIDMap := map[int]struct{}{}
	for cur.Next(context.TODO()) {
		cityAndIPs := lib.CityIPForMongo{}
		err = cur.Decode(&cityAndIPs)
		if err != nil {
			continue
		}
		allLocIDMap[cityAndIPs.LocIntID] = struct{}{}
	}
	cur.Close(context.TODO())

	//now let's query ip list by city
	go func() {
		for LocIntID, _ := range allLocIDMap {
			cityIPs := []lib.CityIPForMongo{}
			cur, err := p.collection.Find(
				context.TODO(),
				bson.D{
					{"locintid", LocIntID},
				},
			)
			if err != nil {
				logrus.Errorf("error while query iplist for a city, process continued, err = %v\n", err.Error())
				continue
			}
			cur.All(context.TODO(), &cityIPs)
			oneIPGroup := []int{}
			for _, c := range cityIPs {
				for _, ip := range c.IPList {
					oneIPGroup = append(oneIPGroup, ip)
				}
			}
			ipGroupChan <- oneIPGroup
		}
		close(ipGroupChan)
		err = p.collection.Drop(context.TODO())
		if err != nil {
			logrus.Errorf("failed to drop tmp collection, err = %v\n", err.Error())
		}
	}()
	return len(allLocIDMap), nil
}
func (p *Proxy) Query(filter interface{}, options ...*options.FindOptions) ([]interface{}, error) {
	cur, err := p.collection.Find(context.Background(), filter, options...)
	if err != nil {
		logrus.Error("error in find data from mongodb, err = " + err.Error())
		return nil, err
	}
	res := []interface{}{}
	err = cur.All(context.TODO(), &res)
	if err != nil {
		logrus.Error("error in decode result from mongo find, err = " + err.Error())
	}
	return res, err
}
