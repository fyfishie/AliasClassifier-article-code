package aliascek

import (
	// "aliasParseMaster/filtermaker"
	"aliasParseMaster/graph"
	"aliasParseMaster/lib"
	"aliasParseMaster/mongoproxy"
	"aliasParseMaster/utils"
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/fyfishie/ipop"
	"github.com/sirupsen/logrus"
)

/*
 * @description: all fileds should be set manually
 * @return {*}
 */
type checker struct {
	//this field should be reset by washAndSet function, which make sure the data is consistent
	// TargetIPList []int
	consistIntIP []int
	checkTask    *lib.AliasCekTask
	//this field should be set by washAndSet function, which make sure the data is consistent
	ipSpingResMap map[int]map[lib.VPUID]lib.SpingWithoutUID
	vpSpingResMap map[lib.VPUID]map[int]lib.SpingWithoutUID
	ipInDB        map[int]struct{}
	tSets         [][]int
	finalBufWtr   *mongoproxy.BufferWriter
	andVPList     []lib.VPUID
	// AntiAliasMCA   lib.MongoCollectionAddr
	// PingResMCA     lib.MongoCollectionAddr
	// TSetMCA        lib.MongoCollectionAddr
	// MaybeAliasMCA  lib.MongoCollectionAddr
	// FinalResultMCA lib.MongoCollectionAddr
}

const TORM_T int = 0 //belongs to set T
const TORM_M int = 1 //belongs to set M
const TORM_X int = 2 //drop it
func NewChecker(aliasCheckTask *lib.AliasCekTask) *checker {
	return &checker{
		checkTask: aliasCheckTask,
	}
}
func (c *checker) secondInit() error {
	btr, err := mongoproxy.NewBufferWriterFromAddr(c.checkTask.AliasCekResAddr, 500)
	if err != nil {
		return err
	}
	c.finalBufWtr = btr
	return nil
}
func (c *checker) secondDestroy() error {
	err := c.finalBufWtr.Flush()
	if err != nil {
		return err
	}
	return c.finalBufWtr.Close()
}
func (c *checker) CheckStart() error {
	fmt.Println("start consistent ip")
	c.setConsistentIP()
	spingResult, err := c.readSpingResult()
	if err != nil {
		logrus.Errorf("error in read detect result, err = %v\n", err.Error())
		return err
	}
	fmt.Println("start setmap")
	err = c.setMap(spingResult)
	if err != nil {
		logrus.Errorf("error while set map, err = %v\n", err.Error())
		return err
	}
	fmt.Println("start classify")
	err = c.classify()
	if err != nil {
		logrus.Error("error in classify ip pairs, err = " + err.Error())
		return err
	}
	err = c.secondInit()
	if err != nil {
		logrus.Error("error in init second parse, err = " + err.Error())
	}
	fmt.Println("start second")
	err = c.Second()
	if err != nil {
		logrus.Error("error in running second parse, err = " + err.Error())
	}
	err = c.secondDestroy()
	if err != nil {
		logrus.Error("error in destory second parse, err = " + err.Error())
	}
	fmt.Println("done!")
	return err

}
func (c *checker) setMap(spingRes []lib.SpingWithUID) error {
	//find all VPs
	vpsMap := map[lib.VPUID]struct{}{}
	for _, uid := range c.andVPList {
		vpsMap[uid] = struct{}{}
	}
	ipItemMap := map[int]map[lib.VPUID]lib.SpingWithoutUID{}
	vpSpingResMap := map[lib.VPUID]map[int]lib.SpingWithoutUID{}
	for _, item := range spingRes {
		if !utils.IsInt(item.Ttl) {
			return errors.New("ttl of sping result is not int")
		}

		withoutUID := lib.SpingWithoutUID{

			IP:  ipop.String2Int(item.IP),
			Ttl: utils.MustInt(item.Ttl),
		}
		_, exist := ipItemMap[ipop.String2Int(item.IP)]
		if !exist {
			ipItemMap[ipop.String2Int(item.IP)] = map[lib.VPUID]lib.SpingWithoutUID{}
		}
		ipItemMap[ipop.String2Int(item.IP)][item.VPUID] = withoutUID

		if _, exist = vpSpingResMap[item.VPUID]; !exist {
			vpSpingResMap[item.VPUID] = map[int]lib.SpingWithoutUID{}
		}
		vpSpingResMap[item.VPUID][ipop.String2Int(item.IP)] = withoutUID
	}
	c.ipSpingResMap = ipItemMap
	c.vpSpingResMap = vpSpingResMap
	return nil
	// nextIP:
	// 	for ip, results := range ipItemMap {
	// 		for _, res := range results {
	// 			for _, vp := range c.andVPList {
	// 				if !utils.VPEquals(vp, res.VP) {
	// 					delete(ipItemMap, ip)
	// 					break nextIP
	// 				}
	// 			}
	// 		}
	// 	}
}
func (c *checker) readSpingResult() ([]lib.SpingWithUID, error) {
	proxy := mongoproxy.NewProxy(c.checkTask.PingResAddr)
	err := proxy.Connect()
	defer proxy.Disconnect()
	if err != nil {
		return nil, err
	}
	ires, err := proxy.ReadSpingWithUID(ipop.Ints2Strings(c.consistIntIP))
	if err != nil {
		logrus.Error("error in query ping result data from mongodb, err = " + err.Error())
		return nil, err
	}
	return ires, nil
}

