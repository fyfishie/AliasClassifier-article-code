/*
 * @Author: fyfishie
 * @Date: 2023-04-18:07
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-07-20:16
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package lib

import (
	"time"
)

type RDNSResItem struct {
	IP      string   `json:"ip"`
	Domains []string `json:"domains"`
}

type Location struct {
	Country string
	City    string
}

type IP = string
type IPInt = int
type Node struct {
	NodeID         string                         `json:"node_id"`
	AS             string                         `json:"as"`
	IPList         []IP                           `json:"ip_list"`
	Ttls           []int                          `json:"-"`
	TtlSimilarity  int                            `json:"-"`
	Rtts           []int                          `json:"-"`
	RttGap         int                            `json:"-"`
	OpenPorts      [][]int                        `json:"-"`
	PortSimilarity int                            `json:"-"`
	PortAndNum     int                            `json:"-"`
	PortOrNum      int                            `json:"-"`
	Traces         []*RawTrace                    `json:"-"`
	IPNeighborsMap map[string]map[string]struct{} `json:"-"`
	// IPIDs  map[IP]string `json:"ip_identification"`
}
type DomainSta struct {
	NodeID             string `json:"node_id"`
	IPTotal            int    `json:"ip_total"`
	DomainVisibleNum   int    `json:"domain_visible_number"`
	DomainVisibleScale int    `json:"domain_visible_scale"`
	Similarity         int    `json:"domain_similarity"`
}

type IPAID struct {
	IP                 string `json:"ip"`
	ID                 int    `json:"identification"`
	SequenceNum        int    `json:"sequence_number"`
	OriginateTimeStamp int    `json:"originate_timestamp"`
	ReceiveTimeStamp   int    `json:"receive_timestamp"`
	TransmitTimeStamp  int    `json:"transmit_timestamp"`
}
type TtlAndRtt struct {
	IP  string `json:"ip"`
	Rtt int    `json:"rtt"`
	Ttl int    `json:"ttl"`
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

type PortScanItem struct {
	IP        string  `json:"ip"`
	Timestamp string  `json:"timestamp"`
	Ports     []Ports `json:"ports"`
}
type Ports struct {
	Port   int    `json:"port"`
	Proto  string `json:"proto"`
	Status string `json:"status"`
	Reason string `json:"reason"`
	TTL    int    `json:"ttl"`
}
type IDPoint struct {
	IPID int `json:"ipid"`
	Time int `json:"time"`
}
type IDTimeArray []IDPoint

type IPIDTimes struct {
	IP          string      `json:"ip"`
	IDTimeArray IDTimeArray `json:"id_time_array"`
}

func (s *IDTimeArray) Len() int {
	return len(*s)
}
func (s *IDTimeArray) Less(i, j int) bool {
	if (*s)[i].Time == (*s)[j].Time {
		return (*s)[i].IPID < (*s)[j].IPID
	}
	return (*s)[i].Time < (*s)[j].Time
}
func (s *IDTimeArray) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

type NodeAndHisIPIDTimes struct {
	NodeID    string      `json:"node_id"`
	IPList    []string    `json:"ip_list"`
	IPIDTimes []IPIDTimes `json:"ip_id_time_arrays"`
}
type WhoisInfo struct {
	Belong       string
	Owner        string
	LastModified string
}
type WhoisRes struct {
	IP                string            `json:"ip"`
	WhoisInfoInstance WhoIsInfoInstance `json:"whois_info"`
}
type WhoIsInfoInstance struct {
	Domain         Domain  `json:"domain,omitempty"`
	Registrar      Contact `json:"registrar,omitempty"`
	Registrant     Contact `json:"registrant,omitempty"`
	Administrative Contact `json:"administrative,omitempty"`
	Technical      Contact `json:"technical,omitempty"`
	Billing        Contact `json:"billing,omitempty"`
}
type Domain struct {
	ID                   string    `json:"id,omitempty"`
	Domain               string    `json:"domain,omitempty"`
	Punycode             string    `json:"punycode,omitempty"`
	Name                 string    `json:"name,omitempty"`
	Extension            string    `json:"extension,omitempty"`
	WhoisServer          string    `json:"whois_server,omitempty"`
	Status               []string  `json:"status,omitempty"`
	NameServers          []string  `json:"name_servers,omitempty"`
	DNSSec               bool      `json:"dnssec,omitempty"`
	CreatedDate          string    `json:"created_date,omitempty"`
	CreatedDateInTime    time.Time `json:"created_date_in_time,omitempty"`
	UpdatedDate          string    `json:"updated_date,omitempty"`
	UpdatedDateInTime    time.Time `json:"updated_date_in_time,omitempty"`
	ExpirationDate       string    `json:"expiration_date,omitempty"`
	ExpirationDateInTime time.Time `json:"expiration_date_in_time,omitempty"`
}

// Contact storing domain contact info
type Contact struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Organization string `json:"organization,omitempty"`
	Street       string `json:"street,omitempty"`
	City         string `json:"city,omitempty"`
	Province     string `json:"province,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	Country      string `json:"country,omitempty"`
	Phone        string `json:"phone,omitempty"`
	PhoneExt     string `json:"phone_ext,omitempty"`
	Fax          string `json:"fax,omitempty"`
	FaxExt       string `json:"fax_ext,omitempty"`
	Email        string `json:"email,omitempty"`
	ReferralURL  string `json:"referral_url,omitempty"`
}
type NodePairs struct {
	NodeID string `json:"node_id"`
	Pairs  []Pair `json:"ip_pair"`
}

// {"ip":"209.159.171.239","rtt":16109300,"ttl":242}
type Sping struct {
	IP  string `json:"ip"`
	Rtt int    `json:"rtt"`
	Ttl int    `json:"ttl"`
}
type TraceSameFactor struct {
	IP     string  `json:"ip"`
	Factor float64 `json:"factor"`
}
type Point struct {
	X int
	Y int
}
type Points = []Point

// {"max_depth": 4, "max_feature": 5, "n_estimators": 100, "criterion": "gini", "oob_score": 0.9038739056729721}
type MArg struct {
	MaxDepth   int `json:"max_depth"`
	MaxFeature int `json:"max_feature"`
	//NEstimators int     `json:"n_estimators"`
	Criterion string `json:"criterion"`
	//OobScore    float64 `json:"oob_score"`
	Recall float64 `json:"recall"`
	Pre    float64 `json:"pre"`
	F1     float64 `json:"f1"`
}

type MArgs []MArg

func (a MArgs) Len() int      { return len(a) }
func (a MArgs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// func (a MArgs) Less(i, j int) bool { return a[i].OobScore < a[j].OobScore }
func (a MArgs) Less(i, j int) bool { return a[i].Pre < a[j].Pre }
