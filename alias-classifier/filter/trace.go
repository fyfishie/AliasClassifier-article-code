package filter

import (
	"alias_article/lib"
	"bufio"
	"encoding/json"
	"github.com/fyfishie/dorapock/store"
	"os"
)

func MakeTrace(ipPath string, OriginTrace string, wtpath string) {
	rfi, err := os.OpenFile(ipPath, os.O_RDONLY, 0000)
	if err != nil {
		panic(err)
	}
	scaner := bufio.NewScanner(rfi)
	allIP := map[string]struct{}{}
	for scaner.Scan() {
		allIP[scaner.Text()] = struct{}{}
	}
	ldr := store.NewLoader[lib.RawTrace](OriginTrace)
	err = ldr.Open()
	if err != nil {
		panic(err)
	}
	wfi, err := os.OpenFile(wtpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	for ldr.HasNext() {
		t := ldr.Next()
		if _, ok := allIP[t.Ip]; ok {
			bs, _ := json.Marshal(t)
			wfi.Write(bs)
			wfi.WriteString("\n")
		}
	}
	wfi.Close()
}
func MakeSping(ipPath string, OriginSping string, wtpath string) {
	rfi, err := os.OpenFile(ipPath, os.O_RDONLY, 0000)
	if err != nil {
		panic(err)
	}
	scaner := bufio.NewScanner(rfi)
	allIP := map[string]struct{}{}
	for scaner.Scan() {
		allIP[scaner.Text()] = struct{}{}
	}
	spings, err := store.LoadAny[lib.Sping](OriginSping)
	if err != nil {
		panic(err)
	}
	wfi, err := os.OpenFile(wtpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	for _, t := range spings {
		if _, ok := allIP[t.IP]; ok {
			bs, _ := json.Marshal(t)
			wfi.Write(bs)
			wfi.WriteString("\n")
		}
	}
	wfi.Close()
}
