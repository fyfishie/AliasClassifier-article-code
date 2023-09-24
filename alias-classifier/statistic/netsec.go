/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-22:21
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package statistic

import (
	"alias_article/lib"
	"alias_article/utils"
	"math"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
	"github.com/fyfishie/ipop"
)

// statistic 'Spatial Distance of the IP pair' charactor
func NetSecSta(pairsPath, wtPath string) {
	pairs, err := store.LoadAny[lib.Pair](pairsPath)
	esyerr.AutoPanic(err)
	secMap := map[int]int{}
	for _, pair := range pairs {
		secMap[secGap(pair.IPA, pair.IPB)]++
	}
	WriteAccuDataFrom(secMap, wtPath, "")
}
func secGap(a, b string) int {
	ia := ipop.String2Int(a)
	ib := ipop.String2Int(b)
	gap := utils.Abs(ia - ib)
	lg := math.Log2(float64(gap))
	return int(math.Round(lg))
}