/*
 * @description: classifies all ip pairs into T,M or drop it
 */
func (c *checker) classify() error {
	graph := graph.NewGraph()
	bufWriter, err := mongoproxy.NewBufferWriterFromAddr(c.checkTask.MaybeAliasMCA, 100)
	if err != nil {
		logrus.Error("error in make new bufwriter, err = " + err.Error())
		return err
	}
	proxy := mongoproxy.NewProxy(c.checkTask.AntiAliasResAddr)
	err = proxy.Connect()
	defer proxy.Disconnect()
	if err != nil {
		return err
	}
	readChan := proxy.RunReadPipeline(c.consistIntIP, 1000)
	err = c.qryIPInDB()
	if err != nil {
		logrus.Error("error in query ip in db, err = " + err.Error())
		return err
	}
	maybeWriteOutChan := make(chan lib.MaybeAlias, 1000)
	doneChan := make(chan struct{})
	go c.asyncWriter(maybeWriteOutChan, doneChan, bufWriter)
	antiMap := map[int]struct{}{}
	total := len(c.consistIntIP)
	start := time.Now().Unix()
	for i := 0; i < total; i++ {
		if i%1000 == 0 {
			fmt.Printf("i: %v\n", i*100/total)
			fmt.Println((time.Now().Unix() - start) / 60)
		}
		ipA := c.consistIntIP[i]
		maybeAlias := lib.MaybeAlias{IP: ipA, MaybeAliasSet: []int{}}
		_, hasAnti := c.ipInDB[ipA]
		if hasAnti {
			antiMap = <-readChan
		}
		for j := i + 1; j < len(c.consistIntIP); j++ {
			ipB := c.consistIntIP[j]
			if _, anti := antiMap[ipB]; anti {
				continue
			}
			torm := c.tORm(ipA, ipB)
			switch torm {
			case TORM_T:
				graph.AddEdge(ipA, ipB)
			case TORM_M:
				maybeAlias.MaybeAliasSet = append(maybeAlias.MaybeAliasSet, ipB)
			default:
			}
		}
		maybeWriteOutChan <- maybeAlias
	}
	close(maybeWriteOutChan)
	<-doneChan
	if err != nil {
		logrus.Error("error in flush maybeAlias set into db, err = " + err.Error())
	}
	c.tSets = graph.MaxConnectedSubGraph()
	return err
}
func (c *checker) asyncWriter(dataInChan chan lib.MaybeAlias, doneChan chan struct{}, bufwtr *mongoproxy.BufferWriter) {
	for data := range dataInChan {
		err := bufwtr.WriteOne(mongoproxy.MaybeAliasSetModel(data))
		if err != nil {
			logrus.Error("error in write maybeAlias set into db, err = " + err.Error())
		}
	}
	err := bufwtr.Flush()
	if err != nil {
		logrus.Error("error in flush maybeAlias set into db, err = " + err.Error())
	}
	doneChan <- struct{}{}
}

