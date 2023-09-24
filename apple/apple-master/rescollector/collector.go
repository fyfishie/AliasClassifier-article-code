/*
 * @Author: fyfishie
 * @Date: 2023-04-03:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-12:09
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package rescollector

import (
	"aliasParseMaster/lib"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type waiting struct {
	//childID should be in range [0,ChildID]
	MaxChildID int
	//map[childID]nothing
	ChildBackSet  map[int]struct{}
	ResultOutChan chan *lib.SlaveTask
}
type waitReq struct {
	TopTaskID  int
	MaxChildID int
	ResOutChan chan *lib.SlaveTask
}
type Collector struct {
	//map[topTaskID]waiting
	waitingMap         map[int]*waiting
	waitReqChan        chan waitReq
	targetDeliveryChan chan amqp091.Delivery
}

func NewCollector(targetDeliveryChan chan amqp091.Delivery) *Collector {
	return &Collector{
		waitingMap:         map[int]*waiting{},
		waitReqChan:        make(chan waitReq),
		targetDeliveryChan: targetDeliveryChan,
	}
}

func (c *Collector) Run() *Collector {
	go func() {
		for {
			select {
			case req := <-c.waitReqChan:
				if _, ok := c.waitingMap[req.TopTaskID]; ok {
					logrus.Errorf("waiting of topTask=%v has not done, we believe that something sends two request of the same topTask\n")
				}
				c.waitingMap[req.TopTaskID] = &waiting{MaxChildID: req.MaxChildID, ChildBackSet: map[int]struct{}{}}
			case delivery := <-c.targetDeliveryChan:
				c.handleAndNotifyTaskBack(delivery)
				c.cleanOneDoneTopTask()
			}
		}
	}()
	return c
}
func (c *Collector) cleanOneDoneTopTask() {
	for topTaskID, waiting := range c.waitingMap {
		if len(waiting.ChildBackSet) == waiting.MaxChildID+1 {
			close(waiting.ResultOutChan)
			delete(c.waitingMap, topTaskID)

			//delete only one topTask
			return
		}
	}
}

// clean all done topTask, not nessary now
func (c *Collector) cleanDoneTopTask() {
	for topTaskID, waiting := range c.waitingMap {
		if len(waiting.ChildBackSet) == waiting.MaxChildID+1 {
			close(waiting.ResultOutChan)
			delete(c.waitingMap, topTaskID)
		}
	}
}
func (c *Collector) handleAndNotifyTaskBack(delivery amqp091.Delivery) {
	slaveTask := lib.SlaveTask{}
	err := json.Unmarshal(delivery.Body, &slaveTask)
	if err != nil {
		logrus.Errorf("error while unmarshal delivery data into slaveTask, err = %v\n", err.Error())
		return
	}
	if childWaiting, ok := c.waitingMap[slaveTask.TopTaskID]; ok {
		if childWaiting.MaxChildID < slaveTask.ChildID {
			logrus.Error(
				("found a slaveTask back whose child id is out of range record, there maybe some error\n"),
				fmt.Sprintf("topTaskName:%v, topTaskID:%v, childTaskID:%v\n", slaveTask.TopTaskName, slaveTask.TopTaskID, slaveTask.ChildID))
			return
		}
		childWaiting.ChildBackSet[slaveTask.ChildID] = struct{}{}
		childWaiting.ResultOutChan <- &slaveTask
	}
}

/*
  - @description:
    this function processes thread safty problem automically
    if all task in idList is collected this function nil
    other cases... i don't kwow

  - @param {*} topTaskID

- @param {int} maxChildID

- @return {*}

- @block: false
*/
func (c *Collector) Wait(topTaskID, maxChildID int) (callBackChan chan *lib.SlaveTask) {
	resultOutChan := make(chan *lib.SlaveTask)
	c.waitReqChan <- waitReq{
		TopTaskID:  topTaskID,
		MaxChildID: maxChildID,
		ResOutChan: resultOutChan,
	}
	return resultOutChan
}
