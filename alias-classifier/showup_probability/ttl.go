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
	"fmt"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// calculates scale of ip who has ttl data in origin ip list
func TTLScale(spingPath, pairPath string) {
	spings, err := store.LoadAny[lib.Sping](spingPath)
	esyerr.AutoPanic(err)
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	spingMap := map[string]lib.Sping{}
	for _, sping := range spings {
		spingMap[sping.IP] = sping
	}
	ipMap := map[string]struct{}{}
	for _, pair := range pairs {
		ipMap[pair.IPA] = struct{}{}
		ipMap[pair.IPB] = struct{}{}
	}
	pairNum := 0
	for _, pair := range pairs {
		if _, ok := spingMap[pair.IPA]; ok {
			if _, ok := spingMap[pair.IPB]; ok {
				pairNum++
			}
		}
	}
	fmt.Printf("pair scale:%v\n", pairNum*100/len(pairs))
	fmt.Printf("ip scale:%v\n", len(spingMap)*100/len(ipMap))
}
