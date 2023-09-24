/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-28:10
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package data4sk

import (
	"alias_article/lib"
	"alias_article/utils"
	"bufio"
	"os"
	"strconv"
)

type FieldChoice struct {
	LengthDiff          bool
	DirectDiff          bool
	TraceSameFactor     bool
	RttGap              bool
	ReplyTtl            bool
	TopDomainConsist    bool
	SecondDomainConsist bool
	SubDomainConsist    bool
	DomainDistance      bool
	IPDistance          bool
}

type DataGener struct {
	Alias       bool
	Field       FieldChoice
	SmarkPath   string
	PairsPath   string
	DomainPath  string
	SpingPath   string
	SavePath    string
	pairs       []lib.Pair
	ipTraceMap  map[string]lib.RawTrace
	ipSpingMap  map[string]lib.Sping
	ipDomainMap map[string]string
	WriteHead   bool

	DomainDistanceMap          map[string]int
	pairLengthDiffMap          map[string]int
	pairDirectDiffMap          map[string]int
	pairTraceFactorMap         map[string]int
	pairRttGapMap              map[string]int
	pairReplyTtlMap            map[string]int
	pairTopDomainConsistMap    map[string]bool
	pairSecondDomainConsistMap map[string]bool
	pairSubDomainConsistMap    map[string]bool
	NetSecMap                  map[string]int
}

func (g *DataGener) Run() {
	g.dataParse()
	g.dataWrite()
}

