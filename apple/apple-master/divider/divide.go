/*
* @Author: fyfishie
* @Date: 2023-05-09:20

  - @LastEditors: fyfishie

  - @LastEditTime: 2023-05-13:15

* @@email: fyfishie@outlook.com
* @Description: :)
*/
package divider

import (
	"aliasParseMaster/geocode"
	"aliasParseMaster/lib"
	"aliasParseMaster/mongoproxy"
	"bufio"
	"os"

	"github.com/fyfishie/ipop"
	"github.com/sirupsen/logrus"
)

/*
* @description:
in case of more and more and more ip input, we put them in a tmp file instead of passing an array of ip
* @param {string} rdpath

* @param {lib.MCA} tmpMCA

* @param {chan[]lib.IP} ipGroupChan

* @return {*}

* @block: false
*/
func (d *Divider) DivideWithMongo(rdpath string, tmpMCA lib.MCA, ipGroupChan chan []lib.IntIP) (groupNum int, err error) {
	s := geocode.NewXdbSearcher()
	err = s.SearchInit("./divider/ip2location.xdb")
	if err != nil {
		logrus.Errorf("error while init searcher, err = %v\n", err.Error())
		return 0, err
	}
	buffer := []string{}
	rfi, err := os.OpenFile(rdpath, os.O_RDONLY, 0000)
	if err != nil {
		return 0, err
	}
	defer rfi.Close()

	inChan := make(chan map[string]lib.LocationID)
	writeDone, err := d.prepareMongoDB(inChan, tmpMCA)

	scaner := bufio.NewScanner(rfi)
	if err != nil {
		logrus.Errorf("error while prepare database, err = %v\n", err.Error())
		return 0, err
	}
	for {
		for i := 0; i < d.bufferLen && scaner.Scan(); i++ {
			buffer = append(buffer, scaner.Text())
		}
		if len(buffer) == 0 {
			break
		}
		locations, err := s.GetCitysIDByXdb(buffer)
		if err != nil {
			logrus.Errorf("error while request city info, err = %v\n", err.Error())
			continue
		}
		inChan <- locations
		buffer = []string{}
	}
	close(inChan)
	<-writeDone
	groupNum, err = d.readFromMongo(tmpMCA, ipGroupChan)
	if err != nil {
		logrus.Errorf("error while start ip group read process, err = %v\n", err.Error())
		return 0, err
	}
	return groupNum, nil
}

// type CityIPForMongo struct {
// 	LocIntID int   `json:"city"`
// 	IPList   []int `json:"ip_list"`
// }

func (d *Divider) prepareMongoDB(inChan chan map[string]lib.LocationID, tmpMCA lib.MCA) (writeDone chan struct{}, err error) {
	p := mongoproxy.NewProxy(tmpMCA)

	err = p.Connect()
	if err != nil {
		logrus.Errorf("error while connect to tmp mongo, err = %v\n", err.Error())
		return nil, err
	}
	doneChan := make(chan struct{})
	go func() {
		for locations := range inChan {
			docs := d.locs2CityIPForMongo(locations)
			err = p.InsertManyInterfaces(docs)
			if err != nil {
				logrus.Errorf("error while insert docs, err = %v\n", err.Error())
			}
		}
		doneChan <- struct{}{}
		p.Disconnect()
	}()
	return doneChan, nil
}

func (d *Divider) locs2CityIPForMongo(locations map[string]lib.LocationID) []interface{} {
	datas := []lib.CityIPForMongo{}
	cityIPMap := map[string][]lib.IP{}
	for ip, locID := range locations {
		if _, ok := cityIPMap[locID]; !ok {
			cityIPMap[locID] = []lib.IP{}
		}
		cityIPMap[locID] = append(cityIPMap[locID], ip)
	}
	for locID, ipList := range cityIPMap {
		if _, ok := d.locationDisMap[locID]; !ok {
			d.locationDisMap[locID] = d.locationIDGen.NewID()
		}
		locIntID := d.locationDisMap[locID]
		ipl := []int{}
		for _, ip := range ipList {
			ipl = append(ipl, ipop.String2Int(ip))
		}
		datas = append(datas, lib.CityIPForMongo{
			LocIntID: locIntID,
			IPList:   ipl,
		})
	}
	docs := []interface{}{}
	for _, r := range datas {
		docs = append(docs, r)
	}
	return docs
}

/*
@description:

@param {lib.MCA} tmpMCA

@param {chan[]lib.IP} ipGroupChan

@return {*}

@block: false
*/
func (d *Divider) readFromMongo(tmpMCA lib.MCA, ipGroupChan chan []lib.IntIP) (groupNum int, err error) {
	//query all country\|city
	p := mongoproxy.NewProxy(tmpMCA)
	err = p.Connect()
	if err != nil {
		logrus.Errorf("error while connect to mongo, err = %v\n", err.Error())
		return 0, err
	}
	groupNum, err = p.QueryIPGroupByLocID(ipGroupChan)
	if err != nil {
		logrus.Errorf("error while query ip group be locatio id, err = %v\n", err.Error())
		return 0, err
	}
	return groupNum, nil
}
