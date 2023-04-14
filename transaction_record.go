package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type TransactionRecord struct {
}

type TransactionRecordContract struct {
}

func (m *TransactionRecordContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

func (m *TransactionRecordContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

func (m *TransactionRecordContract) InvokeContract(method string) protogo.Response {
	switch method {
	default:
		return sdk.Error("invalid method")
	}
}
