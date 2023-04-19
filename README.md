# medicine_data
长安链药品合约，包括购买记录和药品信息两个部分

## Data

### 购买记录
药品经销商上链操作属性包括：
1. 药品经销商单位名称
2. 药品经销商的药品经营许可证
3. 药品经销商法人
4. 购买药品种类
5. 购买药品批次
6. 购买药品数量
7. 购买药品价格（可选）
8. 发货详细地址
9. 交货时间
10. 收货地点
11. 收货时间
12. 责任人姓名
13. 责任人电子邮件
14. 责任人联系电话
15. 备注（可选）

### 药品信息


## Run
使用`./build.sh TransactionRecordContract`编译交易记录智能合约。

## Requires
* Go version 1.20
* 7zip: 编译后的文件需要经过7zip压缩形成最后的合约文件，使用`sudo apt install p7zip-full p7zip-rar`在Ubuntu上安装7zip。
