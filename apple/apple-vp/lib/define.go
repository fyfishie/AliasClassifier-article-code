/*
 * @Author: fyfishie
 * @Date: 2023-03-01:16
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-17:10
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package lib

type ID = int
type TaskType = string
type SlaveTask struct {
	TopTaskName string   `json:"top_task_name"`
	TopTaskID   int      `json:"top_task_id"`
	ChildID     int      `json:"child_id"`
	IPToDoList  []string `json:"ip_list"`
	TaskType    TaskType `json:"task_type"`

	TargetVPUID      string              `json:"target_vp_uid"`
	AntiAliasResAddr MongoCollectionAddr `json:"trace_result_mongo_addr"`
	PingResAddr      MongoCollectionAddr `json:"ping_result_mongo_addr"`
	// AntialiasResMongoAddr MongoCollectionAddr `json:"antialias_res_mongoaddr"`
}

// smark output but not parsed
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
type TraceRoute struct {
	Target   int
	End      int
	TTLStart int
	End2     int
	Trace    []int
	Times    []float64
	Status   []int
	TTLS     []int
}
type MongoCollectionAddr struct {
	IP             string `json:"ip"`
	Port           string `json:"port"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	DBName         string `json:"dbname"`
	CollectionName string `json:"collection_name"`
}

type Message struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
type MessageWithID struct {
	ID      ID     `json:"id"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
type IpPair struct {
	Left  int
	Right int
	Next  *IpPair
}
type AntiAliasSet []int

// sping detect output but not parsed
type RawSping struct {
	IP  string `json:"ip"`
	Rtt int    `json:"rtt"`
	Ttl int    `json:"ttl"`
}
type AntiDescriptor struct {
	Ip           int   `json:"ip"`
	AntiAliasSet []int `json:"antiAliasSet"`
}
type AntiDescriptorForDB struct {
	Ip           int64   `json:"ip"`
	AntiAliasSet []int64 `json:"antiAliasSet"`
}
type MCA = MongoCollectionAddr
type End2AndEnds struct {
	End2 int   `json:"end2"`
	Ends []int `json:"ends"`
}
type SpingWithUID struct {
	VPUID string `json:"vp_uid"`
	IP    string `json:"ip"`
	Ttl   string `json:"ttl"`
}
