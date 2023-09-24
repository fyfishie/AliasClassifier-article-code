package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"whoareu/lib"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func main() {
	whoisQuery("./input.txt", "output.txt", 10000)
}
func whoisQuery(ipPath string, wtpath string, runtineNum int) {
	inChan := make(chan string, 1000)
	outChan := make(chan lib.WhoisRes, 1000)
	doneChan := make(chan struct{})
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	go func() {
		rfi, err := os.Open(ipPath)
		if err != nil {
			panic(err.Error())
		}
		defer rfi.Close()
		scaner := bufio.NewScanner(rfi)
		cnt := 0
		for scaner.Scan() {
			cnt++
			if cnt%10000 == 0 {
				fmt.Printf("cnt: %v\n", cnt)
			}
			inChan <- scaner.Text()
		}
		close(inChan)
	}()
	for i := 0; i < runtineNum; i++ {
		go multiWhoisInfo(inChan, outChan, doneChan)
	}
	doneCnt := 0
	for {
		select {
		case info := <-outChan:
			bs, _ := json.Marshal(info)
			wtr.Write(bs)
			wtr.WriteString("\n")
		case <-doneChan:
			doneCnt++
			if doneCnt == runtineNum {
				wtr.Flush()
				return
			}
		}
	}
}

func multiWhoisInfo(ipIn chan string, outChan chan lib.WhoisRes, doneChan chan struct{}) {
	for ip := range ipIn {
		info, err := whoisInfo(ip)
		if err != nil {
			continue
		}
		outChan <- lib.WhoisRes{
			IP:                ip,
			WhoisInfoInstance: info2Instance(info),
		}
	}
	doneChan <- struct{}{}
}
func info2Instance(info whoisparser.WhoisInfo) lib.WhoIsInfoInstance {
	inst := lib.WhoIsInfoInstance{}
	if info.Administrative != nil {
		inst.Administrative = lib.Contact(*info.Administrative)
	}
	if info.Registrar != nil {
		inst.Registrar = lib.Contact(*info.Registrar)
	}
	if info.Registrant != nil {
		inst.Registrant = lib.Contact(*info.Registrant)
	}
	if info.Technical != nil {
		inst.Technical = lib.Contact(*info.Technical)
	}
	if info.Billing != nil {
		inst.Billing = lib.Contact(*info.Billing)
	}
	inst.Domain = lib.Domain{
		ID:          info.Domain.ID,
		Domain:      info.Domain.Domain,
		Punycode:    info.Domain.Punycode,
		Name:        info.Domain.Name,
		Extension:   info.Domain.Extension,
		WhoisServer: info.Domain.WhoisServer,
		Status:      info.Domain.Status,
		NameServers: info.Domain.NameServers,
		DNSSec:      info.Domain.DNSSec,
		CreatedDate: info.Domain.CreatedDate,
		// CreatedDateInTime:    *info.Domain.CreatedDateInTime,
		UpdatedDate: info.Domain.UpdatedDate,
		// UpdatedDateInTime:    *info.Domain.UpdatedDateInTime,
		ExpirationDate: info.Domain.ExpirationDate,
		// ExpirationDateInTime: *info.Domain.ExpirationDateInTime,
	}
	// if info.Domain.CreatedDateInTime != nil {
	// 	inst.Domain.UpdatedDateInTime = *info.Domain.UpdatedDateInTime
	// }
	// if info.Domain.UpdatedDateInTime != nil {
	// 	inst.Domain.UpdatedDateInTime = *info.Domain.UpdatedDateInTime
	// }
	// if info.Domain.ExpirationDateInTime != nil {
	// 	inst.Domain.ExpirationDateInTime = *info.Domain.ExpirationDateInTime
	// }
	return inst
}
func whoisInfo(ip string) (whoisparser.WhoisInfo, error) {
	domains := domains(ip)
	if domains == nil {
		return whoisparser.WhoisInfo{}, errors.New("no domain")
	}
	res, err := whois.Whois(domains[0])
	if err != nil {
		return whoisparser.WhoisInfo{}, err
	}
	pres, err := whoisparser.Parse(res)
	return pres, err
}
func domains(ip string) []string {
	ptr, err := net.LookupAddr(ip)
	if err != nil {
		return nil
	}
	domains := []string{}
	for _, rdnsRes := range ptr {
		if strings.HasPrefix(rdnsRes, "lookup") {
			continue
		}
		domains = append(domains, strings.TrimSuffix(rdnsRes, "."))
	}
	return domains
}
