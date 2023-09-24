/*
 * @Author: fyfishie
 * @Date: 2023-05-08:07
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-12:10
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package main

import (
	"context"
	"io"
	"net/http"
	"vp/rabbitmq"
	"vp/utils"

	"github.com/sirupsen/logrus"
)

func listenToBeActive(w http.ResponseWriter, r *http.Request) {
	//incase of two active request in a narrow duration arise two active process runtine
	lock.Lock()
	if actived {
		w.Write(utils.MakeMessageBytes(true, "already actived"))
		return
	}
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("err while read active request body, err = %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MakeMessageBytes(false, "failed to read request body"))
		return
	}
	msg, err := utils.ParseMessage(bs)
	if err != nil {
		logrus.Errorf("err while parse message from master, err = %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.MakeMessageBytes(false, "unable to parse message"))
		return
	}
	if msg.Status {
		logrus.Infof("active request received, weeking up...")
	}
	err = weekup()
	if err != nil {
		logrus.Errorf("error while week up, err = %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MakeMessageBytes(false, "error while active machine..."))
		return
	}
	go processDetectTask()
	w.Write(utils.MakeMessageBytes(true, "active success!"))
	lock.Unlock()
}
func weekup() error {
	var err error
	publisher := rabbitmq.NewPublisher(mqURI,
		rabbitRetry,
		rabbitmq.ExchangeInfo{Name: resultExchangeName},
		rabbitPublishTimeout, context.TODO())
	err = publisher.Run()
	if err != nil {
		logrus.Errorf("failed to start publisher, err = %v\n", err.Error())
		return err
	}
	consumer = rabbitmq.NewConsumer(mqURI, rabbitRetry, myUID, rabbitmq.QueueInfo{Name: consumerQueueName}, rabbitmq.ExchangeInfo{Name: taskExchangeName})
	taskIn, err = consumer.RunAsPushMod(2)
	if err != nil {
		logrus.Errorf("failed to start consumer, err = %v\n", err.Error())
		publisher.Clean()
		return err
	}
	actived = true
	return nil
}

func fallSleep() {
	slaveResultPublisher.Clean()
	consumer.Clean()
	actived = false
}
