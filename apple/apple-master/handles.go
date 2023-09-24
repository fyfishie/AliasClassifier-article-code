/*
 * @Author: fyfishie
 * @Date: 2023-04-03:16
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-13:20
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package main

import (
	aliascek "aliasParseMaster/aliaschecker"
	"aliasParseMaster/divider"
	"aliasParseMaster/lib"
	"aliasParseMaster/utils"
	"aliasParseMaster/vpcluster"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func statusHandleFunc(w http.ResponseWriter, r *http.Request) {

}
func VPRegistHandleFunc(w http.ResponseWriter, r *http.Request) {
	logrus.Println("new req")
	slaveURL, slaveName, err := parseVPRegist(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.MakeMessageBytes(false, err.Error()))
	}
	UID, err := vpCluster.AddVP(slaveURL, slaveName)
	if err != nil {
		if err != vpcluster.ErrVPAddrAlreadyExist {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.MakeMessageBytes(false, err.Error()))
			return
		}
	}
	w.Write(utils.MakeMessageBytes(true, UID))
	fmt.Println(vpCluster.AllVP())
}

/*
* @description: we return top-task-id immediately when we get it
and process should go on later
*/
func StreamWorkHandleFunc(w http.ResponseWriter, r *http.Request) {
	userInput, err := StreamWorkReqCheck(w, r)
	//======================================
	//TODO:tmp vp set
	allvp := vpCluster.AllVP()
	vpList := []lib.VP{}
	for _, vp := range allvp {
		vpList = append(vpList, vp)
	}
	userInput.VPSelected = activeVP(vpList)
	// =====================================
	topTaskID := taskIDGen.NewID()
	if err == nil {
		w.Write(utils.MakeMessageBytes(true, strconv.Itoa(topTaskID)))
		taskStatusMap[topTaskID] = &lib.TaskStatus{
			StatusCode: lib.TASK_STATUS_WAITING,
		}
	} else {
		taskIDGen.RecycleID(topTaskID)
	}
	go streamParse(userInput.TaskName, userInput.IPToDoList, userInput.VPSelected, userInput.MaxBufferLen, topTaskID)
}

func streamParse(taskName string, ipList []lib.IP, vpSelected []lib.VP, bufferLen int, topTaskID int) {
	ipGroupChan := make(chan []lib.IntIP)
	//groupNum is accordance for task result waiting
	groupNum, err := makeDividedIP(ipList, ipGroupChan)
	if err != nil {
		logrus.Errorf("error while make divided ip, err = %v\n", err.Error())
	}

	balanceDispatchVPTask(ipGroupChan, vpSelected, topTaskID, bufferLen)
	// parser, taskID, err := streamparser.NewParser(ipList, taskName, taskIDGen, childIDGen, vpCluster, vpPublisher, collector, aliasCekPublisher)
	if err != nil {
		taskStatusMap[topTaskID].Message = "can not dispatch vp task, err = " + err.Error()
		return
	}

	//map[city(int id format)]map[VPKey]struct{}{}
	taskWaitMap := map[int]map[string]*lib.SlaveTask{}
	for i := 0; i < groupNum; i++ {
		taskWaitMap[i] = map[string]*lib.SlaveTask{}
	}
	callBackChan := collector.Wait(topTaskID, groupNum)

	for taskBack := range callBackChan {
		if _, ok := taskWaitMap[taskBack.ChildID]; ok {
			taskWaitMap[taskBack.ChildID][taskBack.TargetVPUID] = taskBack
		}
		doneCity := anyCityDone(len(vpSelected), taskWaitMap)
		if doneCity != -1 {
			aliasCheckResultChan := make(chan lib.AliasCheckResult)
			SafeAliasChecker.TaskInChan <- &aliascek.CheckRequest{
				ResultOut: aliasCheckResultChan,
				Task:      MakeAliasCheckTask(taskBack),
			}
			res := <-aliasCheckResultChan
			updateTaskStatusMap(res.TopTaskID, res.ChildID, res.Status, res.ErrMessage)
		}
	}
	return
}

/*
* @description: there is some exist check for map
incase of some error happened on topTaskID and ChildID which will cause nil pointer error

* @block: true
*/
func updateTaskStatusMap(TopTaskID, ChildID int, ChildStatus int, ChildTaskMsg string) {

}

/*
* @description: because each time this function is called
the deep number of taskWaitMap only increase by one,
so we returns only one city id each time

* @param {int} vpNum

* @param {map[int]map[string]struct{}} taskWaitMap

* @return return doneCity >= 0 if any done city found, or return doneCity = -1

* @block:
*/
func anyCityDone(vpNum int, taskWaitMap map[int]map[string]*lib.SlaveTask) (doneCity int) {
	for cityID, vpMap := range taskWaitMap {
		if len(vpMap) == vpNum {
			delete(taskWaitMap, cityID)
		}
		return cityID
	}
	return -1
}

/*
* @description: here, attention, maybe updated to use "up load a file of ip list"

* @param {[]string} ipList

* @param {chan[]lib.IP} ipGroupChan

* @return {*}

* @block: false
 */
