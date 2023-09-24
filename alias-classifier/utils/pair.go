/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-21:15
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package utils

import (
	"alias_article/lib"
	"bufio"
	"encoding/json"
	"os"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// make some anti-alias pairs
func MakeUnPair(nodePath string, wtpath string) {
	nodeslr := store.NewLoader[lib.Node](nodePath)
	esyerr.AutoPanic(nodeslr.Open())
	cnt := 0
	nodes := []lib.Node{}
	for nodeslr.HasNext() {
		cnt++
		cnt = cnt % 250
		if cnt == 0 {
			nodes = append(nodes, nodeslr.Next())
		}
	}
	allPairs := map[string]struct{}{}
	allIPMap := map[string]struct{}{}
	for _, node := range nodes {
		for i := 0; i < len(node.IPList); i++ {
			allIPMap[node.IPList[i]] = struct{}{}
			for j := i + 1; j < len(node.IPList); j++ {
				p := lib.Pair{IPA: node.IPList[i], IPB: node.IPList[j]}
				allPairs[p.ID()] = struct{}{}
			}
		}
	}
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	allip := []string{}
	for ip := range allIPMap {
		allip = append(allip, ip)
	}
	for i := 0; i < len(allip); i++ {
		for j := i + 1; j < len(allip); j++ {
			up := lib.Pair{IPA: allip[i], IPB: allip[j]}
			if _, ok := allPairs[up.ID()]; !ok {
				bs, _ := json.Marshal(up)
				wtr.Write(bs)
				wtr.WriteString("\n")
			}
		}
	}
	esyerr.AutoPanic(wtr.Flush())

}

// get all alias pairs from itdk data
func AllPair(nodePath string, wtpath string) {
	nodes, err := store.LoadAny[lib.Node](nodePath)
	if err != nil {
		PanicOnError(err)
	}
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	for _, node := range nodes {
		if len(node.IPList) != 0 {
			for i := 0; i < len(node.IPList); i++ {
				for j := i + 1; j < len(node.IPList); j++ {
					p := lib.Pair{IPA: node.IPList[i], IPB: node.IPList[j], ANode: node.NodeID, BNode: node.NodeID}
					bs, _ := json.Marshal(p)
					wtr.Write(bs)
					wtr.WriteString("\n")
				}
			}
		}
	}
	PanicOnError(wtr.Flush())
}
