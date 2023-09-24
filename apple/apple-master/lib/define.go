/*
  - @Author: fyfishie
  - @Date: 2023-02-20:09

* @LastEditors: fyfishie

* @LastEditTime: 2023-05-17:08
  - @Description: :)
  - @email: muren.zhuang@outlook.com
*/
package lib

type DetectTask struct {
	IPToDoList   []string            `json:"ip_to_do_list"`
	PingResAddr  MongoCollectionAddr `json:"ping_result_mongo_addr"`
	TraceResAddr MongoCollectionAddr `json:"trace_result_mongo_addr"`
}
type AntiAliasTask struct {
	//who detect our target ips
	VPIP string `json:"vp_ip"`

	//which ips are devided into the area the vp located
	TargetIPList []string `json:"target_ip_list"`

	//in which collection the ping result are stored
	PingResAddr MongoCollectionAddr `json:"ping_result_addr"`

	//in which collection the trace result are storeTaskStatusCoded
	TraceResAddr MongoCollectionAddr `json:"trace_result_addr"`

	//to which the antialias parse result will be stored
	AntiAliasResAddr MongoCollectionAddr `json:"antialias_result_addr"`
}
type AliasCekTask struct {
	//ip list for this alias check
	// TargetIPList []int `json:"target_ip"`
	//vp who is responsible for ip list above
	// TargetVPUID string `json:"target_vp"`
	//collection that contains ping result of target ip list
	PingResAddr MongoCollectionAddr `json:"ping_result_addr"`
	//collection that contains antialias parse result
	AntiAliasResAddr MongoCollectionAddr `json:"antialias_result_addr"`
	//to which collection the final result will be writen
	AliasCekResAddr MongoCollectionAddr `json:"AliasCekResAddr"`
	//tmp collection for T Set store
	TSetMCA MongoCollectionAddr
	//tmp collection for M Set store
	MaybeAliasMCA MongoCollectionAddr
}

// 该结构体用于唯一标识一个分析任务
// 并且保存该任务的数据
type SlaveTask struct {
	TopTaskName string `json:"top_task_name"`

	// TopTaskID全网唯一[doge],标识一个任务
	TopTaskID int `json:"top_task_id"`

	// childID is presentative of city
	ChildID int `json:"child_id"`

	IPToDoList []int `json:"ip_to_do_list"`

	TaskType string `json:"task_type"`

	//TargetVPUID is presentative of vp, for queue in rabbitmq
	TargetVPUID string `json:"target_vp_uid"`

	PingResAddr MongoCollectionAddr `json:"ping_result_mongo_addr"`

	AntiAliasResAddr MongoCollectionAddr `json:"trace_result_mongo_addr"`
}
type MongoCollectionAddr struct {
	IP             string `json:"ip"`
	Port           string `json:"port"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	DBName         string `json:"dbname"`
	CollectionName string `json:"collection_name"`
}

type MessageWithID struct {
	WorkID  int    `json:"work_id"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type Message struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
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

type UserInput struct {
	TaskName     string   `json:"task_name"`
	IPToDoList   []string `json:"ip_to_do_list"`
	VPSelected   []VP     `json:"vp_selected"`
	MaxBufferLen int      `json:"max_buffer_len"`
}

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type RegistResp struct {
	MQIP    string `json:"rabbitmq_ip"`
	MQPort  string `json:"rabbitmq_port"`
	MQVHost string `json:"rabbitmq_vhost"`
	//where to get task
	ConsumerQueueName   string `json:"consumer_queue_name"`
	PublishExchangeName string `json:"publsher_exchange_name"`
	PublishExchangeKind string `json:"publisher_exchange_kind"`
	PublishQueueName    string `json:"publisher_queue_name"`
}

type TaskStatus struct {
	StatusCode     TaskStatusCode
	Message        string
	ChildStatusMap map[int]TaskStatusCode
	ChildStatusMsg map[int]string
}

type IP = string
type IntIP = int
type MCA = MongoCollectionAddr
type VP struct {
	UID string
	URL string
}

type LocationID = string
type CityIPForMongo struct {
	LocIntID int   `json:"city"`
	IPList   []int `json:"ip_list"`
}
type TaskStatusCode = int

//	type SpingResult struct {
//		VP              VP       `json:"vp"`
//		SpingDetectItem RawSping `json:"sping_item"`
//	}
type MaybeAlias struct {
	IP            int   `json:"ip"`
	MaybeAliasSet []int `json:"maybe_alias_set"`
}
type SureAlias struct {
	IP           int   `json:"ip"`
	SureAliasSet []int `json:"maybe_alias_set"`
}
type AliasCheckResult struct {
	//you know?...
	TopTaskID int `json:"top_task_id"`
	ChildID   int `json:"child_id"`

	//representative of AntiAlias Parse Rersult
	Status     int    `json:"status"`
	ErrMessage string `json:"error_message"`
}

type SpingWithUID struct {
	VPUID string `json:"vp_uid"`
	IP    string `json:"ip"`
	Ttl   string `json:"ttl"`
}
type SpingWithoutUID struct {
	IP  int `json:"ip"`
	Ttl int `json:"ttl"`
}
type RawSping struct {
	IP  string `json:"ip"`
	Rtt int    `json:"rtt"`
	Ttl int    `json:"ttl"`
}
type VPUID = string
