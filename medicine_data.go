package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type MedicineData struct {
}

type MedicineDataContract struct {
}

func (m *MedicineDataContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

func (m *MedicineDataContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

func (m *MedicineDataContract) InvokeContract(method string) protogo.Response {
	switch method {
	default:
		return sdk.Error("invalid method")
	}
}
