/*
 * @Author: fyfishie
 * @Date: 2023-05-10:21
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-17:09
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package mongoproxy

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AllColumnsByIPFilter(IPList []string) (bson.D, *(options.FindOptions)) {
	filter := bson.D{{"ip", bson.D{{"$in", IPList}}}}
	option := options.Find().SetProjection(bson.D{{"_id", 0}})
	return filter, option
}

func IPInDBFilter(IPList []int) (bson.D, *(options.FindOptions)) {
	filter := bson.D{{"ip", bson.D{{"$in", IPList}}}}
	option := options.Find().SetProjection(bson.D{{"ip", 1}})
	option.SetSort(bson.D{{"ip", 1}})
	return filter, option
}

func MPairsIndexFilter() (bson.D, *options.FindOptions) {
	filter := bson.D{}
	option := options.Find().SetProjection(bson.D{{"ip", 1}})
	return filter, option
}

func MAliasPairsFilter(ipList []int) (bson.D, *options.FindOptions) {
	filter := bson.D{{"ip", bson.D{{"$in", ipList}}}}
	option := options.Find().SetProjection(bson.D{{"_id", 0}})
	return filter, option
}
