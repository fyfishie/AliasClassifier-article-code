/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-16:11
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

// calculates scale of ip of all pairs in origin ip list
func IPTotal(pairsPath string) {
	pairs, err := store.LoadAny[lib.Pair](pairsPath)
	esyerr.AutoPanic(err)
	m := map[string]struct{}{}
	for _, pair := range pairs {
		m[pair.IPA] = struct{}{}
		m[pair.IPB] = struct{}{}
	}
	fmt.Println(len(m))
}
