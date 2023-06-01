package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type MedicineDataContract struct {
}

type MedicineData struct {
	Name              string    `json:"name"`               // 药品名称
	Company           string    `json:"company"`            // 药品公司
	Licence           string    `json:"license"`            // 文书
	Specification     string    `json:"specification"`      // 规格
	Material          string    `json:"material"`           // 原料
	Batch             string    `json:"batch"`              // 生产批次
	ProductionAddr    string    `json:"production_address"` // 生产地址
	ProducerId        int       `json:"producer_id"`        // 生产商id
	BatchNumber       int       `json:"batch_number"`       // 批次数量
	PackingFrim       string    `json:"packing_firm"`       // 包装企业
	Gmp               string    `json:"gmp"`                // GMP 标号
	Responsible       string    `json:"responsible_person"` // 负责人
	Status            string    `json:"status"`             // 药品状态
	AdminId           string    `json:"admin_id"`           // 批准审核人的ID
	RemainingQuantity int       `json:"remaining_quantity"` // 剩余数量
	ProductionDate    time.Time `json:"production_date"`    // 生产时间
	ExpiredDate       time.Time `json:"expired_date"`       // 过期时间
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
	case "buy":
		return m.buy()
	case "changeStatus":
		return m.changeStatus()
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

func (m *MedicineDataContract) buy() protogo.Response {
	id := string(sdk.Instance.GetArgs()["medicine_id"])
	number, err := strconv.ParseInt(string(sdk.Instance.GetArgs()["number"]), 10, 32)
	if err != nil {
		return sdk.Error("the arg 'number' is not int")
	}
	result, err := sdk.Instance.GetStateByte("medicine_id", id)
	if err != nil {
		return sdk.Error("failed to call get_state")
	}
	var medicine MedicineData
	if err = json.Unmarshal(result, &medicine); err != nil {
		return sdk.Error(fmt.Sprintf("unmarshal record failed, err: %s", err))
	}
	medicine.RemainingQuantity -= int(number)
	medicineStr, _ := json.Marshal(medicine)
	// 发送事件
	sdk.Instance.EmitEvent(id, []string{string(medicineStr)})
	// 保存数据
	err = sdk.Instance.PutState("medicine_id", id, string(medicineStr))
	if err != nil {
		errMsg := fmt.Sprintf("put new transaction record failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	sdk.Instance.Infof("[update] medicine_id = " + id + " number " + strconv.Itoa(int(number)))
	return sdk.Success([]byte(id + string(medicineStr)))
}

func (m *MedicineDataContract) changeStatus() protogo.Response {
	id := string(sdk.Instance.GetArgs()["medicine_id"])
	status := string(sdk.Instance.GetArgs()["status"])

	result, err := sdk.Instance.GetStateByte("medicine_id", id)
	if err != nil {
		return sdk.Error("failed to call get_state")
	}
	var medicine MedicineData
	if err = json.Unmarshal(result, &medicine); err != nil {
		return sdk.Error(fmt.Sprintf("unmarshal record failed, err: %s", err))
	}
	medicine.Status = status
	medicineStr, _ := json.Marshal(medicine)
	// 发送事件
	sdk.Instance.EmitEvent(id, []string{string(medicineStr)})
	// 保存数据
	err = sdk.Instance.PutState("medicine_id", id, string(medicineStr))
	if err != nil {
		errMsg := fmt.Sprintf("put new transaction record failed, %s", err)
		sdk.Instance.Errorf(errMsg)
		return sdk.Error(errMsg)
	}
	sdk.Instance.Infof("[update] medicine_id = " + id + " status " + status)
	return sdk.Success([]byte(id + string(medicineStr)))
}

func main() {
	err := sandbox.Start(new(MedicineDataContract))
	if err != nil {
		log.Fatal(err)
	}
}
