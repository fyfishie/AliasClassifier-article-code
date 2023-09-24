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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// statistic domain consistence charactor
func DomainConsistenceSta(domainPath string, pairsPath string, wtPath string) {
	pairs, err := store.LoadAny[lib.Pair](pairsPath)
	esyerr.AutoPanic(err)
	domains := loadDomain(domainPath)

	pairValidCnt := 0
	topConsistCnt := 0
	secondValidCnt := 0
	secondConsistCnt := 0
	subValidCnt := 0
	subConsistCnt := 0
	for _, pair := range pairs {
		if _, ok := domains[pair.IPA]; !ok {
			continue
		}
		if _, ok := domains[pair.IPB]; !ok {
			continue
		}
		pairValidCnt++
		if topConsist(domains[pair.IPA], domains[pair.IPB]) {
			topConsistCnt++
		}
		secValid, secConsist := secondConsist(domains[pair.IPA], domains[pair.IPB])
		if secValid {
			secondValidCnt++
			if secConsist {
				secondConsistCnt++
			}
		}
		subValid, subConsist := subConsist(domains[pair.IPA], domains[pair.IPB])
		if subValid {
			subValidCnt++
			if subConsist {
				subConsistCnt++
			}
		}
	}
	wfi, err := os.OpenFile(wtPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	wtr.WriteString(fmt.Sprintf("top consist:\nvalid scale: %v\n", pairValidCnt*100/len(pairs)))
	wtr.WriteString(fmt.Sprintf("consist scale: %v\n\n", topConsistCnt*100/pairValidCnt))
	wtr.WriteString(fmt.Sprintf("second consist:\nvalid scale: %v\n", secondValidCnt*100/len(pairs)))
	wtr.WriteString(fmt.Sprintf("consist scale: %v\n\n", secondConsistCnt*100/secondValidCnt))
	wtr.WriteString(fmt.Sprintf("sub valid scale: %v\n", subValidCnt*100/len(pairs)))
	wtr.WriteString(fmt.Sprintf("sub consist scale: %v\n\n", subConsistCnt*100/subValidCnt))
	wtr.Flush()
}

func topConsist(a, b string) bool {
	ssA := strings.Split(a, ".")
	ssB := strings.Split(b, ".")
	return ssA[len(ssA)-1] == ssB[len(ssB)-1]
}
func secondConsist(a, b string) (valid bool, consist bool) {
	ssA := strings.Split(a, ".")
	ssB := strings.Split(b, ".")
	if len(ssA) < 2 || len(ssB) < 2 {
		return false, false
	}
	return true, ssA[len(ssA)-2] == ssB[len(ssB)-2]
}
func subConsist(a, b string) (valid, consist bool) {
	ssA := strings.Split(a, ".")
	ssB := strings.Split(b, ".")
	if len(ssA) < 3 || len(ssB) < 3 {
		return false, false
	}
	return true, ssA[len(ssA)-3] == ssB[len(ssB)-3]
}

// map[ip]domain
func loadDomain(path string) map[string]string {
	rfi, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer rfi.Close()
	scaner := bufio.NewScanner(rfi)
	res := map[string]string{}
	for scaner.Scan() {
		line := scaner.Text()
		ss := strings.Split(line, "\t")
		res[ss[1]] = ss[2]
	}
	return res
}
