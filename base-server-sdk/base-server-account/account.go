package base_server_account

import (
	"encoding/json"
	common "github.com/becent/golang-common"
	"github.com/becent/golang-common/base-server-sdk"
	"strconv"
	"strings"
)

type OpType int
type AccountStatus int

const (
	OP_TYPE_AVAIL_ADD  OpType = 1 //可用-加
	OP_TYPE_AVAIL_SUB  OpType = 2 //可用-减
	OP_TYPE_FREEZE_ADD OpType = 3 //冻结-加
	OP_TYPE_FREEZE_SUB OpType = 4 //冻结-减
	OP_TYPE_UN_FREEZE  OpType = 5 //解冻-冻结进可用

	ACCOUNT_STATUS_NORMAL AccountStatus = 1 //正常
	ACCOUNT_STATUS_FREEZE AccountStatus = 2 //冻结
)

type Account struct {
	AccountId    int64  `json:"accountId"`
	OrgId        int    `json:"orgId"`
	UserId       int64  `json:"userId"`
	Currency     string `json:"currency"`
	AvailAmount  string `json:"availAmount"`
	FreezeAmount string `json:"freezeAmount"`
	Status       int    `json:"status"`
	CreateTime   int64  `json:"createTime"`
	UpdateTime   int64  `json:"updateTime"`
}

type LogList struct {
	LogId      int64  `json:"logId"`
	UserId     int64  `json:"userId"`
	Currency   string `json:"currency"`
	LogType    int    `json:"logType"`
	Amount     string `json:"amount"`
	CreateTime int64  `json:"createTime"`
}

type TaskDetail struct {
	OpType        OpType `json:"opType"`
	BsType        int    `json:"bsType"`
	AccountId     int64  `json:"accountId"`
	UserId        int64  `json:"userId"`
	Currency      string `json:"currency"`
	AllowNegative int    `json:"allowNegative"`
	Amount        string `json:"amount"`
	Detail        string `json:"detail"`
	Ext           string `json:"ext"`
}

type TaskCallBack struct {
	CallBackUrl string            `json:"callBackUrl"`
	Data        map[string]string `json:"data"`
}

//  POST CreateAccount 创建账户
//
//	注意:
//	1. orgId必须大于0
//
//	异常错误:
//	1001 参数错误
//	2001 账户已存在
//	2002 账户创建失败
func CreateAccount(orgId int, userId int64, currency []string) ([]*Account, *base_server_sdk.Error) {
	if orgId <= 0 || userId <= 0 || len(currency) == 0 {
		return nil, base_server_sdk.ErrInvalidParams
	}
	params := make(map[string]string)
	params["orgId"] = strconv.Itoa(orgId)
	params["userId"] = strconv.FormatInt(userId, 10)
	params["currency"] = strings.Join(currency, ",")

	client := base_server_sdk.Instance
	data, err := client.DoRequest(client.Hosts.AccountServerHost, "account", "createAccount", params)
	if err != nil {
		return nil, err
	}
	var account []*Account
	if err := json.Unmarshal(data, &account); err != nil {
		common.ErrorLog("baseServerSdk_CreateAccount", params, "unmarshal account fail"+string(data))
		return nil, base_server_sdk.ErrServiceBusy
	}
	return account, nil
}

//	账户信息
//	POST account/AccountInfo
//
//	异常错误:
//	1001 参数错误
//	2003 账户不存在
func AccountInfo(orgId int, userIds []int64, currency string) ([]*Account, *base_server_sdk.Error) {
	if orgId <= 0 || len(userIds) <= 0 {
		return nil, base_server_sdk.ErrInvalidParams
	}
	params := make(map[string]string)
	params["orgId"] = strconv.Itoa(orgId)
	userIdsMarshal, _ := json.Marshal(userIds)
	params["userIds"] = string(userIdsMarshal)
	params["currency"] = currency

	client := base_server_sdk.Instance
	data, err := client.DoRequest(client.Hosts.AccountServerHost, "account", "accountInfo", params)
	if err != nil {
		return nil, err
	}

	var account []*Account
	if err := json.Unmarshal(data, &account); err != nil {
		common.ErrorLog("baseServerSdk_AccountInfo", params, "unmarshal account fail"+string(data))
		return nil, base_server_sdk.ErrServiceBusy
	}
	return account, nil
}

//	状态变更
//	POST account/updateStatus
//
//	status 状态 1:正常 2:禁用
//
//	异常错误:
//	1001 参数错误
//	2003 账户不存在
//	2004 更新状态失败
func UpdateStatus(orgId int, accountId int64, status AccountStatus) *base_server_sdk.Error {
	if orgId <= 0 || accountId <= 0 || status <= 0 {
		return base_server_sdk.ErrInvalidParams
	}
	params := make(map[string]string)
	params["orgId"] = strconv.Itoa(orgId)
	params["accountId"] = strconv.FormatInt(accountId, 10)
	params["status"] = strconv.Itoa(int(status))

	client := base_server_sdk.Instance
	_, err := client.DoRequest(client.Hosts.AccountServerHost, "account", "updateStatus", params)
	if err != nil {
		return err
	}
	return nil
}

