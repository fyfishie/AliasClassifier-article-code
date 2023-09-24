package assess

import (
	"alias_article/lib"
	"fmt"
	"github.com/fyfishie/dorapock/store"
)

// calculate precious, recall and f1 score from file
func Assess_f1_pre_rec(p_pair_path, t_pair_path, f_pair_path string) {
	ppairs, _ := store.LoadAny[lib.Pair](p_pair_path)
	treemap := map[string]lib.Pair{}
	for _, p := range ppairs {
		treemap[p.ID()] = p
	}
	tpairs, err := store.LoadAny[lib.Pair](t_pair_path)
	if err != nil {
		panic(err)
	}
	tpairsMap := map[string]lib.Pair{}
	for _, p := range tpairs {
		tpairsMap[p.ID()] = p
	}
	fpairs, err := store.LoadAny[lib.Pair](f_pair_path)
	if err != nil {
		panic(err)
	}
	fpairsmap := map[string]lib.Pair{}
	for _, p := range fpairs {
		fpairsmap[p.ID()] = p
	}
	tpCnt := 0
	for _, p := range ppairs {
		if _, ok := tpairsMap[p.ID()]; ok {
			tpCnt++
		}
	}
	//tnCnt := len(tpairs) - tpCnt
	fpCnt := 0
	for _, p := range ppairs {
		if _, ok := fpairsmap[p.ID()]; ok {
			fpCnt++
		}
	}
	fnCnt := len(fpairs) - fpCnt
	pre := float64(tpCnt) / float64(tpCnt+fpCnt)
	recall := float64(tpCnt) / float64(tpCnt+fnCnt)
	f1 := 2 * pre * recall / (pre + recall)

	fmt.Printf("pre:%v\n", pre)
	fmt.Printf("recall:%v\n", recall)
	fmt.Printf("fi:%v\n", f1)
}
