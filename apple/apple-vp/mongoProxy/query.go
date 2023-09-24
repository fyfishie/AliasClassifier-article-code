/*
 * @Author: fyfishie
 * @Date: 2023-03-21:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-14:18
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package mongoProxy

import (
	"context"
	"fmt"
	"vp/lib"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
func AntiDescriptorModel(descriptor []lib.AntiDescriptor) []mongo.WriteModel {
	models := []mongo.WriteModel{}
	for _, item := range descriptor {
		m := mongo.NewInsertOneModel().SetDocument(item)
		models = append(models, m)
	}
	return models
}
func (p *Proxy) QueryAllEnd2IP() ([]int, error) {
	res := []int{}
	cursor, err := p.collection.Find(context.TODO(), bson.D{}, options.Find().SetProjection(bson.D{}))
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		item := lib.End2AndEnds{}
		err := cursor.Decode(&item)
		if err != nil {
			continue
		}
		res = append(res, item.End2)
	}
	return res, nil
}

/*
	type End2AndItsEnds struct {
		End2 int   `json:"end2"`
		Ends []int `json:"ends"`
	}
*/
func (p *Proxy) QueryEndsOfEnd2(end2List []int, outChan chan lib.End2AndEnds) {
	total := len(end2List)
	for index, end2 := range end2List {
		if index%1000 == 0 {
			fmt.Printf("%v/%v\n", index, total)
		}
		cursor, err := p.collection.Find(
			context.TODO(),
			bson.D{
				{"end2", end2},
			},
		)
		if err != nil {
			continue
		}
		items := []lib.End2AndEnds{}
		cursor.All(context.TODO(), &items)
		end2List := lib.End2AndEnds{
			End2: end2,
			Ends: []int{},
		}
		for _, item := range items {
			end2List.Ends = append(end2List.Ends, item.Ends...)
		}
		outChan <- end2List
	}
	close(outChan)
}
