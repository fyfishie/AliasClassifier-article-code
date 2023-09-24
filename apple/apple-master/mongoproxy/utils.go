/*
 * @Author: fyfishie
 * @Date: 2023-05-09:20
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-17:08
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package mongoproxy

import "go.mongodb.org/mongo-driver/mongo"

func ConvT2Module[T any](raw []T) []mongo.InsertOneModel {
	res := []mongo.InsertOneModel{}
	for _, r := range raw {
		res = append(res, *mongo.NewInsertOneModel().SetDocument(r))
	}
	return res
}
