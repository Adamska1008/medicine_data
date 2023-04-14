package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type TransactionRecord struct {
	Retailer Retailer `json:"retailer"`
	// TODO
}

type Retailer struct {
	Name        string `json:"name"`
	License     string `json:"license"`
	LegalPerson string `json:"legal_person"`
}

type TransactionRecordContract struct {
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
	case "queryHistory":
		return t.queryHistory()
	default:
		return sdk.Error("invalid method")
	}
}

func (t *TransactionRecordContract) save() protogo.Response {
	// TODO
}

func (t *TransactionRecordContract) queryHistory() protogo.Response {
	// TODO
}
