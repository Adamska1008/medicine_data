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
	Number string `json:"number"`
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
	// case "queryHistory":
	// 	return t.queryHistory()
	default:
		return sdk.Error("invalid method")
	}
}

func (t *TransactionRecordContract) save() protogo.Response {
	params := sdk.Instance.GetArgs()
	transactionId := string(params["transaction_id"])
	transactionRecordStr := string(params["transaction_record"])
	var transactionRecord TransactionRecord
	err := json.Unmarshal([]byte(transactionRecordStr), &transactionRecord)
	if err != nil {
		errMsg := fmt.Sprintf("unmarshall transaction record failed: %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	sdk.Instance.EmitEvent(transactionId, []string{transactionRecordStr})
	err = sdk.Instance.PutState(transactionId, "", transactionRecordStr)
	if err != nil {
		errMsg := fmt.Sprintf("put new transaction record failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	return sdk.Success([]byte(transactionId + transactionRecordStr))
}

// func (t *TransactionRecordContract) queryHistory() protogo.Response {
// 	// TODO
// }
