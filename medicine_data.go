package main

import (
	"encoding/json"
	"fmt"
	"time"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type MedicineData struct {
	TraceNo     string `json:traceno`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Licence           string    `json:"license"`
	Specification     string    `json:"specification"`      // 规格
	Material          string    `json:"material"`           // 原料
	Batch             string    `json:"batch"`              // 生产批次
	ProductionAddr    string    `json:"production_address"` // 生产地址
	ProducerId        int       `json:"producer_id"`        // 生产商id
	BatchNumber       int       `json:"batch_number"`       // 批次数量
	SalePrice         float64   `json:"sale_price"`         // 出售价格
	PackingFrim       string    `json:"packing_firm"`       // 包装企业
	Gmp               string    `json:"gmp"`                // GMP 标号
	Responsible       string    `json:"responsible_person"` // 负责人
	Remark            string    `json:"remark"`             // 备注
	Status            string    `json:"status"`             // 药品状态
	AdminId           string    `json:"admin_id"`           // 批准审核人的ID
	RemainingQuantity int       `json:"remaining_quantity"` // 剩余数量
	ProductionDate    time.Time `json:"production_date"`    // 生产时间
	ExpiredDate       time.Time `json:"expired_date"`       // 过期时间

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
	case "save":
		return m.save()
	case "queryById":
		return m.queryById()
	default:
		return sdk.Error("invalid method")
	}
}

func (m *MedicineDataContract) save() protogo.Response {
	params := sdk.Instance.GetArgs()
	// 获取 id 参数
	medicineId := string(params["medicine_id"])
	// 获取 transaction 字符串
	medicineStr := string(params["medicine_data"])
	var medicineData MedicineData
	// 反序列化
	err := json.Unmarshal([]byte(medicineStr), &medicineData)
	if err != nil {
		errMsg := fmt.Sprintf("unmarshall transaction record failed: %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	// 发送事件
	sdk.Instance.EmitEvent(medicineId, []string{medicineStr})
	// 保存数据
	err = sdk.Instance.PutState("medicine_id", medicineId, medicineStr)
	if err != nil {
		errMsg := fmt.Sprintf("put new transaction record failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	sdk.Instance.Infof("[save] medicine_id = " + medicineId)
	return sdk.Success([]byte(medicineId + medicineStr))
}

func (m *MedicineDataContract) queryById() protogo.Response {
	id := string(sdk.Instance.GetArgs()["medicine_id"])
	result, err := sdk.Instance.GetStateByte("medicine_id", id)
	if err != nil {
		return sdk.Error("failed to call get_state")
	}
	var medicine MedicineData
	if err = json.Unmarshal(result, &medicine); err != nil {
		return sdk.Error(fmt.Sprintf("unmarshal record failed, err: %s", err))
	}
	sdk.Instance.Infof("[queryById] medicine_id = " + id)
	return sdk.Success(result)
}
