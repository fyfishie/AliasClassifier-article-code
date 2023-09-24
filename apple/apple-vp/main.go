/*
 * @Author: fyfishie
 * @Date: 2023-03-01:16
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-15:15
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package main

import (
	"encoding/json"
	"net/http"
	"vp/lib"
	"vp/mongoProxy"
	"vp/smarkDetect"
	"vp/spingdetct"
	"vp/utils"

	"github.com/sirupsen/logrus"
)

func main_serve() {
	registVP2Master()
	// go ImAlive()
	http.HandleFunc("/active", listenToBeActive)
	http.ListenAndServe(":"+myPort, nil)

	//block main goruntine
	// <-make(chan struct{})
}

// ping detect function
func runSping(task lib.SlaveTask) error {
	d := spingdetct.NewDetector(task.IPToDoList, "./bin")
	err := d.Run()
	if err != nil {
		logrus.Error("error in running sping detector, err = " + err.Error())
		return err
	}
	bfwtr, err := mongoProxy.NewBufferWriterFromAddr(task.PingResAddr, 5000)
	if err != nil {
		logrus.Errorf("error while make new buffer writer for mongo, err = %v\n", err.Error())
		return err
	}
	bfwtr.MvSpingRes2Mongo("./bin/sping_output", task.TargetVPUID)
	return nil
}

// traceroute detect function
func traceAndParse(task lib.SlaveTask) error {
	d := smarkDetect.NewSmarkDetector("./bin")
	d.SetTargerAddr(task.IPToDoList)
	err := d.DetectorRun()
	if err != nil {
		logrus.Error("error in smark detect, err = " + err.Error())
		return err
	}
	err = parseOneTraceAndMvResult2Mongo(task.AntiAliasResAddr, "")
	if err != nil {
		logrus.Errorf("error while parsing one trace data, err = %v\n", err.Error())
	}
	err = parseEnd2AndMv2Mongo(task.AntiAliasResAddr, "")
	if err != nil {
		logrus.Errorf("error while parsing end2 data, err = %v\n", err.Error())
	}
	err = parseSubnet(task.AntiAliasResAddr, utils.StringIPs2Ints(task.IPToDoList))
	if err != nil {
		logrus.Errorf("error while parsing subnet data, err = %v\n", err.Error())
	}
	return nil
}

func serve(task lib.SlaveTask) []byte {
	if len(task.IPToDoList) == 0 {
		return utils.MakeMessageBytesWithID(false, task.TopTaskID, "zero length of ip list")
	}
	err := traceAndParse(task)
	if err != nil {
		return utils.MakeMessageBytesWithID(false, task.TopTaskID, err.Error())
	}
	err = runSping(task)
	if err != nil {
		return utils.MakeMessageBytesWithID(false, task.TopTaskID, err.Error())
	}
	return utils.MakeMessageBytesWithID(true, task.TopTaskID, "ok")
}

// func ImAlive() {
// 	for {
// 		time.Sleep(time.Hour)
// 		_, err := http.Get("http://" + masterURL + "/api/vp_status?status=Im_ok")
// 		if err != nil {
// 			logrus.Error("error in send alive package to master, err = " + err.Error())
// 			return
// 		}
// 	}
// }

func processDetectTask() {
	for {
		delivery, ok := <-taskIn
		if !ok {
			logrus.Error("failed to get delivery from rabbitmq")
			return
		}
		bs := delivery.Body
		vpTask := lib.SlaveTask{}
		err := json.Unmarshal(bs, &vpTask)
		if err != nil {
			logrus.Error("error in unmarshar bytes into Task, err = " + err.Error())
			return
		}
		insertMongoAuth(&vpTask)
		switch vpTask.TaskType {
		case lib.TASK_TYPE_DETECT_AND_PARSE:
			bs = serve(vpTask)
			err = slaveResultPublisher.Publish(map[string][]byte{"vp_detect_result": bs})
			if err != nil {
				delivery.Ack(false)
			} else {
				logrus.Errorf("err while publish result, err = %v\n", err.Error())
			}
		case lib.TASK_TYPE_SLEEP:
			logrus.Infof("fall in sleep...")
			fallSleep()
			return
		default:
			bs = utils.MakeMessageBytesWithID(false, vpTask.TopTaskID, "unsupported task type")
			err = slaveResultPublisher.Publish(map[string][]byte{"vp_detect_result": bs})
			if err != nil {
				delivery.Ack(false)
			} else {
				logrus.Errorf("err while publish result, err = %v\n", err.Error())
			}
		}
	}
}

func insertMongoAuth(task *lib.SlaveTask) {
	task.AntiAliasResAddr.Username = mongoUsername
	task.AntiAliasResAddr.Password = mongoPassword
	task.PingResAddr.Username = mongoUsername
	task.PingResAddr.Password = mongoPassword
}
