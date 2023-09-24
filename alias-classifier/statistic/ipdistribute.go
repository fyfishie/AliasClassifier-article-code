package statistic

import (
	"alias_article/geocode"
	"alias_article/lib"
	"fmt"
	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

func IpDistribute() {
	pairs, err := store.LoadAny[lib.Pair]("./db/pairs/all_pairs-host-port-whois.json")
	esyerr.AutoPanic(err)
	ipMap := map[string]struct{}{}
	for _, pair := range pairs {
		ipMap[pair.IPA] = struct{}{}
		ipMap[pair.IPB] = struct{}{}
	}
	ipList := []string{}
	for k, _ := range ipMap {
		ipList = append(ipList, k)
	}
	searcher := geocode.NewXdbSearcher()
	searcher.SearchInit("./geocode/ip2location.xdb")
	res, _ := searcher.GetCountrysIDByXdb(ipList)
	cityMap := map[string]struct{}{}
	for _, v := range res {
		cityMap[v] = struct{}{}
	}
	if _, ok := cityMap[""]; ok {
		fmt.Println("?")
	}
	fmt.Println(len(cityMap))
}
