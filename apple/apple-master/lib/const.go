/*
 * @Author: fyfishie
 * @Date: 2023-02-20:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-12:10
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package lib

const MQ_VP_TASK_EXCHANGE_NAME = "vp_task_exchange"
const MQ_VP_RESULT_QUEUE_BINDKEY = "vp_result_bindkey"
const MQ_VP_RESULT_QUEUE_NAME = "vp_result_consume_queue"
const MQ_VP_RESULT_EXCHANGE_NAME = "vp_result_exchange"
const (
	TASK_STATUS_WAITING = iota
	TASK_STATUS_RUNNING
	TASK_STATUS_DONE
	TASK_STATUS_ERROR
)
const (
	TASK_RESULT_MSG_OK   = "ok"
	TASK_RESULT_MSG_DONE = "done"
)
const (
	MongoDBName_AntiAliasParseResult = "aliasparse_antialias_result"
	MongoDBName_PingResult           = "aliasparse_ping_result"
	MongoDBName_AliasCheckResult     = "aliasparse_alias_check_result"
	MongoDBName_MaybeAliasSet        = "aliasparse_maybe_alias_set"
)
const (
	TASK_TYPE_DETECT_AND_PARSE = "detect_and_parse"
	TASK_TYPE_SLEEP            = "sleep"
)