func makeDividedIP(ipList []string, ipGroupChan chan []lib.IntIP) (groupNum int, err error) {
	//write into tmp file
	rwfi, err := os.CreateTemp(os.TempDir(), "alias_parse_divide_iplist")
	if err != nil {
		logrus.Errorf("error while create tmp file for divide ip, err = %v\n", err.Error())
		return 0, err
	}
	wtr := bufio.NewWriter(rwfi)
	for _, ip := range ipList {
		wtr.WriteString(ip + "\n")
	}
	wtr.Flush()

	//I have a simple and 100w buflen maybe tolerent for mongodb
	// which has a limit of single document in 16M
	d := divider.NewDivider(200000)
	rwfi.Name()
	return d.DivideWithMongo(rwfi.Name(), lib.MongoCollectionAddr{
		IP:             globalMongoAddr.IP,
		Port:           globalMongoAddr.Port,
		DBName:         "alias_parse_master_divide_ip",
		CollectionName: rwfi.Name()},
		ipGroupChan)
}

/*
* @description:

* @param {chan[]string} ipGroupChan

* @param {[]lib.VP} vpList

* @param {int} topTaskID

* @param {int} bufferLen

* @return {*}

* @block: false
 */
func balanceDispatchVPTask(ipGroupChan chan []lib.IntIP, vpList []lib.VP, topTaskID int, bufferLen int) {
	buffer := [][]lib.IntIP{}
	leftIndex := 0
	for group := range ipGroupChan {
		buffer = append(buffer, group)
		if len(buffer) >= bufferLen {
			bufferFlush(buffer, vpList, topTaskID, leftIndex)
			buffer = [][]lib.IntIP{}
			leftIndex += len(buffer)
		}
	}
	if len(buffer) > 0 {
		bufferFlush(buffer, vpList, topTaskID, leftIndex)
		buffer = [][]lib.IntIP{}
		leftIndex += len(buffer)
	}
}

func bufferFlush(taskIPLists [][]lib.IntIP, vpList []lib.VP, topTaskID int, baskIndex int) {
	warp := 0
	for _, vp := range vpList {
		for i := 0; i < len(taskIPLists); i++ {
			warp++
			warp %= len(taskIPLists)
			task := makeSlaveTask(topTaskID, warp, taskIPLists[warp], vp)
			bs, _ := json.Marshal(task)
			err := vpPublisher.Publish(map[string][]byte{task.TargetVPUID: bs})
			if err != nil {
				logrus.Errorf("error while publish task, err = %v\n", err.Error())
			}
		}
		warp++
	}
}

func makeSlaveTask(topTaskID, childID int, ipList []lib.IntIP, targetVP lib.VP) lib.SlaveTask {
	slaveTask := lib.SlaveTask{
		TopTaskID:  topTaskID,
		ChildID:    childID,
		TaskType:   lib.TASK_TYPE_DETECT_AND_PARSE,
		IPToDoList: ipList,
		AntiAliasResAddr: lib.MCA{
			DBName:         lib.MongoDBName_AntiAliasParseResult,
			CollectionName: strconv.Itoa(topTaskID) + "_" + strconv.Itoa(childID),
		},
		PingResAddr: lib.MCA{
			DBName:         lib.MongoDBName_PingResult,
			CollectionName: strconv.Itoa(topTaskID) + "_" + strconv.Itoa(childID),
		},
		TargetVPUID: targetVP.UID,
	}
	return slaveTask
}
func MakeAliasCheckTask(slaveTask *lib.SlaveTask) *lib.AliasCekTask {
	return &lib.AliasCekTask{
		PingResAddr: lib.MCA{
			DBName:         lib.MongoDBName_PingResult,
			CollectionName: utils.UniqueMongoCollection(slaveTask.TopTaskID, slaveTask.ChildID),
		},
		AntiAliasResAddr: lib.MCA{
			DBName:         lib.MongoDBName_AntiAliasParseResult,
			CollectionName: utils.UniqueMongoCollection(slaveTask.TopTaskID, slaveTask.ChildID),
		},
		AliasCekResAddr: lib.MCA{
			DBName:         lib.MongoDBName_AliasCheckResult,
			CollectionName: utils.UniqueMongoCollection(slaveTask.TopTaskID, slaveTask.ChildID),
		},
		MaybeAliasMCA: lib.MCA{
			DBName:         lib.MongoDBName_MaybeAliasSet,
			CollectionName: utils.UniqueMongoCollection(slaveTask.TopTaskID, slaveTask.ChildID),
		},
	}
}

func activeVP(vpList []lib.VP) (activedVP []lib.VP) {
	activedVP = []lib.VP{}
	for _, vp := range vpList {
		msgbs := utils.MakeMessageBytes(true, "active request")
		resp, err := http.Post(vp.URL+"/active", "application/json", bytes.NewReader(msgbs))
		if err != nil {
			logrus.Infof("failed to active vp, err = %v\n", err.Error())
			continue
		}
		rbs, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Infof("failed to read active response vp, err = %v\n", err.Error())
			continue
		}
		msg, err := utils.ParseMessage(rbs)
		if err != nil {
			logrus.Infof("failed to parse message from vp active response, err = %v\n", err.Error())
			continue
		}
		if msg.Status {
			activedVP = append(activedVP, vp)
		}
	}
	return activedVP
}
