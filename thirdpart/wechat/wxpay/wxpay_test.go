package wxpay

import (
	"fmt"
)

import "testing"

func TestWxPay(t *testing.T) {

	config := &Config{
		AppId: "wxd83a941d5e18f70f",
		MchId: "1493063982",
		ApiKey: "jishangqazwsxedcrphoenixfvtgbyhnujmikol",
		CertFile: "cert/apiclient_cert.pem",
		KeyFile: "cert/apiclient_key.pem",
		RootCAFile: "cert/rootca.pem",
		NotifyUrl: "http://mall-api.uwely.com/notify/pay/wxpub",
	}

	c := NewWxPay(config)

	ret, err := c.QueryTransfer("abcdefg")
	fmt.Println("查询企业转账订单信息", ret, err)

	// ret, err = c.Charge("o_s5ks3Vydm7R8FyatauKR4DGPlI", "100", "abcdefg", "测试", DEVICE_INFO_WEB, TRADE_TYPE_JSAPI, "127.0.0.1")
	// fmt.Println("创建支付单", ret, err)

	// ret, err = c.Refund("pay_transaction_id_001", "100", "f_abcdefg", "100")
	// fmt.Println("创建退款单", ret, err)

	ret, err = c.QueryCharge("abcdefg")
	fmt.Println("查询支付订单", ret, err)

}
