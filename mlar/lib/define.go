/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-08-03:17
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package lib

import (
	"strconv"

	"github.com/fyfishie/ipop"
)

type ISP struct {
	IP  string `json:"ip"`
	Isp string `json:"isp"`
}
type Pair struct {
	IPA   string `json:"ipa"`
	IPB   string `json:"ipb"`
	ANode string `json:"a_node"`
	BNode string `json:"b_node"`
}

func (i *Pair) ID() string {
	ipa := ipop.String2Int(i.IPA)
	ipb := ipop.String2Int(i.IPB)
	if ipa < ipb {
		return strconv.Itoa(ipa) + strconv.Itoa(ipb)
	} else {
		return strconv.Itoa(ipb) + strconv.Itoa(ipa)
	}
}

type RawTrace struct {
	Ip      string `json:"ip"`
	Results []Hop  `json:"results"`
	Circle  bool   `json:"circle"`
}
type Hop struct {
	TTL    int    `json:"ttl"`
	Status int    `json:"status"`
	Ip     string `json:"ip"`
	Rtt    int    `json:"rtt"`
}

/*
	{
		"data": [
	    {
	      "accuracy": "CN",
	      "areacode": "CN",
	      "asnumber": "4538",
	      "continent": "亚洲",
	      "country": "中国",
	      "ip": "222.194.15.1",
	      "isp": "中国教育网",
	      "multiAreas": [{
	        "city": "威海市",
	        "district": "环翠区",
	        "latbd": "37.540047", // 百度坐标系
	        "latwgs": "37.533252", //WGS坐标系
	        "lngbd": "122.089909",
	        "lngwgs": "122.078279",
	        "prov": "山东省",
	        "radius": "0.1813" // 定位精度半径
	      }],
	      "timezone": "UTC+8",
	      "zipcode": "264200",
	      "owner":"中国教育网",
	    },
	  ],
		"status": true
	}
*/
type ISPRes struct {
	Data   []Data `json:"data"`
	Status bool   `json:"status"`
}
type Data struct {
	Accuracy   string      `json:"accuracy"`
	Areacode   string      `json:"areacode"`
	Continent  string      `json:"continent"`
	Country    string      `json:"country"`
	IP         string      `json:"ip"`
	ISP        string      `json:"isp"`
	MultiAreas []MultiArea `json:"multiareas"`
	TimeZone   string      `json:"timezone"`
	ZipCode    string      `json:"zipcode"`
	Owner      string      `json:"woner"`
}
type MultiArea struct {
	City     string `json:"city"`
	District string `json:"district"`
	Latbd    string `json:"latbd"`
	Latwgs   string `json:"latwgs"`
	Lngbd    string `json:"lngbd"`
	Lngwgs   string `json:"lngwgs"`
	Prov     string `json:"prov"`
	Radius   string `json:"radius"`
}
type Sping struct {
	IP  string `json:"ip"`
	Rtt int    `json:"rtt"`
	Ttl int    `json:"ttl"`
}
