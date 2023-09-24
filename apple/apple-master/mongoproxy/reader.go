/*
 * @Author: fyfishie
 * @Date: 2023-05-10:21
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-17:19
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
/*
 * @Author: fyfishie
 * @Date: 2023-03-02:08
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-03-21:17
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package mongoproxy

import (
	// "aliasParseMaster/filtermaker"
	"aliasParseMaster/lib"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// read antialias ip set for each ip of ipList, returns result by readChan.
func (p *Proxy) RunReadPipeline(ipList []int, blockSize int) (readChan chan map[int]struct{}) {
	readChan = make(chan map[int]struct{}, blockSize*2)
	go p.readPipeline(ipList, blockSize, readChan)
	return readChan
}
func (p *Proxy) readPipeline(ipToDoList []int, blockSize int, readChan chan map[int]struct{}) {
	options := options.FindOptions{}
	options.SetSort(bson.D{{"ip", 1}})
	blockNum := len(ipToDoList) / blockSize
	rightEdge := blockNum * blockSize
	i := 0
	for i = 0; i < rightEdge; i += blockSize {
		// fmt.Println("query block")
		ips := ipToDoList[i : i+blockSize]
		queryRes := p.QueryMany(ips, &options)
		if len(queryRes) > 0 {
			sendReadRes(readChan, queryRes)
		}
	}
	queryRes := p.QueryMany(ipToDoList[i:], &options)
	if len(queryRes) > 0 {
		sendReadRes(readChan, queryRes)
	}
}

func sendReadRes(sendChan chan map[int]struct{}, queryRes []IpAndArray) {
	lastIP := queryRes[0].Ip
	m := map[int]struct{}{}
	for _, item := range queryRes {
		if lastIP == item.Ip {
			for _, ip := range item.AntiAliasSet {
				m[ip] = struct{}{}
			}
		} else {
			sendChan <- m
			lastIP = item.Ip
			m = map[int]struct{}{}
			for _, ip := range item.AntiAliasSet {
				m[ip] = struct{}{}
			}
		}
	}
	//最后一个m
	sendChan <- m
}

/*
 * @description: not a good implement :)
 */
type MPairsBufReader struct {
	mPairsAddr  lib.MongoCollectionAddr
	mPairProxy  *Proxy
	bufLen      int
	currentLen  int
	buffer      []lib.MaybeAlias
	cursor      int
	ipIndex     []int
	mongoCursor int
}

func NewMPairsBufReader(mPairsAddr lib.MongoCollectionAddr, buflen int) (*MPairsBufReader, error) {
	b := MPairsBufReader{
		mPairsAddr:  mPairsAddr,
		bufLen:      buflen,
		currentLen:  0,
		buffer:      make([]lib.MaybeAlias, buflen),
		cursor:      -1,
		mongoCursor: 0,
	}
	p := NewProxy(b.mPairsAddr)
	err := p.Connect()
	if err != nil {
		return nil, err
	}
	b.mPairProxy = &p
	err = b.allIndexIP()
	return &b, err
}

func (r *MPairsBufReader) NextOne() (lib.MaybeAlias, error) {
	if r.currentLen > 0 {
		r.cursor++
		r.currentLen--
		return r.buffer[r.cursor], nil
	}
	news, err := r.getNum(r.bufLen)
	if err != nil {
		return lib.MaybeAlias{}, err
	}
	r.currentLen = len(news)
	r.cursor = -1
	r.buffer = news
	if r.currentLen == 0 {
		return lib.MaybeAlias{}, errors.New("no more")
	}
	r.currentLen--
	r.cursor++
	return r.buffer[r.cursor], nil
}

// 返回num条结果，不足num个则把全部返回
func (r *MPairsBufReader) getNum(num int) ([]lib.MaybeAlias, error) {
	if r.mongoCursor == len(r.ipIndex) {
		return []lib.MaybeAlias{}, nil
	}
	if r.mongoCursor+num-1 >= len(r.ipIndex) {
		num = len(r.ipIndex) - r.mongoCursor - 1
	}
	nextList := r.ipIndex[r.mongoCursor : r.mongoCursor+num]
	r.mongoCursor += num
	// filter, option := MAliasPairsFilter(nextList)
	res, err := r.mPairProxy.QueryMaybeAlias(nextList, context.TODO())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *MPairsBufReader) allIndexIP() error {
	m, err := r.mPairProxy.QueryIPLeaderinMaybe(context.TODO())
	if err != nil {
		return err
	}
	r.ipIndex = m
	return nil
}

func (p *Proxy) RunReadPipeline1(filter bson.D, options *options.FindOptions) (chan interface{}, error) {
	resultChan := make(chan interface{})
	cursor, err := p.collection.Find(
		context.TODO(),
		filter,
		options)
	if err != nil {
		log.Printf("error in get cursor of result, err = %v\n", err.Error())
		return resultChan, err
	}

	go func() {
		for {
			var val interface{}
			err := cursor.Decode(&val)
			if err != nil {
				close(resultChan)
				break
			}
			resultChan <- val
		}
	}()
	return resultChan, nil
}
