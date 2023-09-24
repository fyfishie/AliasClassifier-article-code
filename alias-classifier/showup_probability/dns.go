/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-28:12
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package showup_probability

import (
	"alias_article/lib"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// calculates scale of ip who has dns data in origin ip list
func DNSScale(dnspath string, pairPath string) {
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	ipMap := map[string]struct{}{}
	for _, pair := range pairs {
		ipMap[pair.IPA] = struct{}{}
		ipMap[pair.IPB] = struct{}{}
	}
	domains := loadDomain(dnspath)
	pairNum := 0
	for _, pair := range pairs {
		if _, ok := domains[pair.IPA]; ok {
			if _, ok := domains[pair.IPB]; ok {
				pairNum++
			}
		}
	}
	ipNum := 0
	for ip := range ipMap {
		if _, ok := domains[ip]; ok {
			ipNum++
		}
	}
	fmt.Printf("pair scale:%v\n", pairNum*100/len(pairs))
	fmt.Printf("ip scale:%v\n", ipNum*100/len(ipMap))
}

func loadDomain(path string) map[string]string {
	rfi, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer rfi.Close()
	scaner := bufio.NewScanner(rfi)
	res := map[string]string{}
	for scaner.Scan() {
		line := scaner.Text()
		ss := strings.Split(line, "\t")
		res[ss[1]] = ss[2]
	}
	return res
}
