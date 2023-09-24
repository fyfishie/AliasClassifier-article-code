/*
 * @Author: fyfishie
 * @Date: 2023-04-26:19
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-22:16
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package statistic

import (
	"alias_article/lib"
	"alias_article/utils"
	"fmt"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// calculates 'Difference Value of Reply TTL ' charactor
func TTLScale(spingPath string, pairsPath string, wtpath string) {
	spings, err := store.LoadAny[lib.Sping](spingPath)
	esyerr.AutoPanic(err)
	pairs, err := store.LoadAny[lib.Pair](pairsPath)
	esyerr.AutoPanic(err)

	ipSpingMap := map[string]lib.Sping{}
	for _, s := range spings {
		ipSpingMap[s.IP] = s
	}

	ttlGapMap := map[int]int{}
	validCnt := 0
	for _, pair := range pairs {
		if _, ok := ipSpingMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := ipSpingMap[pair.IPB]; !ok {
			continue
		}
		validCnt++
		ttlGapMap[utils.Abs(ipSpingMap[pair.IPA].Ttl-ipSpingMap[pair.IPB].Ttl)]++
	}
	WriteAccuDataFrom(ttlGapMap, wtpath, fmt.Sprintf("valid scale: %v\n", validCnt*100/len(pairs)))
}