//  金额操作
//  POST account/operateAmount
//	类型枚举:
//	1	//可用-加
//	2	//可用-减
//	3	//冻结-加
//	4	//冻结-减
//	5	//解冻-冻结进可用
//
//	异常错误:
//	1001 参数错误
//	2003 账户不存在
//	1009 BC操作失败
//	2005 账户可用增加失败
//	2007 可用余额不足
//	2008 解冻失败
//	2009 账户可用减少失败
//	2010 账户冻结减少失败
//	2011 账户日志创建失败
func OperateAmount(orgId int, accountId int64, opType OpType, bsType, allowNegative int, amount, detail, ext, callback string) *base_server_sdk.Error {
	if orgId <= 0 || accountId <= 0 || opType <= 0 || bsType <= 0 || amount == "" {
		return base_server_sdk.ErrInvalidParams
	}
	params := make(map[string]string)
	params["orgId"] = strconv.Itoa(orgId)
	params["accountId"] = strconv.FormatInt(accountId, 10)
	params["opType"] = strconv.Itoa(int(opType))
	params["bsType"] = strconv.Itoa(bsType)
	params["allowNegative"] = strconv.Itoa(allowNegative)
	params["amount"] = amount
	params["detail"] = detail
	params["ext"] = ext
	params["callback"] = callback

	client := base_server_sdk.Instance
	_, err := client.DoRequest(client.Hosts.AccountServerHost, "account", "operateAmount", params)
	if err != nil {
		return err
	}
	return nil
}

// 账户日志列表
// post account/accountLogList
//
//	类型枚举:
//	1	可用-加
//	2	可用-减
//	3	冻结-加
//	4	冻结-减
//	5	解冻-冻结进可用
//
//	异常错误:
//	1001 参数错误
//	2003 账户不存在
func AccountLogList(orgId int, userId int64, opType OpType, bsType int, currency string, beginTime, endTime int, page, limit int) ([]*LogList, *base_server_sdk.Error) {
	if orgId <= 0 || userId <= 0 || opType < 0 || page <= 0 || limit <= 0 || limit > 1000 {
		return nil, base_server_sdk.ErrInvalidParams
	}

	params := make(map[string]string)
	params["orgId"] = strconv.Itoa(orgId)
	params["userId"] = strconv.FormatInt(userId, 10)
	params["opType"] = strconv.Itoa(int(opType))
	params["bsType"] = strconv.Itoa(bsType)
	params["currency"] = currency
	params["beginTime"] = strconv.Itoa(beginTime)
	params["endTime"] = strconv.Itoa(endTime)
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	client := base_server_sdk.Instance
	data, err := client.DoRequest(client.Hosts.AccountServerHost, "account", "accountLogList", params)
	if err != nil {
		return nil, err
	}

	var logList []*LogList
	if err := json.Unmarshal(data, &logList); err != nil {
		common.ErrorLog("baseServerSdk_AccountLogList", params, "unmarshal account list fail"+string(data))
		return nil, base_server_sdk.ErrServiceBusy
	}
	return logList, nil
}

// 日志总和
func SumLog(orgId int, userId int64, opType OpType, bsType int, currency string, beginTime, endTime int) (string, *base_server_sdk.Error) {
	if orgId <= 0 || userId <= 0 || opType < 0 {
		return "0", base_server_sdk.ErrInvalidParams
	}
	params := make(map[string]string)
	params["orgId"] = strconv.Itoa(orgId)
	params["userId"] = strconv.FormatInt(userId, 10)
	params["opType"] = strconv.Itoa(int(opType))
	params["bsType"] = strconv.Itoa(bsType)
	params["currency"] = currency
	params["beginTime"] = strconv.Itoa(beginTime)
	params["endTime"] = strconv.Itoa(endTime)

	client := base_server_sdk.Instance
	data, err := client.DoRequest(client.Hosts.AccountServerHost, "account", "sumLog", params)
	if err != nil {
		return "0", err
	}

	var sumAmount string
	if err := json.Unmarshal(data, &sumAmount); err != nil {
		common.ErrorLog("baseServerSdk_SumLog", params, "unmarshal sumAmount fail"+string(data))
		return "0", base_server_sdk.ErrServiceBusy
	}

	return sumAmount, nil
}

// 批量操作金额
func BatchOperateAmount(orgId, isAsync int, details []*TaskDetail, callback *TaskCallBack) *base_server_sdk.Error {
	if orgId <= 0 || len(details) == 0 {
		return base_server_sdk.ErrInvalidParams
	}
	params := make(map[string]string)
	params["orgId"] = strconv.Itoa(orgId)
	params["isAsync"] = strconv.Itoa(isAsync)
	taskDetailByte, _ := json.Marshal(details)
	params["detail"] = string(taskDetailByte)
	callbackData, _ := json.Marshal(callback)
	params["callback"] = string(callbackData)

	client := base_server_sdk.Instance
	_, err := client.DoRequest(client.Hosts.AccountServerHost, "account", "batchOperateAmount", params)
	if err != nil {
		return err
	}
	return nil
}
