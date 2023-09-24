/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-07-20:16
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package data4sk

import (
	"alias_article/lib"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// load origin data before data parse
func (g *DataGener) LoadData() {
	traces, err := store.LoadAny[lib.RawTrace](g.SmarkPath)
	esyerr.AutoPanic(err)
	ipTraceMap := map[string]lib.RawTrace{}
	for _, trace := range traces {
		ipTraceMap[trace.Ip] = trace
	}
	g.ipTraceMap = ipTraceMap

	spings, err := store.LoadAny[lib.Sping](g.SpingPath)
	esyerr.AutoPanic(err)
	ipSpingMap := map[string]lib.Sping{}
	for _, sping := range spings {
		ipSpingMap[sping.IP] = sping
	}
	g.ipSpingMap = ipSpingMap

	domains := loadDomain(g.DomainPath)
	g.ipDomainMap = domains

	if g.PairsPath != "" {
		pairs, err := store.LoadAny[lib.Pair](g.PairsPath)
		esyerr.AutoPanic(err)
		g.pairs = pairs
	}
}

func loadDomain(path string) map[string]string {
	domains, err := store.LoadAny[lib.RDNSResItem](path)
	if err != nil {
		panic(err)
	}
	res := map[string]string{}
	for _, domain := range domains {
		if len(domain.Domains) > 0 {
			res[domain.IP] = domain.Domains[0]
		}
	}
	return res
}
