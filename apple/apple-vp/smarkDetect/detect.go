/*
 * @Author: fyfishie
 * @Date: 2023-03-06:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-08:10
 * @Description: traceroute detect Implemented bt exec smark
 * @email: muren.zhuang@outlook.com
 */
package smarkDetect

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"vp/status"

	"github.com/sirupsen/logrus"
)

type SmarkDetector struct {
	IPAddrList []string
	smarkHome  string
	statusPort status.StatusEye
}
type RawTraceResult struct {
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

func NewSmarkDetector(smarkHome string) *SmarkDetector {
	return &SmarkDetector{
		statusPort: status.StatusEye{},
		smarkHome:  smarkHome,
	}
}
func (s *SmarkDetector) SetTargerAddr(IPAddrList []string) {
	s.IPAddrList = IPAddrList
}

/*
* @description: run smark with task.IPBytes being input of smark
and task.smarkPath being path of smark
* @return {*}
*/
func (s *SmarkDetector) DetectorRun() error {
	_, err := os.Create(s.smarkHome + "/smark_output")
	if err != nil {
		logrus.Error("error create smark out file, err = " + err.Error())
		return err
	}
	s.statusPort.Status = "writing ip addr into smark input file..."
	err = os.Truncate(s.smarkHome+"/smark_input", 0)
	if err != nil {
		logrus.Error("error in truncate smark input file, err = " + err.Error())
	}
	err = s.writeInput()
	if err != nil {
		// logrus.Error("error in write ip addr into smark input file, err = %v\n", err.Error())
		return err
	}
	s.statusPort.Status = "traceRoute running..."
	cmd := exec.Command(s.smarkHome+"/smark", "-if", s.smarkHome+"/smark_input", "-of", s.smarkHome+"/smark_output")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error(fmt.Sprintf("error in exec smark, err = %v\noutput = %v", err.Error(), string(output)))
		return err
	}
	// s.statusPort.Status = "reading traceRoute result..."
	// res, err := s.readDetectResult()
	// if err != nil {
	// 	logrus.Error("error in read smark output, err = " + err.Error())
	// 	return nil, err
	// }
	s.statusPort.Status = "traceRoute done!"
	// return utils.RawTraces2Traces(res), err
	return nil
}

/*
* @description: write ip addr into smark input file, if any error occurred,
this function stop and returns the error untouched

the input file will be cleared at the very beginning of this function
* @return {*}
*/
func (d *SmarkDetector) writeInput() error {
	fi, err := os.OpenFile(d.smarkHome+"/smark_input", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		logrus.Error("error in open smark input file, err = " + err.Error())
		return err
	}
	wt := bufio.NewWriter(fi)
	for _, addr := range d.IPAddrList {
		_, err := wt.WriteString(addr + "\n")
		if err != nil {
			logrus.Error("error in writing ip addr into smark input, err = " + err.Error())
			return err
		}
	}
	err = wt.Flush()
	if err != nil {
		logrus.Error("error in flush data into smark input file, err = " + err.Error())
	}
	return err
}

/*
* @description: read smark output only, if any error encounted, function stops and return the error untouched
 */
// func (d *SmarkDetector) readDetectResult() (nextBlock []lib.RawTrace, err error) {
// 	fi, err := os.OpenFile(d.smarkHome+"/smark_output", os.O_RDONLY, 0444)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rdr := bufio.NewReader(fi)
// 	for {
// 		line, _, err := rdr.ReadLine()
// 		if err != nil {
// 			break
// 		}
// 		rawTrace := lib.RawTrace{}
// 		err = json.Unmarshal(line, &rawTrace)
// 		if err != nil {
// 			break
// 		}
// 		nextBlock = append(nextBlock, rawTrace)
// 	}
// 	return nextBlock, err
// }