// pair should be in set T or set M?
// check defination of values
func (c *checker) tORm(ipA, ipB int) int {
	or := false
	and := true
	for _, spingMap := range c.vpSpingResMap {
		//I don not want to name it, i is good
		i := (spingMap[ipA].Ttl == spingMap[ipB].Ttl)
		or = or || i
		and = and && i
	}
	if and {
		return TORM_T
	} else if or {
		return TORM_M
	} else {
		return TORM_X
	}
}

// make sure the spingResult field data is consistent
func (c *checker) setConsistentIP() error {
	p := mongoproxy.NewProxy(c.checkTask.PingResAddr)
	err := p.Connect()
	defer p.Disconnect()
	if err != nil {
		return err
	}
	consistIP, andVP, err := p.ConsistenceIP()
	if err != nil {
		logrus.Errorf("error while get consist ip, err = %v\n", err.Error())
		return err
	}
	consistentIntIP := []int{}
	for _, ip := range consistIP {
		consistentIntIP = append(consistentIntIP, ipop.String2Int(ip))
	}
	sort.Ints(consistentIntIP)
	c.consistIntIP = consistentIntIP
	c.andVPList = andVP
	return nil
}

func conv2SpingWithUID(ires []interface{}) []lib.SpingWithUID {
	res := []lib.SpingWithUID{}
	for _, i := range ires {
		res = append(res, i.(lib.SpingWithUID))
	}
	return res
}

func (c *checker) qryIPInDB() error {
	proxy := mongoproxy.NewProxy(c.checkTask.AntiAliasResAddr)
	err := proxy.Connect()
	if err != nil {
		return err
	}
	ires, err := proxy.QueryIPLeaderinDB(c.consistIntIP, context.TODO())
	if err != nil {
		return err
	}
	c.ipInDB = ires
	proxy.Disconnect()
	return nil
}

func (c *checker) Second() error {
	//record which set IPi belongs to
	tIndex := map[int][]int{}
	for i := 0; i < len(c.tSets); i++ {
		TC := c.tSets[i]
		for _, tIP := range TC {
			tIndex[tIP] = TC
		}
	}
	mReader, err := mongoproxy.NewMPairsBufReader(c.checkTask.MaybeAliasMCA, 1000)
	if err != nil {
		return err
	}
	count := 0
	total := len(tIndex)
	for {
		count++
		if count%1000 == 0 {
			fmt.Printf("count_scale: %v\n", count*100/total)
		}
		m, err := mReader.NextOne()
		if err != nil {
			break
		}
		TC, hasTC := tIndex[m.IP]
		if !hasTC {
			continue
		}

		//format of IP_j and PC_j is consistent with that in patent article
		IP_js := []int{}
		for _, IP_j := range m.MaybeAliasSet {
			PC_j := append(TC, IP_j)
			if c.checkSecond(PC_j) {
				IP_js = append(IP_js, IP_j)
			}
		}
		c.writeOutAlias(m.IP, IP_js)
	}
	return nil
}
func (c *checker) checkSecond(PC_j []int) bool {
	k := len(PC_j)
	a := k * (k - 1) / 2
	r := c.r(PC_j)
	if a < 1 || r < 0 {
		return false
	}
	v_ := c.v_(float64(a), float64(r))
	vcnt := 0
	for i, ipA := range PC_j {
		for _, ipB := range PC_j[i+1:] {
			for _, list := range c.vpSpingResMap {
				if list[ipA].Ttl == list[ipB].Ttl {
					vcnt++
				}
			}
		}
	}
	return float64(vcnt) >= v_
}

func (c *checker) writeOutAlias(ipA int, ipBs []int) error {
	return c.finalBufWtr.WriteOne(mongoproxy.SureAliasModel(lib.SureAlias{IP: ipA, SureAliasSet: ipBs}))
}

func (c *checker) r(PC_j []int) int {
	r := 0
	for _, ip := range PC_j {
		if m, ok := c.ipSpingResMap[ip]; ok {
			for _, v := range m {
				if v.Ttl > r {
					r = v.Ttl
				}
			}
		}
	}
	return r
}

// r should be greater than 0 and a shoud be greater than 1
func (c *checker) v_(a, r float64) float64 {
	lna := math.Log(a)
	lna_1 := math.Log(a - 1)
	lnx := lna - math.Log(lna-lna_1)
	lnr := math.Log(r)
	return math.Ceil(lnx / lnr)
}
