/*
* @Author: fyfishie
 * @Date: 2023-03-27:10
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-09:19
 * @Description: :)
 * @email: fyfishie@outlook.com
*/
package geocode

import (
	"aliasParseMaster/lib"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type CtyCity struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

var Err_Can_Not_Get_City = errors.New("can not get city")

// returns location of each ip in ipList,
// but only location of ips those who have iplocation info in db will be added to result
func GetCitys(ipList []string) (map[lib.IP]*lib.Location, error) {
	res := map[lib.IP]*lib.Location{}
	bs, err := json.Marshal(ipList)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(ip2CityUrl(), "application/json", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	rbs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rbs, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TODO:correct ./config.int path after debug
func ip2CityUrl() string {
	cfg, err := ini.Load("../config.ini")
	if err != nil {
		logrus.Panicf("error in load config.ini, err = %v\n", err.Error())
		os.Exit(1)
	}
	section := cfg.Section("iplocate")
	url := section.Key("ip2city_url").String()
	return url
}

// func city2GPSUrl() string {
// 	cfg, err := ini.Load("./config.ini")
// 	if err != nil {
// 		logrus.Panicf("error in load config.ini, err = %v\n", err.Error())
// 		os.Exit(1)
// 	}
// 	section := cfg.Section("iplocate")
// 	url := section.Key("city2gps_url").String()
// 	return url
// }

// type GPSReqAndResp struct {
// 	Exist     bool    `json:"exist"`
// 	Country   string  `json:"country"`
// 	City      string  `json:"city"`
// 	Latitude  float64 `json:"lat"`
// 	Longitude float64 `json:"lng"`
// }

// locations: map[country+city]lib.Location
// fills locations with latitude and longitude
// func GetGPS(locations map[string]*lib.Location) error {
// 	url := city2GPSUrl()
// 	reqData := []GPSReqAndResp{}
// 	for _, location := range locations {
// 		reqData = append(reqData, GPSReqAndResp{
// 			Country: location.Country,
// 			City:    location.City,
// 		})
// 	}
// 	bs, err := json.Marshal(reqData)
// 	if err != nil {
// 		return err
// 	}
// 	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
// 	if err != nil {
// 		return err
// 	}
// 	rbs, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	respData := []GPSReqAndResp{}
// 	err = json.Unmarshal(rbs, &respData)
// 	if err != nil {
// 		return err
// 	}
// 	for _, r := range respData {
// 		if _, ok := locations[r.Country+r.City]; ok {
// 			if r.Exist {
// 				locations[r.Country+r.City].GPS.Latitude = r.Latitude
// 				locations[r.Country+r.City].GPS.Longitude = r.Longitude
// 			}
// 		}
// 	}
// 	return nil
// }

// func FillGPS(ipLocMap map[lib.IP]*lib.Location) error {
// 	locations := map[string]*lib.Location{}
// 	for _, locPtr := range ipLocMap {
// 		loc := *locPtr
// 		locations[locPtr.Country+locPtr.City] = &loc
// 	}
// 	err := GetGPS(locations)
// 	if err != nil {
// 		return err
// 	}
// 	for _, locPtr := range ipLocMap {
// 		if loc, ok := locations[locPtr.Country+locPtr.City]; ok {
// 			locPtr.GPS.Latitude = loc.GPS.Latitude
// 			locPtr.GPS.Longitude = loc.GPS.Longitude
// 		}
// 	}
// 	return nil
// }
