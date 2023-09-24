/*
 * @Author: fyfishie
 * @Date: 2023-03-03:11
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-08:10
 * @Description: run sping detect only, should not have any chance contacting with master
 * @email: muren.zhuang@outlook.com
 */
package spingdetct

import (
	"os"
	"os/exec"
	"vp/status"

	"github.com/sirupsen/logrus"
)

type Detector struct {
	ipAddrList []string
	spingHome  string
	StatusEye status.StatusEye
}

/*
 * @description: @IPList: ip format or CIDR format
 */
type SpingDetectTask struct {
	IPAddrList []string
	SpingHome  string
}

func NewDetector(ipAddrList []string, spingHome string) *Detector {
	d := Detector{
		ipAddrList: ipAddrList,
		spingHome:  spingHome,
		StatusEye: status.StatusEye{},
	}
	return &d
}

func (d *Detector) Run() error {
	d.StatusEye.Status = "writing ip addr into sping input file..."
	err := d.write2SpingInput()
	if err != nil {
		return err
	}
	d.StatusEye.Status = "start sping and running detect..."
	cmd := exec.Command(d.spingHome+"/sping4", "-if", d.spingHome+"/sping_input", "-of", d.spingHome+"/sping_output")
	output, err := cmd.CombinedOutput()
	d.StatusEye.Status = "done"
	if err != nil {
		logrus.Error("error in exec sping4, err = %v\noutput = %v\n", err.Error(), string(output))
	}
	return err
}

/*
@description: write ip addr(point or CIDR) into sping input file,
if any error encounted while open the file or writing into it,
this function stops and return the error

each time the function called, it clears all content before writing ip addr into.

@return {error}
*/
func (d *Detector) write2SpingInput() error {
	fi, err := os.OpenFile(d.spingHome+"/sping_input", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Printf("error in open sping input file, err = %v\n", err.Error())
		return err
	}
	defer fi.Close()
	for _, addr := range d.ipAddrList {
		_, err := fi.WriteString(addr + "\n")
		if err != nil {
			logrus.Error("error in write addr into sping input file, err = %v\n", err.Error())
			return err
		}
	}
	return err
}

/*
* @description: read sping output, if any error encounted while read lien from output file,
stop reading file when io.EOF it is, skip the line when others.

if failed to open the output file, it returns the error untouched.
* @return {*}
*/
// func (d *Detector) readSpingOutput() ([]lib.RawSping, error) {
// 	fi, err := os.OpenFile(d.spingHome+"/sping_output", os.O_RDONLY, 0444)
// 	if err != nil {
// 		logrus.Printf("error in open sping output file, err = %v\n", err.Error())
// 		return nil, err
// 	}
// 	defer fi.Close()
// 	rdr := bufio.NewReader(fi)
// 	res := []lib.RawSping{}
// 	for {
// 		line, _, err := rdr.ReadLine()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			} else {
// 				continue
// 			}
// 		}
// 		rawSping := lib.RawSping{}
// 		err = json.Unmarshal(line, &rawSping)
// 		if err != nil {
// 			logrus.Printf("error in unmarshal data bytes to lib.RawSping{}, err = %v\n", err.Error())
// 			continue
// 		}
// 		res = append(res, rawSping)
// 	}
// 	return res, nil
// }
