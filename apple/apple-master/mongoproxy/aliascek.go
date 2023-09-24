/*
 * @Author: fyfishie
 * @Date: 2023-05-15:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-17:19
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package mongoproxy

import (
	"aliasParseMaster/lib"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//this file is function for aliascek

// /*
// get ip list whose sping result is contained by all vp
// */
func (p *Proxy) AndVPList() []string {
	// get all uid
	qres, err := p.collection.Distinct(context.TODO(), "vpuid", bson.D{}, nil)
	if err != nil {
		return nil
	}
	res := []string{}
	for _, q := range qres {
		res = append(res, q.(string))
	}
	return res
}

// find IP who is  detected by all vp using sping
func (p *Proxy) ConsistenceIP() ([]string, []lib.VPUID, error) {
	andVP := p.AndVPList()
	if len(andVP) < 2 {
		return nil, nil, nil
	}
	andMap := map[string]struct{}{}
	filter := bson.D{{"vpuid", andVP[0]}}
	qres, err := p.collection.Distinct(context.TODO(), "ip", filter, nil)
	if err != nil {
		return nil, nil, err
	}
	for _, qr := range qres {
		andMap[qr.(string)] = struct{}{}
	}
	for _, vpUID := range andVP[1:] {
		filter := bson.D{{"vpuid", vpUID}}
		// projection := bson.D{{"ip", 1}}
		qres, err := p.collection.Distinct(context.TODO(), "ip", filter, nil)
		if err != nil {
			return nil, nil, err
		}
		for _, qr := range qres {
			if _, ok := andMap[qr.(string)]; !ok {
				delete(andMap, qr.(string))
			}
		}
	}

	res := []string{}
	for k := range andMap {
		res = append(res, k)
	}
	return res, andVP, nil
}

func (p *Proxy) ReadSpingWithUID(ipList []string) ([]lib.SpingWithUID, error) {
	filter := bson.D{{"ip", bson.D{{"$in", ipList}}}}
	corsur, err := p.collection.Find(context.TODO(), filter, nil)
	if err != nil {
		return nil, err
	}
	res := []lib.SpingWithUID{}
	err = corsur.All(context.TODO(), &res)
	return res, err
}

func (p *Proxy) QueryIPLeaderinDB(consistentIPList []int, ctx context.Context) (map[int]struct{}, error) {
	ires, err := p.collection.Distinct(ctx,
		"ip",
		bson.D{{"ip", bson.D{{"$in", consistentIPList}}}},
		nil)
	if err != nil {
		return nil, err
	}
	res := map[int]struct{}{}
	for _, r := range ires {
		switch r.(type) {
		case int64:
			res[int(r.(int64))] = struct{}{}
		case int32:
			res[int(r.(int32))] = struct{}{}
		}
	}
	return res, nil
}
func (p *Proxy) QueryIPLeaderinMaybe(ctx context.Context) ([]int, error) {
	ires, err := p.collection.Distinct(ctx, "ip", bson.D{}, nil)
	if err != nil {
		return nil, err
	}
	res := []int{}
	for _, r := range ires {
		switch r.(type) {
		case int64:
			res = append(res, int(r.(int64)))
		case int32:
			res = append(res, int(r.(int32)))
		}
	}
	return res, nil
}
func (p *Proxy) QueryMaybeAlias(ipList []int, ctx context.Context) ([]lib.MaybeAlias, error) {
	corsor, err := p.collection.Find(ctx, bson.D{{"ip", bson.D{{"$in", ipList}}}}, options.Find().SetProjection(bson.D{{"_id", 0}}))
	if err != nil {
		return nil, err
	}
	res := []lib.MaybeAlias{}
	err = corsor.All(ctx, &res)
	return res, err

}
