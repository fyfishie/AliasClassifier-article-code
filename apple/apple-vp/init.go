/*
 * @Author: fyfishie
 * @Date: 2023-03-03:16
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-13:10
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"vp/rabbitmq"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	myName        string
	myUID         string
	myIP          string
	myPort        string
	taskIn        chan amqp091.Delivery
	consumer      *rabbitmq.Consumer
	actived       bool
	lock          sync.Mutex
	mongoUsername string
	mongoPassword string

	// var resultOutChan chan farm.PublishTask
	slaveResultPublisher *rabbitmq.Publisher
	mqURI                string
	rabbitIP             string
	rabbitPort           string
	rabbitRetry          int
	rabbitPublishTimeout int
	rabbitUsername       string
	rabbitPassword       string
	rabbitVHost          string
	taskExchangeName     string
	consumerQueueName    string
	resultExchangeName   string
	masterURL            string
	myURL                string
)

func init() {
	lock = sync.Mutex{}
	logInit()
	configInit()
}
func logInit() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	fi, err := os.OpenFile("./log.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error in init logrus, err = %v\n", err.Error())
		os.Exit(1)
	}
	logrus.SetOutput(fi)
}
func configInit() {
	//slave
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	section := cfg.Section("slave")
	myName = section.Key("my_name").String()

	//rabbitmq
	section = cfg.Section("rabbitmq")
	rabbitPublishTimeout = section.Key("publish_timeout").MustInt(10)
	if err != nil {
		rabbitPublishTimeout = 10
	}
	rabbitIP = section.Key("ip").String()
	rabbitPort = section.Key("port").String()
	rabbitVHost = section.Key("vhost").String()
	rabbitRetry = section.Key("retry").MustInt(10)
	rabbitUsername = section.Key("username").String()
	rabbitPassword = section.Key("password").String()
	taskExchangeName = section.Key("task_exchange_name").String()
	resultExchangeName = section.Key("result_exchange_name").String()
	mqURI = "amqp://" + rabbitUsername + ":" + rabbitPassword + "@" + rabbitIP + ":" + rabbitPort + "/" + rabbitVHost
	//master
	section = cfg.Section("master")
	masterURL = section.Key("master_url").String()

	//local
	section = cfg.Section("local")
	myIP = section.Key("my_ip").String()
	myPort = section.Key("my_port").String()
	myURL = "http://" + myIP + ":" + myPort

	section = cfg.Section("mongo")
	mongoUsername = section.Key("username").String()
	mongoPassword = section.Key("password").String()
}

// func defaultLocation() (lat, lng string) {
// 	cfg, err := ini.Load("config.ini")
// 	if err != nil {
// 		logrus.Fatalf("can not load config.ini, err = %v\n", err.Error())
// 	}
// 	section := cfg.Section("location")
// 	lati, err := section.Key("latitude").Int()
// 	if err != nil {
// 		logrus.Fatalf("latitude of location in config.ini is invalid, check please\n")
// 	}
// 	lngi, err := section.Key("longitude").Int()
// 	if err != nil {
// 		logrus.Fatalf("latitude of lonitude in config.ini is invalid, check please\n")
// 	}
// 	lat = strconv.Itoa(lati)
// 	lng = strconv.Itoa(lngi)
// 	return lat, lng
// }
