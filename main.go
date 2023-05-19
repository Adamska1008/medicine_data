package main

import (
	"log"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
)

// sdk代码中，有且仅有一个main()方法
func main() {
	// main()方法中，下面的代码为必须代码，不建议修改main()方法当中的代码
	// 其中，TestContract为用户实现合约的具体名称
	err := sandbox.Start(new(TransactionRecordContract))
	if err != nil {
		log.Fatal(err)
	}
}
