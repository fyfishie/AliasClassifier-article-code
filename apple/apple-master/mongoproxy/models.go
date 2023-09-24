package mongoproxy

import (
	"aliasParseMaster/lib"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func appendArraysDoc(ipAndArrays []IpAndArray) []mongo.WriteModel {
	models := []mongo.WriteModel{}
	for _, item := range ipAndArrays {
		update := bson.D{
			{"$push", bson.D{
				{"antiAliasSet", bson.D{{"$each", item.AntiAliasSet}}},
			}},
		}

		m := mongo.NewUpdateOneModel().SetFilter(bson.D{{"ip", item.Ip}}).SetUpdate(update).SetUpsert(true)
		models = append(models, m)
	}
	return models
}

func MaybeAliasSetModel(maybe lib.MaybeAlias) mongo.WriteModel {
	return mongo.NewInsertOneModel().SetDocument(maybe)
}

func SureAliasModel(sure lib.SureAlias) mongo.WriteModel {
	return mongo.NewInsertOneModel().SetDocument(sure)
}
