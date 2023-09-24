/*
 * @Author: fyfishie
 * @Date: 2023-03-22:10
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-08:21
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package main

import (
	"fmt"
	"io"
	"net/http"
	"vp/utils"

	"github.com/sirupsen/logrus"
)

// type registResp struct {
// 	MQIP    string `json:"rabbitmq_ip"`
// 	MQPort  string `json:"rabbitmq_port"`
// 	MQVHost string `json:"rabbitmq_vhost"`
// 	//where to get task
// 	ConsumerQueueName   string `json:"consumer_queue_name"`
// 	PublishExchangeName string `json:"publsher_exchange_name"`
// 	PublishExchangeKind string `json:"publisher_exchange_kind"`
// 	PublishQueueName    string `json:"publisher_queue_name"`
// }

/*
* @description: regist runs in this way:
vp send a regist request with its own public ip address to master,
then get a rabbitmq descriptor which contains the infomation required to connect rabbitmq.
how does master know live status of vp?
vp sends packages every once in a while,
if master doesn't receive packages multiple times in a row,
it thinks this vp is dead.
But I haven't make a decision about the length sequence above.
* @param {string} masterAddr: 1.2.3.4:5678
* @return {*}
*/
func registVP2Master() {
	resp, err := http.Get(masterURL + "/api/regist/vp?slave_url=" + myURL + "&name=" + myName)
	if err != nil {
		logrus.Panic("error in regist vp to master, err = " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		logrus.Panic("bad request of regist")
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Panicf("error in read master response while regist,err = %v\n", err.Error())
	}
	msg, err := utils.ParseMessage(bs)
	if err != nil {
		logrus.Panicf("error in parse master response while regist,err = %v\n", err.Error())
	}
	if msg.Status {
		logrus.Infof("regist success!")
		fmt.Println("regist success!")
		myUID = msg.Message
		consumerQueueName = myUID
	} else {
		logrus.Panicf("regist failed, error message from master: %v\n", err.Error())
	}
}
