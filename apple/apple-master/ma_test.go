/*
 * @Author: fyfishie
 * @Date: 2023-05-03:16
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-17:19
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package main

import (
	aliascek "aliasParseMaster/aliaschecker"
	"aliasParseMaster/lib"
	"aliasParseMaster/mongoproxy"
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

/*
uid:
qingdao:0
weihai:1
newark:2
london:3
*/
func Test_mv2mongo(t *testing.T) {
	fs, err := os.ReadDir("./db/spingresult/")
	if err != nil {
		panic(err.Error())
	}
	cityMap := map[string]struct{}{}
	for _, f := range fs {
		ss := strings.Split(f.Name(), "_")
		cityMap[ss[0]] = struct{}{}
	}
	cityList := []string{}
	for city := range cityMap {
		cityList = append(cityList, city)
	}
	sort.Strings(cityList)
	fmt.Println(cityList)
	for index, city := range cityList {
		fmt.Printf("index: %v\n", index)
		fmt.Println(city)
		for _, f := range fs {
			ss := strings.Split(f.Name(), "_")
			if ss[0] == city {
				mv2mongo(f.Name(), ss[1], city)
			}
		}
	}
}
func mv2mongo(rdName string, UID string, collectionName string) {
	bufwtr, err := mongoproxy.NewBufferWriterFromAddr(lib.MongoCollectionAddr{
		IP:             "127.0.0.1",
		Port:           "27017",
		DBName:         "taiwan_manual_sping_result",
		CollectionName: collectionName,
	}, 5000)
	if err != nil {
		panic(err.Error())
	}
	bufwtr.MvSpingRes2Mongo("./db/spingresult/"+rdName, UID)
}

func ipList() []string {
	rfi, err := os.OpenFile("./midar-iff.nodes.cleaned.ip", os.O_RDONLY, 0000)
	if err != nil {
		panic(err)
	}
	defer rfi.Close()
	rdr := bufio.NewReader(rfi)
	res := []string{}
	for {
		line, _, err := rdr.ReadLine()
		if err != nil {
			break
		}
		res = append(res, string(line))
	}
	return res
}
func Test_checker(t *testing.T) {
	fs, err := os.ReadDir("./db/spingresult/")
	if err != nil {
		panic(err.Error())
	}
	cityMap := map[string]struct{}{}
	for _, f := range fs {
		ss := strings.Split(f.Name(), "_")
		cityMap[ss[0]] = struct{}{}
	}
	cityList := []string{}
	for city := range cityMap {
		cityList = append(cityList, city)
	}
	sort.Strings(cityList)
	for index, cityName := range cityList {
		fmt.Printf("index: %v\n", index)
		fmt.Printf("%v start! and index is %v\n", cityName, index)
		checker_onecity(cityName)
		fmt.Println(cityName + " done!")
	}
}
func checker_onecity(cityName string) {
	task := lib.AliasCekTask{
		PingResAddr: lib.MongoCollectionAddr{
			IP:             "127.0.0.1",
			Port:           "27017",
			DBName:         "taiwan_manual_sping_result",
			CollectionName: cityName,
		},
		AntiAliasResAddr: lib.MongoCollectionAddr{
			IP:             "127.0.0.1",
			Port:           "27017",
			DBName:         "taiwan_manual_antialias",
			CollectionName: cityName,
		},
		AliasCekResAddr: lib.MongoCollectionAddr{
			IP:             "127.0.0.1",
			Port:           "27017",
			DBName:         "taiwan_manual_alias_result",
			CollectionName: cityName,
		},
		TSetMCA: lib.MongoCollectionAddr{
			IP:             "127.0.0.1",
			Port:           "27017",
			DBName:         "taiwan_manual_Tset",
			CollectionName: cityName,
		},
		MaybeAliasMCA: lib.MongoCollectionAddr{
			IP:             "127.0.0.1",
			Port:           "27017",
			DBName:         "taiwan_manual_maybe",
			CollectionName: cityName,
		},
	}
	cek := aliascek.NewChecker(&task)
	cek.CheckStart()
}
