/*
 * @Author: fyfishie
 * @Date: 2023-04-22:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-22:21
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package statistic

import (
	"alias_article/lib"
	"alias_article/utils"
	"fmt"
	"github.com/agnivade/levenshtein"
	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// statistic domain distance charactor
func DomainDistanceSta(domainPath string, pairPath string, wtpath string) {
	domains := loadDomain(domainPath)
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)

	validCnt := 0
	distanceMap := map[int]int{}
	for _, pair := range pairs {
		if _, ok := domains[pair.IPA]; !ok {
			continue
		}
		if _, ok := domains[pair.IPB]; !ok {
			continue
		}
		validCnt++
		distance := levenshtein.ComputeDistance(utils.ReverseDomain(domains[pair.IPA]), utils.ReverseDomain(domains[pair.IPB]))
		distanceMap[distance]++
	}
	WriteAccuDataFrom(distanceMap, wtpath, fmt.Sprintf("valid scale: %v\n", validCnt*100/len(pairs)))
}
