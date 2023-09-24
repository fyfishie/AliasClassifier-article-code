/*
 * @Author: fyfishie
 * @Date: 2023-05-09:21
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-13:11
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package geocode

import (
	"aliasParseMaster/lib"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
)

type XdbSearcher struct {
	cBuff    []byte
	searcher *xdb.Searcher
}

func NewXdbSearcher() *XdbSearcher {
	return &XdbSearcher{}
}
func (s *XdbSearcher) SearchInit(dbPath string) error {
	s.cBuff, _ = xdb.LoadContentFromFile(dbPath)
	xdbs, err := xdb.NewWithBuffer(s.cBuff)
	if err != nil {
		logrus.Errorf("error while init searcher, err = %v\n", err.Error())
		return err
	}
	s.searcher = xdbs
	return nil
}
func (s *XdbSearcher) QueryIP(ip string) (lib.Location, error) {
	region, _ := s.searcher.SearchByStr(ip)
	ss := strings.Split(region, "|")
	if ss[0] == "0" {
		ss[0] = ""
	}
	if ss[3] == "0" {
		ss[3] = ""
	}
	return lib.Location{
		Country: ss[0],
		City:    ss[3],
	}, nil
}

func (s *XdbSearcher) SearchClean() {
	s.searcher.Close()
}

// map[IP](country\|city)
func (s *XdbSearcher) GetCitysIDByXdb(ipList []string) (map[lib.IP]string, error) {
	res := map[lib.IP]string{}
	for _, ip := range ipList {
		region, _ := s.searcher.SearchByStr(ip)
		ss := strings.Split(region, "|")
		if ss[0] == "0" {
			ss[0] = ""
		}
		if ss[3] == "0" {
			ss[3] = ""
		}
		res[ip] = ss[0] + "|" + ss[3]
	}
	return res, nil
}