// write charactor vectors into specified file
func (g *DataGener) dataWrite() {
	wfi, err := os.OpenFile(g.SavePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	if g.WriteHead {
		g.writeHead(wtr)
	}
	for index, pair := range g.pairs {
		g.writePair(pair, wtr, index)
	}
	wtr.Flush()
}
func (g *DataGener) writeHead(wtr *bufio.Writer) {
	str := ""
	if g.Field.LengthDiff {
		str += "len_diff,"
	}
	if g.Field.DirectDiff {
		str += "dir_diff,"
	}
	if g.Field.TraceSameFactor {
		str += "same_factor,"
	}
	if g.Field.RttGap {
		str += "rtt_gap,"
	}
	if g.Field.ReplyTtl {
		str += "reply_ttl_gap,"
	}
	if g.Field.TopDomainConsist {
		str += "top_domain_consist,"
	}
	if g.Field.SecondDomainConsist {
		str += "second_domain_consist,"
	}
	if g.Field.SubDomainConsist {
		str += "sub_domain_consist,"
	}
	if g.Field.DomainDistance {
		str += "domain_distance,"
	}
	if g.Field.IPDistance {
		str += "ip_distance,"
	}
	str += "alias"
	// str = str[:len(str)-1]
	str += "\n"
	wtr.WriteString(str)
}

// write vector of one pair into file
func (g *DataGener) writePair(pair lib.Pair, wtr *bufio.Writer, index int) {
	// str := strconv.Itoa(index) + ","
	str := ""
	if g.Field.LengthDiff {
		if v, ok := g.pairLengthDiffMap[pair.ID()]; ok {
			str += strconv.Itoa(v)
		} else {
			str += "-1"
		}
		str += ","
	}
	if g.Field.DirectDiff {
		if v, ok := g.pairDirectDiffMap[pair.ID()]; ok {
			str += strconv.Itoa(v)
		} else {
			str += "-1"
		}
		str += ","
	}
	if g.Field.TraceSameFactor {
		if v, ok := g.pairTraceFactorMap[pair.ID()]; ok {
			str += strconv.Itoa(v)
		} else {
			str += "-1"
		}
		str += ","
	}
	if g.Field.RttGap {
		if v, ok := g.pairRttGapMap[pair.ID()]; ok {
			str += strconv.Itoa(v)
		} else {
			str += "-1"
		}
		str += ","
	}
	if g.Field.ReplyTtl {
		if v, ok := g.pairReplyTtlMap[pair.ID()]; ok {
			str += strconv.Itoa(v)
		} else {
			str += "-1"
		}
		str += ","
	}
	if g.Field.TopDomainConsist {
		if v, ok := g.pairTopDomainConsistMap[pair.ID()]; ok {
			if v {
				str += "1"
			} else {
				str += "0"
			}
		} else {
			str += "-1"
		}
		str += ","
	}
	if g.Field.SecondDomainConsist {
		if v, ok := g.pairSecondDomainConsistMap[pair.ID()]; ok {
			if v {
				str += "1"
			} else {
				str += "0"
			}
		} else {
			str += "-1"
		}
		str += ","
	}
	if g.Field.SubDomainConsist {
		if v, ok := g.pairSubDomainConsistMap[pair.ID()]; ok {
			if v {
				str += "1"
			} else {
				str += "0"
			}
		} else {
			str += "-1"
		}
		str += ","
	}
	if g.Field.DomainDistance {
		if v, ok := g.DomainDistanceMap[pair.ID()]; ok {
			str += strconv.Itoa(v)
		} else {
			str += "-1"
		}
		str += ","
	}
	if g.Field.IPDistance {
		str += strconv.Itoa(g.NetSecMap[pair.ID()])
		str += ","
	}
	// str = str[:len(str)-1]
	if g.Alias {
		str += "1"
	} else {
		str += "0"
	}
	str += "\n"
	wtr.WriteString(str)
}

// parse charactor vectors of all pairs first
func (g *DataGener) dataParse() {
	g.LoadData()
	if g.Field.LengthDiff {
		g.pairLengthDiff()
	}
	if g.Field.DirectDiff {
		g.pairDirectDiff()
	}
	if g.Field.TraceSameFactor {
		g.pairFactor()
	}
	if g.Field.RttGap {
		g.pairRttGap()
	}
	if g.Field.ReplyTtl {
		g.pairReplyTtlGap()
	}
	if g.Field.TopDomainConsist {
		g.pairTopDomainConsist()
	}
	if g.Field.SecondDomainConsist {
		g.pairSecondDomainConsist()
	}
	if g.Field.SubDomainConsist {
		g.pairSubDomainConsist()
	}
}

// calculates 'Difference Value of Path Length' of all pairs
func (g *DataGener) pairLengthDiff() {
	d := map[string]int{}
	for _, pair := range g.pairs {
		if _, ok := g.ipTraceMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := g.ipTraceMap[pair.IPB]; !ok {
			continue
		}
		traceA := g.ipTraceMap[pair.IPA]
		traceB := g.ipTraceMap[pair.IPB]
		if utils.ValidTracePair(traceA, traceB) {
			abs := difLength(traceA, traceB)
			d[pair.ID()] = abs
		}
	}
	g.pairLengthDiffMap = d
}

// calculate 'Difference Value of Path Direction' of all pairs
func (g *DataGener) pairDirectDiff() {
	d := map[string]int{}
	for _, pair := range g.pairs {
		if _, ok := g.ipTraceMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := g.ipTraceMap[pair.IPB]; !ok {
			continue
		}
		traceA := g.ipTraceMap[pair.IPA]
		traceB := g.ipTraceMap[pair.IPB]
		if utils.ValidTracePair(traceA, traceB) {
			abs := difDirect(traceA, traceB)
			d[pair.ID()] = abs
		}
	}
	g.pairDirectDiffMap = d
}

// calculates 'Path Similarity Coefficient' for all pairs
func (g *DataGener) pairFactor() {
	d := map[string]int{}
	for _, pair := range g.pairs {
		if _, ok := g.ipTraceMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := g.ipTraceMap[pair.IPB]; !ok {
			continue
		}
		traceA := g.ipTraceMap[pair.IPA]
		traceB := g.ipTraceMap[pair.IPB]
		if utils.ValidTracePair(traceA, traceB) {
			factor := factor(traceA, traceB)
			d[pair.ID()] = factor
		}
	}
	g.pairTraceFactorMap = d
}

// calculates 'Difference Value of Relative round-trip time' of all pairs
func (g *DataGener) pairRttGap() {
	d := map[string]int{}
	for _, pair := range g.pairs {
		if _, ok := g.ipTraceMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := g.ipTraceMap[pair.IPB]; !ok {
			continue
		}
		traceA := g.ipTraceMap[pair.IPA]
		traceB := g.ipTraceMap[pair.IPB]
		if utils.ValidTracePair(traceA, traceB) {
			gap := rttGap(traceA, traceB)
			d[pair.ID()] = gap
		}
	}
	g.pairRttGapMap = d
}

// calculates 'Difference Value of Reply TTL' for all pairs
func (g *DataGener) pairReplyTtlGap() {
	ttlGapMap := map[string]int{}
	validCnt := 0
	for _, pair := range g.pairs {
		if _, ok := g.ipSpingMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := g.ipSpingMap[pair.IPB]; !ok {
			continue
		}
		validCnt++
		g := utils.Abs(g.ipSpingMap[pair.IPA].Ttl - g.ipSpingMap[pair.IPB].Ttl)
		if g >= 15 {
			g = 15
		}
		ttlGapMap[pair.ID()] = g
	}
	g.pairReplyTtlMap = ttlGapMap
}

// calculates 'Top-Level Domain Name Consistency' for all pairs
func (g *DataGener) pairTopDomainConsist() {
	g.pairTopDomainConsistMap = map[string]bool{}
	for _, pair := range g.pairs {
		if _, ok := g.ipDomainMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := g.ipDomainMap[pair.IPB]; !ok {
			continue
		}
		consist := topConsist(g.ipDomainMap[pair.IPA], g.ipDomainMap[pair.IPB])
		g.pairTopDomainConsistMap[pair.ID()] = consist
	}
}

// calculates 'Second-Level Domain Name Consistency' for all pairs
func (g *DataGener) pairSecondDomainConsist() {
	g.pairSecondDomainConsistMap = map[string]bool{}
	for _, pair := range g.pairs {
		if _, ok := g.ipDomainMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := g.ipDomainMap[pair.IPB]; !ok {
			continue
		}
		valid, consist := secondConsist(g.ipDomainMap[pair.IPA], g.ipDomainMap[pair.IPB])
		if valid {
			g.pairSecondDomainConsistMap[pair.ID()] = consist
		}
	}
}

// calculates 'Third-Level Domain Name Consistency' for all pairs
func (g *DataGener) pairSubDomainConsist() {
	g.pairSubDomainConsistMap = map[string]bool{}
	for _, pair := range g.pairs {
		if _, ok := g.ipDomainMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := g.ipDomainMap[pair.IPB]; !ok {
			continue
		}
		valid, consist := subConsist(g.ipDomainMap[pair.IPA], g.ipDomainMap[pair.IPB])
		if valid {
			g.pairSubDomainConsistMap[pair.ID()] = consist
		}
	}
}

// calculates 'Character Edit Distance of the Domain Name' for all pairs
func (g *DataGener) DomainDistance() {
	g.DomainDistanceMap = map[string]int{}
	for _, pair := range g.pairs {
		if _, ok := g.ipDomainMap[pair.IPA]; ok {
			if _, ok := g.ipDomainMap[pair.IPB]; ok {
				g.DomainDistanceMap[pair.ID()] = domainDistance(g.ipDomainMap[pair.IPA], g.ipDomainMap[pair.IPB])
			}
		}
	}
}

// calculates 'Spatial Distance of the IP pair' for all pairs
func (g *DataGener) NetSec() {
	g.NetSecMap = map[string]int{}
	for _, pair := range g.pairs {
		g.NetSecMap[pair.ID()] = SecGap(pair.IPA, pair.IPB)
	}
}

// this method is for pipeline mode, which generates charactor vectors one by one
func (g *DataGener) OnePair(ipA, ipB string) []int {
	res := []int{}
	if g.Field.LengthDiff {
		if _, ok := g.ipTraceMap[ipA]; !ok {
			res = append(res, -1)
		}
		if _, ok := g.ipTraceMap[ipB]; !ok {
			res = append(res, -1)
		}
		traceA := g.ipTraceMap[ipA]
		traceB := g.ipTraceMap[ipB]
		if utils.ValidTracePair(traceA, traceB) {
			abs := difLength(traceA, traceB)
			res = append(res, abs)
		}
	}
	if g.Field.DirectDiff {
		if _, ok := g.ipTraceMap[ipA]; !ok {
			res = append(res, -1)
		}
		if _, ok := g.ipTraceMap[ipB]; !ok {
			res = append(res, -1)
		}
		traceA := g.ipTraceMap[ipA]
		traceB := g.ipTraceMap[ipB]
		if utils.ValidTracePair(traceA, traceB) {
			abs := difDirect(traceA, traceB)
			res = append(res, abs)
		}
	}
	if g.Field.TraceSameFactor {
		if _, ok := g.ipTraceMap[ipA]; !ok {
			res = append(res, -1)
		}
		if _, ok := g.ipTraceMap[ipB]; !ok {
			res = append(res, -1)
		}
		traceA := g.ipTraceMap[ipA]
		traceB := g.ipTraceMap[ipB]
		if utils.ValidTracePair(traceA, traceB) {
			factor := factor(traceA, traceB)
			res = append(res, factor)
		}
	}
	if g.Field.RttGap {
		if _, ok := g.ipTraceMap[ipA]; !ok {
			res = append(res, -1)
		}
		if _, ok := g.ipTraceMap[ipB]; !ok {
			res = append(res, -1)
		}
		traceA := g.ipTraceMap[ipA]
		traceB := g.ipTraceMap[ipB]
		if utils.ValidTracePair(traceA, traceB) {
			gap := rttGap(traceA, traceB)
			res = append(res, gap)
		}
	}
	if g.Field.ReplyTtl {
		if _, ok := g.ipSpingMap[ipA]; !ok {
			res = append(res, -1)
		}
		if _, ok := g.ipSpingMap[ipB]; !ok {
			res = append(res, -1)
		}
		res = append(res, utils.Abs(g.ipSpingMap[ipA].Ttl-g.ipSpingMap[ipB].Ttl))
	}
	if g.Field.TopDomainConsist {
		if _, ok := g.ipDomainMap[ipA]; !ok {
			res = append(res, -1)
		}
		if _, ok := g.ipDomainMap[ipB]; !ok {
			res = append(res, -1)
		}
		consist := topConsist(g.ipDomainMap[ipA], g.ipDomainMap[ipB])
		if consist {
			res = append(res, 1)
		} else {
			res = append(res, 0)
		}
	}
	if g.Field.SecondDomainConsist {
		if _, ok := g.ipDomainMap[ipA]; !ok {
			res = append(res, -1)
		}
		if _, ok := g.ipDomainMap[ipB]; !ok {
			res = append(res, -1)
		}
		valid, consist := secondConsist(g.ipDomainMap[ipA], g.ipDomainMap[ipB])
		if valid {
			if consist {
				res = append(res, 1)
			} else {
				res = append(res, -1)
			}
		} else {
			res = append(res, -1)
		}
	}

	if g.Field.SubDomainConsist {
		if _, ok := g.ipDomainMap[ipA]; !ok {
			res = append(res, -1)
		}
		if _, ok := g.ipDomainMap[ipB]; !ok {
			res = append(res, -1)
		}
		valid, consist := subConsist(g.ipDomainMap[ipA], g.ipDomainMap[ipB])
		if valid {
			if consist {
				res = append(res, 1)
			} else {
				res = append(res, 0)
			}
		} else {
			res = append(res, -1)
		}
	}
	return res
}
