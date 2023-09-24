/*
 * @Author: fyfishie
 * @Date: 2023-05-14:16
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-17:10
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
/*
 * @Author: fyfishie
 * @Date: 2023-05-08:15
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-14:16
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package main

import (
	"bufio"
	"encoding/json"
	"os"
	"vp/end2"
	"vp/lib"
	"vp/mongoProxy"
	"vp/onetrace"
	"vp/subnet"
	"vp/utils"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func parseOneTraceAndMvResult2Mongo(mongoAddr lib.MCA, tracePath string) error {
	traceInChan := make(chan lib.TraceRoute)
	parser := onetrace.NewParser(traceInChan).WithRun()
	go validTraceInput(traceInChan, tracePath)
	err := indb(mongoAddr, parser.ResultOutChan)
	if err != nil {
		logrus.Errorf("error while mv onetrace parse result into mongodb, err = %v\n", err.Error())
	}
	return err
}

func parseEnd2AndMv2Mongo(antiMongoAddr lib.MCA, tracePath string) error {
	traceInChan := make(chan lib.TraceRoute, 1000)
	parser := end2.NewParser(traceInChan).WithRun()
	go validTraceInput(traceInChan, tracePath)
	err := indb(antiMongoAddr, parser.ResultOutChan)
	if err != nil {
		logrus.Errorf("error while mv end2 parse result into mongo, err = %v\n", err.Error())
	}
	return err
}

func parseSubnet(mongoAddr lib.MCA, ipToDoList []int) error {
	parser := subnet.NewParser(ipToDoList, 1).WithRun()
	err := indb(mongoAddr, parser.ResultOutChan)
	if err != nil {
		logrus.Errorf("error while mv subnet parse result into mongo, err = %v\n", err.Error())
	}
	return err
}

func indb(mongoAddr lib.MCA, dataInChan chan lib.AntiDescriptor) error {
	bufwtr, err := mongoProxy.NewBufferWriterFromAddr(mongoAddr, 5000)
	if err != nil {
		return err
	}
	for descriptor := range dataInChan {
		err = bufwtr.WriteOne(mongo.NewInsertOneModel().SetDocument(inoutconv(descriptor)))
		if err != nil {
			logrus.Errorf("error while write antialias result into mongodb, err = %v\n", err.Error())
		}
	}
	bufwtr.Flush()
	return bufwtr.Close()
}
func inoutconv(i lib.AntiDescriptor) lib.AntiDescriptorForDB {
	return lib.AntiDescriptorForDB{
		Ip:           int64(i.Ip),
		AntiAliasSet: inoutSetConv(i.AntiAliasSet),
	}
}
func inoutSetConv(is []int) []int64 {
	res := []int64{}
	for _, i := range is {
		res = append(res, int64(i))
	}
	return res
}
func validTraceInput(traceInChan chan lib.TraceRoute, tracePath string) error {
	rfi, err := os.OpenFile(tracePath, os.O_RDONLY, 0000)
	if err != nil {
		return err
	}
	defer rfi.Close()
	rdr := bufio.NewReader(rfi)
	for {
		line, _, err := rdr.ReadLine()
		if err != nil {
			break
		}
		item := lib.RawTrace{}
		if json.Unmarshal(line, &item) != nil {
			continue
		}
		if len(item.Results) < 2 || item.Ip != item.Results[len(item.Results)-1].Ip {
			continue
		}
		traceInChan <- utils.RawTrace2Trace(item)
	}
	close(traceInChan)
	return nil
}
