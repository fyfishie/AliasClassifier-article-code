/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-21:16
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package filter

import (
	"alias_article/lib"
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// used to generate Ground-real datasets:
// by throwing pairs without same whois information away
func WhoisFilter(pairsPath string, wtpath string, whoisPath string) {
	whoisMap := map[string]lib.WhoIsInfoInstance{}
	whoisldr := store.NewLoader[lib.WhoisRes](whoisPath)
	err := whoisldr.Open()
	if err != nil {
		panic(err.Error())
	}
	for whoisldr.HasNext() {
		wrap := whoisldr.Next()
		whoisMap[wrap.IP] = wrap.WhoisInfoInstance
	}
	whoisldr.Close()
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	pairs, err := store.LoadAny[lib.Pair](pairsPath)
	esyerr.AutoPanic(err)
	filted := 0
	for _, pair := range pairs {
		if _, ok := whoisMap[pair.IPA]; !ok {
			bs, _ := json.Marshal(pair)
			wtr.Write(bs)
			wtr.WriteString("\n")
			continue
		}
		if _, ok := whoisMap[pair.IPB]; !ok {
			bs, _ := json.Marshal(pair)
			wtr.Write(bs)
			wtr.WriteString("\n")
			continue
		}
		if whoisInfoEqule(whoisMap[pair.IPA], whoisMap[pair.IPB]) {
			bs, _ := json.Marshal(pair)
			wtr.Write(bs)
			wtr.WriteString("\n")
			continue
		} else {
			filted++
		}
	}
	fmt.Printf("filted: %v\n", filted)
	esyerr.AutoPanic(wtr.Flush())
}
func whoisInfoEqule(a, b lib.WhoIsInfoInstance) bool {
	e := true
	e = e && a.Domain.UpdatedDate == b.Domain.UpdatedDate
	e = e && a.Domain.CreatedDate == b.Domain.CreatedDate
	e = e && a.Registrant.Organization == b.Registrant.Organization
	e = e && a.Registrant.Name == b.Registrant.Name
	e = e && a.Registrant.Country == b.Registrant.Country
	e = e && a.Registrant.City == b.Registrant.City
	e = e && a.Registrant.Email == b.Registrant.Email
	return e
}
