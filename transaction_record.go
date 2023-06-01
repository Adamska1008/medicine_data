package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type TransactionRecordContract struct{}

type TransactionRecord struct {
	Sum         int       `json:"sum"`         //  交易总金额
	Cnt         int       `json:"cnt"`         //  交易数量
	BuyTime     time.Time `json:"buyTime"`     //  交易时间
	FromAddress string    `json:"fromAddress"` //  发货地址
	ToAddress   string    `json:"toAddress"`   //  收货地址
	MedicineId  string    `json:"medicineId"`  //  药品id
	BuyerId     string    `json:"buyId"`       //  经销商id
}

func (t *TransactionRecordContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

func (t *TransactionRecordContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

func (t *TransactionRecordContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "save":
		return t.save()
	case "queryById":
		return t.queryById()
	default:
		return sdk.Error("invalid method")
	}
}

func (t *TransactionRecordContract) save() protogo.Response {
	params := sdk.Instance.GetArgs()
	// 获取 id 参数
	transactionId := string(params["transaction_id"])
	// 获取 transaction 字符串
	transactionRecordStr := string(params["transaction_record"])
	var transactionRecord TransactionRecord
	// 反序列化
	err := json.Unmarshal([]byte(transactionRecordStr), &transactionRecord)
	if err != nil {
		errMsg := fmt.Sprintf("unmarshall transaction record failed: %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	// 发送事件
	sdk.Instance.EmitEvent(transactionId, []string{transactionRecordStr})
	// 保存数据
	err = sdk.Instance.PutState("transaction_id", transactionId, transactionRecordStr)
	if err != nil {
		errMsg := fmt.Sprintf("put new transaction record failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	sdk.Instance.Infof("[save] transaction_id = " + transactionId)
	return sdk.Success([]byte(transactionId + transactionRecordStr))
}

func (t *TransactionRecordContract) queryById() protogo.Response {
	id := string(sdk.Instance.GetArgs()["transaction_id"])
	result, err := sdk.Instance.GetStateByte("transaction_id", id)
	if err != nil {
		return sdk.Error("failed to call get_state")
	}
	var record TransactionRecord
	if err = json.Unmarshal(result, &record); err != nil {
		return sdk.Error(fmt.Sprintf("unmarshal record failed, err: %s", err))
	}
	sdk.Instance.Infof("[queryById] record_id = " + id)
	return sdk.Success(result)
}

func main() {
	err := sandbox.Start(new(TransactionRecordContract))
	if err != nil {
		log.Fatal(err)
	}
}
