package utils

import (
	"alias_article/lib"
	"fmt"
	"github.com/fyfishie/dorapock/store"
	"math"
	"os"
)

var vpList = map[string]string{
	"1": "./db/database/distect/36.138.22.160/36.138.22.160.1/",
	"2": "./db/database/distect/36.140.40.41/36.140.40.41.1/",
	"3": "./db/database/distect/36.140.14.122/36.140.14.122.1/",
	"4": "./db/database/distect/london.tect/london.1/",
	"5": "./db/database/distect/newark/newark.1/",
	"6": "./db/database/distect/qingdao/qingdao.1/",
}

var allName = []string{
	"6-15",
	"16-25",
	"26-35",
	"36-45",
	"46-55",
	"56-65",
	"66-75",
	"76-85",
	"86-95",
	"96-105",
	"106-115",
	"116-125",
	"136-145",
	"176-185",
	"206-215",
	"216-225",
	"226-235",
	"236-245",
	"246-255",
}

var Num_w = []int{1, 10, 50, 100, 200, 500, 800, 1000}

// cut target ip list into many blocks
// it seems to be a slide window, which has a length of 10(ms)
func Slide_windows(spingPath string, allIP []string) (groups [][]string, unreachable []string) {
	spings, err := store.LoadAny[lib.Sping](spingPath)
	if err != nil {
		panic(err)
	}
	if len(spings) == 0 {
		return nil, allIP
	}
	spingMap := map[string]lib.Sping{}
	for _, sp := range spings {
		spingMap[sp.IP] = sp
	}
	un := []string{}
	for _, ip := range allIP {
		if _, ok := spingMap[ip]; !ok {
			un = append(un, ip)
		}
	}
	min := math.MaxInt
	max := 0
	for _, v := range spingMap {
		if v.Ttl > max {
			max = v.Ttl
		}
		if v.Ttl < min {
			min = v.Ttl
		}
	}
	groups = [][]string{}
	var i int
	for i = min; i < max; i += 10 {
		left := i
		right := i + 9
		os.Mkdir("./db/database/slide_window/"+fmt.Sprintf("%v-%v", left, right), 0777)
		group := get_slide(left, right, spingMap)
		if group != nil {
			groups = append(groups, group)
		}
	}
	i -= 10
	group := get_slide(i, max, spingMap)
	if group != nil {
		groups = append(groups, group)
	}
	return groups, un
}
func get_slide(left, right int, data map[string]lib.Sping) []string {
	cc := 0
	for _, v := range data {
		if v.Ttl <= right && v.Ttl >= left {
			cc++
		}
	}
	if cc == 0 {
		return nil
	}
	//wtr := bufio.NewWriter(wfi)
	group := []string{}
	for k, v := range data {
		if left <= v.Ttl && right >= v.Ttl {
			group = append(group, k)
		}
	}
	return group
}
