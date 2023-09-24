package character

import (
	"github.com/fyfishie/dorapock/store"
	"mlar/lib"
	"mlar/utils"
)

type MlarDataItem struct {
	Pair      lib.Pair
	RttSame   float64
	traceSame float64
}

// it generates data for machine learning program
type Gener struct {
	tracePaths  []string
	pairPath    string
	wtPath      string
	DataOut     chan MlarDataItem
	interips    map[string]struct{}
	ipTracesMap map[string][]lib.RawTrace
	ipRttsMap   map[string][]int
	pairs       []lib.Pair
	tl          int
}

// construct method
func NewGener(tracePaths []string, pairsPath string, wtPath string) *Gener {
	c := make(chan MlarDataItem, 1000)
	return &Gener{
		tracePaths:  tracePaths,
		pairPath:    pairsPath,
		wtPath:      wtPath,
		DataOut:     c,
		tl:          len(tracePaths),
		interips:    map[string]struct{}{},
		ipTracesMap: map[string][]lib.RawTrace{},
		ipRttsMap:   ipTracesMap,
	}
}

// Run!!!
func (g *Gener) Run() (map[string]int, map[string]int) {
	g.ini()
	rttSames := map[string]int{}
	for _, pair := range g.pairs {
		r := rttfactor(g.ipRttsMap[pair.IPA], g.ipRttsMap[pair.IPB], g.tl)
		rttSames[pair.ID()] = r
	}
	traceFactors := TraceFactor(g.tracePaths, g.pairPath)
	return rttSames, traceFactors
}

// load data for anti-alias pair filt
func (g *Gener) ini() {
	pairs, err := store.LoadAny[lib.Pair](g.pairPath)
	if err != nil {
		panic(err)
	}
	g.pairs = pairs
	g.interips = map[string]struct{}{}
	for _, p := range pairs {
		g.interips[p.IPA] = struct{}{}
		g.interips[p.IPB] = struct{}{}
	}
	for ip, _ := range g.interips {
		ori := []int{}
		for i := 0; i < len(g.tracePaths); i++ {
			ori = append(ori, -1)
		}
		g.ipRttsMap[ip] = ori
	}
	for index, tracePath := range g.tracePaths {
		traces := utils.ValidTrace(tracePath)
		for _, trace := range traces {
			if _, ok := g.ipRttsMap[trace.Ip]; ok {
				//array := g.ipRttsMap[trace.Ip]
				//array[index] = trace.Results[len(trace.Results)-1].Rtt
				g.ipRttsMap[trace.Ip][index] = trace.Results[len(trace.Results)-1].Rtt / 1000000
				//g.ipRttsMap[trace.Ip] = array
			}
		}
	}
}
