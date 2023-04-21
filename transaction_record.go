package main

import (
	"encoding/json"
	"fmt"
	"time"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type TransactionRecord struct {
	Retailer     Retailer    `json:"retailer"`
	Medicine     Medicine    `json:"medicine"`
	Responsible  Responsible `json:"resiponsible"`
	DeliveryAddr string      `json:"delivery_address"` // 发货地址
	DeliveryTime time.Time   `json:"delivery_time"`    // 交货时间
	ShippingAddr string      `json:"shipping_addr"`    // 收货地址
	ReceiptTime  time.Time   `json:"receipt_time"`     // 收货时间
	Remark       string      `json:"remark"`           // 备注
}

type Medicine struct {
	Id     int    `json:"id"`
	Batch  string `json:"batch"`
	Number int    `json:"number"`
	Price  int    `json:"price"`
}

type Retailer struct {
	Name        string `json:"name"`
	License     string `json:"license"`
	LegalPerson string `json:"legal_person"`
}

type Responsible struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type TransactionRecordContract struct{}

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
