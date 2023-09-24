/*
* @Author: fyfishie
* @Date: 2023-03-27:08

 * @LastEditors: fyfishie

 * @LastEditTime: 2023-05-13:10

* @Description: :)
* @email: fyfishie@outlook.com
*/
package main

import (
	aliascek "aliasParseMaster/aliaschecker"
	"aliasParseMaster/idgen"
	"aliasParseMaster/lib"
	"aliasParseMaster/rabbitmq"
	"aliasParseMaster/rescollector"
	"aliasParseMaster/vpcluster"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-ini/ini"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var ()

// regist data
// var aliasCekRegistResponse = lib.RegistResp{}
// var aliasCekRegistRespBytes = []byte{}
// var VPRegistResponse = lib.RegistResp{}
// var VPRegistRespBytes = []byte{}
// var sieveRegistResponse = lib.RegistResp{}
// var sieveRegistRespBytes = []byte{}

// for local with username and password
var (
	mqIP             string
	mqPort           string
	mqUri            string
	mqVhost          string
	mqRetry          int
	mqPublishTimeout int
	
	servePort        string
	vpCluster        *vpcluster.Cluster
	taskIDGen        *idgen.Generator
	vpIDGen          *idgen.Generator
	childIDGen       *idgen.Generator
	vpPublisher      *rabbitmq.Publisher
	vpResultConsumer *rabbitmq.Consumer
	vpResultInChan   chan amqp091.Delivery
	collector        *rescollector.Collector
	globalMongoAddr        lib.MCA
	taskStatusMap    map[int]*lib.TaskStatus
	SafeAliasChecker *aliascek.ThreadSafeChecker
)

func startMaster() {
	loginit()
	loadIni()
	initGlobalData()
	fmt.Println("init done!")
	initHandles()
}
func loginit() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	fi, err := os.OpenFile("./log.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error in init logrus, err = %v\n", err.Error())
		os.Exit(1)
	}
	logrus.SetOutput(fi)
}
func initHandles() {
	http.HandleFunc("/api/regist/vp", VPRegistHandleFunc)
	// http.HandleFunc("/api/regist/alias_check", aliasCekRegistHandleFunc)
	// http.HandleFunc("/api/regist/sieve", sieveRegistHandleFunc)
	http.HandleFunc("/api/status", statusHandleFunc)
	http.HandleFunc("/api/stream_work", StreamWorkHandleFunc)
	http.ListenAndServe(":8080", nil)
}

func initGlobalData() {
	var err error
	vpCluster = vpcluster.NewCluster().Run()
	taskIDGen = idgen.NewIdGen().Run()
	childIDGen = idgen.NewIdGen().Run()

	vpPublisher = rabbitmq.NewPublisher(
		mqUri,
		mqRetry,
		rabbitmq.ExchangeInfo{Name: lib.MQ_VP_TASK_EXCHANGE_NAME, Kind: amqp091.ExchangeDirect},
		mqPublishTimeout,
		context.TODO())
	vpPublisher.Run()

	vpResultConsumer = rabbitmq.NewConsumer(mqUri,
		mqRetry,
		lib.MQ_VP_RESULT_QUEUE_BINDKEY,
		rabbitmq.QueueInfo{Name: lib.MQ_VP_RESULT_QUEUE_NAME},
		rabbitmq.ExchangeInfo{Name: lib.MQ_VP_RESULT_EXCHANGE_NAME, Kind: amqp091.ExchangeDirect})
	vpResultInChan, err = vpResultConsumer.RunAsPushMod(2)
	if err != nil {
		logrus.Panicf("error while start vp result consumer, err = %v\n", err.Error())
	}

	collector = rescollector.NewCollector(vpResultInChan).Run()
	SafeAliasChecker = aliascek.NewThreadSafeChecker()
	taskStatusMap = make(map[int]*lib.TaskStatus)
}
func loadIni() {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		logrus.Panicf("error in load config.ini, err = %v\n", err.Error())
		os.Exit(1)
	}
	section := cfg.Section("mongodb")
	globalMongoAddr = lib.MCA{
		IP:       section.Key("ip").String(),
		Port:     section.Key("port").String(),
		Username: section.Key("username").String(),
		Password: section.Key("password").String(),
	}

	section = cfg.Section("rabbitmq")
	//amqp://user:pass@host:10000/vhost
	mqUri = fmt.Sprintf("amqp://%v:%v@%v:%v/%v",
		section.Key("username").String(),
		section.Key("password").String(),
		section.Key("ip").String(),
		section.Key("port").String(),
		section.Key("vhost").String())
	mqRetry = section.Key("retry").MustInt(10)
	mqPublishTimeout = section.Key("publish_timeout").MustInt(10)
	mqIP = section.Key("ip").String()
	mqPort = section.Key("port").String()
	mqVhost = section.Key("vhost").String()
	section = cfg.Section("serve")
	servePort = section.Key("serve_port").String()
}
