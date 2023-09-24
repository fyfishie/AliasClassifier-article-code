/*
  - @Author: fyfishie
  - @Date: 2023-05-11:08

* @LastEditors: fyfishie

* @LastEditTime: 2023-05-11:19

  - @@email: fyfishie@outlook.com

  - @Description: we need something that makes sure only one process running at the same time
    even caller give more than one tasks
*/
package aliascek

import (
	"aliasParseMaster/lib"
	"errors"

	"github.com/sirupsen/logrus"
)

type CheckResult string

var ErrCheckResultOk = errors.New("check_result_ok")

type CheckRequest struct {
	ResultOut chan lib.AliasCheckResult
	Task      *lib.AliasCekTask
}
type ThreadSafeChecker struct {
	TaskInChan chan *CheckRequest
}

func NewThreadSafeChecker() *ThreadSafeChecker {
	return &ThreadSafeChecker{TaskInChan: make(chan *CheckRequest)}
}

/*
 * @block: false
 */
func (c *ThreadSafeChecker) Start() {
	go func() {
		for task := range c.TaskInChan {
			task.ResultOut <- c.Check(task)
		}
	}()
}

func (c *ThreadSafeChecker) Check(req *CheckRequest) lib.AliasCheckResult {
	checker := NewChecker(req.Task)
	err := checker.CheckStart()
	if err != nil {
		logrus.Errorf("error while running checker, err = %v\n", err.Error())
		return lib.AliasCheckResult{
			Status:     lib.TASK_STATUS_ERROR,
			ErrMessage: err.Error(),
		}
	}
	return lib.AliasCheckResult{Status: lib.TASK_STATUS_DONE, ErrMessage: lib.TASK_RESULT_MSG_DONE}
}
