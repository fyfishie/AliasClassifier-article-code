/*
 * @Author: fyfishie
 * @Date: 2023-05-14:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-17:10
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
	"vp/lib"
	"vp/utils"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
	"github.com/fyfishie/ipop"
	"github.com/sirupsen/logrus"
)

func Test_taiwan(t *testing.T) {
	fs, err := os.ReadDir("./db/trace_part")
	esyerr.AutoPanic(err)
	for _, f := range fs {
		fmt.Println(f.Name() + "start")
		if f.Name() != "Taipei" {
			continue
		}
		anti_city("./db/trace_part/"+f.Name(), "./db/taiwanip_fromtrace/"+f.Name(), f.Name())
		fmt.Println(f.Name() + "end")
	}
}
func anti_city(tracePath, ipListpath string, cityName string) {
	task := lib.SlaveTask{
		TopTaskName: cityName,
		TopTaskID:   6,
		IPToDoList:  IPList(ipListpath),
		AntiAliasResAddr: lib.MongoCollectionAddr{
			IP:             "127.0.0.1",
			Port:           "27017",
			DBName:         "taiwan_manual_antialias_int64",
			CollectionName: cityName,
		},
	}
	fmt.Println("onetrace start")
	err := parseOneTraceAndMvResult2Mongo(task.AntiAliasResAddr, tracePath)
	fmt.Println("onetrace done")
	if err != nil {
		logrus.Errorf("error while parsing one trace data, err = %v\n", err.Error())
	}
	fmt.Println("end2 start")
	err = parseEnd2AndMv2Mongo(task.AntiAliasResAddr, tracePath)
	fmt.Println("end2 end")
	if err != nil {
		logrus.Errorf("error while parsing end2 data, err = %v\n", err.Error())
	}
	fmt.Println("subnet start")
	err = parseSubnet(task.AntiAliasResAddr, utils.StringIPs2Ints(task.IPToDoList))
	fmt.Println("subnet end")
	if err != nil {
		logrus.Errorf("error while parsing subnet data, err = %v\n", err.Error())
	}
}

func IPList(path string) []string {
	res := []string{}
	rfi, err := os.OpenFile(path, os.O_RDONLY, 0000)
	if err != nil {
		panic(err.Error())
	}
	defer rfi.Close()
	scaner := bufio.NewScanner(rfi)
	for scaner.Scan() {
		res = append(res, scaner.Text())
	}
	return res
}

func Test_end2(t *testing.T) {
	end2AndEndsMap := map[int][]int{}
	rfi, err := os.OpenFile("./db/trace_part/Banqiao", os.O_RDONLY, 0000)
	if err != nil {
		return
	}
	defer rfi.Close()
	rdr := bufio.NewReader(rfi)
	count := 0
	for {
		count++
		if count%10000 == 0 {
			fmt.Printf("count: %v\n", count)
		}
		line, _, err := rdr.ReadLine()
		if err != nil {
			break
		}
		rt := lib.RawTrace{}
		err = json.Unmarshal(line, &rt)
		if err != nil {
			continue
		}
		if len(rt.Results) < 2 {
			continue
		}
		if rt.Ip != rt.Results[len(rt.Results)-1].Ip {
			continue
		}
		trace := utils.RawTrace2Trace(rt)
		if _, ok := end2AndEndsMap[trace.End2]; !ok {
			end2AndEndsMap[trace.End2] = []int{}
		}
		end2AndEndsMap[trace.End2] = append(end2AndEndsMap[trace.End2], trace.End)
	}
	max := 0
	for _, ends := range end2AndEndsMap {
		if len(ends) > 1000 {
			fmt.Println("?")
			sort.Ints(ends)
		}
		max = utils.Max(max, len(ends))
	}
	fmt.Printf("max: %v\n", max)

}

func Test_traceCount(t *testing.T) {
	ipCountMap := map[int]int{}
	loader := store.NewLoader[lib.RawTrace]("./db/trace_part/Banqiao")
	err := loader.Open()
	if err != nil {
		panic(err.Error())
	}
	count := 0
	for {
		nexts, err := loader.Next(5000)
		if err != nil {
			for _, r := range nexts {
				if len(r.Results) < 2 {
					continue
				}
				if r.Ip != r.Results[len(r.Results)-1].Ip {
					continue
				}
				ipCountMap[ipop.String2Int(r.Ip)]++
			}
			break
		}
		for _, r := range nexts {
			count++
			if count%100000 == 0 {
				fmt.Printf("count: %v\n", count)
			}
			if len(r.Results) < 2 {
				continue
			}
			ipCountMap[ipop.String2Int(r.Ip)]++
		}
	}
	// fmt.Println(end2AndEndsMap)
	max := 0
	for _, count := range ipCountMap {
		if count > 10000 {
			fmt.Println("?")
		}
		max = utils.Max(max, count)
	}
	fmt.Printf("max: %v\n", max)
}
