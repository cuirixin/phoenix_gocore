package wechat

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/cuirixin/phoenix_gocore/thirdpart/gopay"
)

type Client struct {
	AppId       string
	MchId       string
	ApiKey      string
	BaseURL     string
	IsProd      bool
	certificate tls.Certificate
	certPool    *x509.CertPool
	mu          sync.RWMutex
}

// 初始化微信客户端
//    appId：应用ID
//    mchId：商户ID
//    ApiKey：API秘钥值
//    IsProd：是否是正式环境
func NewClient(appId, mchId, apiKey string, isProd bool) (client *Client) {
	return &Client{
		AppId:  appId,
		MchId:  mchId,
		ApiKey: apiKey,
		IsProd: isProd}
}

// 提交付款码支付
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_10&index=1
func (w *Client) Micropay(bm gopay.BodyMap) (wxRsp *MicropayResponse, err error) {
	err = bm.CheckEmptyError("nonce_str", "body", "out_trade_no", "total_fee", "spbill_create_ip", "auth_code")
	if err != nil {
		return nil, err
	}
	var bs []byte
	if w.IsProd {
		bs, err = w.doWeChatPostProd(bm, microPay, nil)
	} else {
		bm.Set("total_fee", 1)
		bs, err = w.doWeChatPostSanBox(bm, sandboxMicroPay)
	}
	if err != nil {
		return nil, err
	}
	wxRsp = new(MicropayResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 授权码查询openid（正式）
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_13&index=9
func (w *Client) AuthCodeToOpenId(bm gopay.BodyMap) (wxRsp *AuthCodeToOpenIdResponse, err error) {
	err = bm.CheckEmptyError("nonce_str", "auth_code")
	if err != nil {
		return nil, err
	}

	bs, err := w.doWeChatPostProd(bm, authCodeToOpenid, nil)
	if err != nil {
		return nil, err
	}
	wxRsp = new(AuthCodeToOpenIdResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 统一下单
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (w *Client) UnifiedOrder(bm gopay.BodyMap) (wxRsp *UnifiedOrderResponse, err error) {
	err = bm.CheckEmptyError("nonce_str", "body", "out_trade_no", "total_fee", "spbill_create_ip", "notify_url", "trade_type")
	if err != nil {
		return nil, err
	}
	var bs []byte
	if w.IsProd {
		bs, err = w.doWeChatPostProd(bm, unifiedOrder, nil)
	} else {
		bm.Set("total_fee", 101)
		bs, err = w.doWeChatPostSanBox(bm, sandboxUnifiedOrder)
	}
	if err != nil {
		return nil, err
	}
	wxRsp = new(UnifiedOrderResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 查询订单
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
func (w *Client) QueryOrder(bm gopay.BodyMap) (wxRsp *QueryOrderResponse, err error) {
	err = bm.CheckEmptyError("nonce_str")
	if err != nil {
		return nil, err
	}
	if bm.Get("out_trade_no") == gopay.NULL && bm.Get("transaction_id") == gopay.NULL {
		return nil, errors.New("out_trade_no and transaction_id are not allowed to be null at the same time")
	}
	var bs []byte
	if w.IsProd {
		bs, err = w.doWeChatPostProd(bm, orderQuery, nil)
	} else {
		bs, err = w.doWeChatPostSanBox(bm, sandboxOrderQuery)
	}
	if err != nil {
		return nil, err
	}
	wxRsp = new(QueryOrderResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 关闭订单
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
func (w *Client) CloseOrder(bm gopay.BodyMap) (wxRsp *CloseOrderResponse, err error) {
	err = bm.CheckEmptyError("nonce_str", "out_trade_no")
	if err != nil {
		return nil, err
	}
	var bs []byte
	if w.IsProd {
		bs, err = w.doWeChatPostProd(bm, closeOrder, nil)
	} else {
		bs, err = w.doWeChatPostSanBox(bm, sandboxCloseOrder)
	}
	if err != nil {
		return nil, err
	}
	wxRsp = new(CloseOrderResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 撤销订单
//    注意：如已使用client.AddCertFilePath()添加过证书，参数certFilePath、keyFilePath、pkcs12FilePath全传空字符串 ""，否则，3证书Path均不可空
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_11&index=3
func (w *Client) Reverse(bm gopay.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRsp *ReverseResponse, err error) {
	err = bm.CheckEmptyError("nonce_str", "out_trade_no")
	if err != nil {
		return nil, err
	}
	var (
		bs        []byte
		tlsConfig *tls.Config
	)
	if w.IsProd {
		if tlsConfig, err = w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath); err != nil {
			return nil, err
		}
		bs, err = w.doWeChatPostProd(bm, reverse, tlsConfig)
	} else {
		bs, err = w.doWeChatPostSanBox(bm, sandboxReverse)
	}
	if err != nil {
		return nil, err
	}
	wxRsp = new(ReverseResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 申请退款
//    注意：如已使用client.AddCertFilePath()添加过证书，参数certFilePath、keyFilePath、pkcs12FilePath全传空字符串 ""，否则，3证书Path均不可空
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
func (w *Client) Refund(bm gopay.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRsp *RefundResponse, err error) {
	err = bm.CheckEmptyError("nonce_str", "out_refund_no", "total_fee", "refund_fee")
	if err != nil {
		return nil, err
	}
	if bm.Get("out_trade_no") == gopay.NULL && bm.Get("transaction_id") == gopay.NULL {
		return nil, errors.New("out_trade_no and transaction_id are not allowed to be null at the same time")
	}
	var (
		bs        []byte
		tlsConfig *tls.Config
	)
	if w.IsProd {
		if tlsConfig, err = w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath); err != nil {
			return nil, err
		}
		bs, err = w.doWeChatPostProd(bm, refund, tlsConfig)
	} else {
		bs, err = w.doWeChatPostSanBox(bm, sandboxRefund)
	}
	if err != nil {
		return nil, err
	}
	wxRsp = new(RefundResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 查询退款
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_5
func (w *Client) QueryRefund(bm gopay.BodyMap) (wxRsp *QueryRefundResponse, err error) {
	err = bm.CheckEmptyError("nonce_str")
	if err != nil {
		return nil, err
	}
	if bm.Get("refund_id") == gopay.NULL && bm.Get("out_refund_no") == gopay.NULL && bm.Get("transaction_id") == gopay.NULL && bm.Get("out_trade_no") == gopay.NULL {
		return nil, errors.New("refund_id, out_refund_no, out_trade_no, transaction_id are not allowed to be null at the same time")
	}
	var bs []byte
	if w.IsProd {
		bs, err = w.doWeChatPostProd(bm, refundQuery, nil)
	} else {
		bs, err = w.doWeChatPostSanBox(bm, sandboxRefundQuery)
	}
	if err != nil {
		return nil, err
	}
	wxRsp = new(QueryRefundResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 下载对账单
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_6
func (w *Client) DownloadBill(bm gopay.BodyMap) (wxRsp string, err error) {
	err = bm.CheckEmptyError("nonce_str", "bill_date", "bill_type")
	if err != nil {
		return gopay.NULL, err
	}
	billType := bm.Get("bill_type")
	if billType != "ALL" && billType != "SUCCESS" && billType != "REFUND" && billType != "RECHARGE_REFUND" {
		return gopay.NULL, errors.New("bill_type error, please reference: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_6")
	}
	var bs []byte
	if w.IsProd {
		bs, err = w.doWeChatPostProd(bm, downloadBill, nil)
	} else {
		bs, err = w.doWeChatPostSanBox(bm, sandboxDownloadBill)
	}
	if err != nil {
		return gopay.NULL, err
	}
	return string(bs), nil
}

// 下载资金账单（正式）
//    注意：如已使用client.AddCertFilePath()添加过证书，参数certFilePath、keyFilePath、pkcs12FilePath全传空字符串 ""，否则，3证书Path均不可空
//    貌似不支持沙箱环境，因为沙箱环境默认需要用MD5签名，但是此接口仅支持HMAC-SHA256签名
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_18&index=7
func (w *Client) DownloadFundFlow(bm gopay.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRsp string, err error) {
	err = bm.CheckEmptyError("nonce_str", "bill_date", "account_type")
	if err != nil {
		return gopay.NULL, err
	}
	accountType := bm.Get("account_type")
	if accountType != "Basic" && accountType != "Operation" && accountType != "Fees" {
		return gopay.NULL, errors.New("account_type error, please reference: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_18&index=7")
	}
	bm.Set("sign_type", SignType_HMAC_SHA256)
	tlsConfig, err := w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath)
	if err != nil {
		return gopay.NULL, err
	}
	bs, err := w.doWeChatPostProd(bm, downloadFundFlow, tlsConfig)
	if err != nil {
		return gopay.NULL, err
	}
	wxRsp = string(bs)
	return
}

// 交易保障
//    文档地址：（JSAPI）https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_8&index=9
//    文档地址：（付款码）https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_14&index=8
//    文档地址：（Native）https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_8&index=9
//    文档地址：（APP）https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_8&index=10
//    文档地址：（H5）https://pay.weixin.qq.com/wiki/doc/api/H5.php?chapter=9_8&index=9
//    文档地址：（微信小程序）https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_8&index=9
func (w *Client) Report(bm gopay.BodyMap) (wxRsp *ReportResponse, err error) {
	err = bm.CheckEmptyError("nonce_str", "interface_url", "execute_time", "return_code", "return_msg", "result_code", "user_ip")
	if err != nil {
		return nil, err
	}
	var bs []byte
	if w.IsProd {
		bs, err = w.doWeChatPostProd(bm, report, nil)
	} else {
		bs, err = w.doWeChatPostSanBox(bm, sandboxReport)
	}
	if err != nil {
		return nil, err
	}
	wxRsp = new(ReportResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 拉取订单评价数据（正式）
//    注意：如已使用client.AddCertFilePath()添加过证书，参数certFilePath、keyFilePath、pkcs12FilePath全传空字符串 ""，否则，3证书Path均不可空
//    貌似不支持沙箱环境，因为沙箱环境默认需要用MD5签名，但是此接口仅支持HMAC-SHA256签名
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_17&index=11
func (w *Client) BatchQueryComment(bm gopay.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRsp string, err error) {
	err = bm.CheckEmptyError("nonce_str", "begin_time", "end_time", "offset")
	if err != nil {
		return gopay.NULL, err
	}
	bm.Set("sign_type", SignType_HMAC_SHA256)
	tlsConfig, err := w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath)
	if err != nil {
		return gopay.NULL, err
	}
	bs, err := w.doWeChatPostProd(bm, batchQueryComment, tlsConfig)
	if err != nil {
		return gopay.NULL, err
	}
	return string(bs), nil
}

// 企业向微信用户个人付款（正式）
//    注意：如已使用client.AddCertFilePath()添加过证书，参数certFilePath、keyFilePath、pkcs12FilePath全传空字符串 ""，否则，3证书Path均不可空
//    注意：此方法未支持沙箱环境，默认正式环境，转账请慎重
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_2
func (w *Client) Transfer(bm gopay.BodyMap, certFilePath, keyFilePath, pkcs12FilePath string) (wxRsp *TransfersResponse, err error) {
	err = bm.CheckEmptyError("nonce_str", "partner_trade_no", "openid", "check_name", "amount", "desc", "spbill_create_ip")
	if err != nil {
		return nil, err
	}
	bm.Set("mch_appid", w.AppId)
	bm.Set("mchid", w.MchId)
	var (
		tlsConfig *tls.Config
		url       = baseUrlCh + transfers
	)
	if tlsConfig, err = w.addCertConfig(certFilePath, keyFilePath, pkcs12FilePath); err != nil {
		return nil, err
	}
	bm.Set("sign", getReleaseSign(w.ApiKey, SignType_MD5, bm))

	httpClient := gopay.NewHttpClient().SetTLSConfig(tlsConfig).Type(gopay.TypeXML)
	if w.BaseURL != gopay.NULL {
		w.mu.RLock()
		url = w.BaseURL + transfers
		w.mu.RUnlock()
	}
	wxRsp = new(TransfersResponse)
	res, errs := httpClient.Post(url).SendString(generateXml(bm)).EndStruct(wxRsp)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %d", res.StatusCode)
	}
	return wxRsp, nil
}

// 公众号纯签约（正式）
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_1&index=1
func (w *Client) EntrustPublic(bm gopay.BodyMap) (wxRsp *EntrustPublicResponse, err error) {
	err = bm.CheckEmptyError("plan_id", "contract_code", "request_serial", "contract_display_account", "notify_url", "version", "timestamp")
	if err != nil {
		return nil, err
	}
	bs, err := w.doWeChatGetProd(bm, entrustPublic, SignType_MD5)
	if err != nil {
		return nil, err
	}
	wxRsp = new(EntrustPublicResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// APP纯签约-预签约接口-获取预签约ID（正式）
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_5&index=2
func (w *Client) EntrustAppPre(bm gopay.BodyMap) (wxRsp *EntrustAppPreResponse, err error) {
	err = bm.CheckEmptyError("plan_id", "contract_code", "request_serial", "contract_display_account", "notify_url", "version", "timestamp")
	if err != nil {
		return nil, err
	}
	bs, err := w.doWeChatPostProd(bm, entrustApp, nil)
	if err != nil {
		return nil, err
	}
	wxRsp = new(EntrustAppPreResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// H5纯签约（正式）
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_16&index=4
func (w *Client) EntrustH5(bm gopay.BodyMap) (wxRsp *EntrustH5Response, err error) {
	err = bm.CheckEmptyError("plan_id", "contract_code", "request_serial", "contract_display_account", "notify_url", "version", "timestamp", "clientip")
	if err != nil {
		return nil, err
	}
	bs, err := w.doWeChatGetProd(bm, entrustH5, SignType_HMAC_SHA256)
	if err != nil {
		return nil, err
	}
	wxRsp = new(EntrustH5Response)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// 支付中签约（正式）
//    文档地址：https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_13&index=5
func (w *Client) EntrustPaying(bm gopay.BodyMap) (wxRsp *EntrustPayingResponse, err error) {
	err = bm.CheckEmptyError("contract_mchid", "contract_appid",
		"out_trade_no", "nonce_str", "body", "notify_url", "total_fee",
		"spbill_create_ip", "trade_type", "plan_id", "contract_code",
		"request_serial", "contract_display_account", "contract_notify_url")
	if err != nil {
		return nil, err
	}
	bs, err := w.doWeChatPostProd(bm, entrustPaying, nil)
	if err != nil {
		return nil, err
	}
	wxRsp = new(EntrustPayingResponse)
	if err = xml.Unmarshal(bs, wxRsp); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal(%s)：%w", string(bs), err)
	}
	return wxRsp, nil
}

// Post请求
func (w *Client) doWeChatPostSanBox(bm gopay.BodyMap, path string) (bs []byte, err error) {
	var url = baseUrlCh + path
	w.mu.RLock()
	defer w.mu.RUnlock()
	bm.Set("appid", w.AppId)
	bm.Set("mch_id", w.MchId)

	if bm.Get("sign") == gopay.NULL {
		bm.Set("sign_type", SignType_MD5)
		sign, err := getSignBoxSign(w.MchId, w.ApiKey, bm)
		if err != nil {
			return nil, err
		}
		bm.Set("sign", sign)
	}

	if w.BaseURL != gopay.NULL {
		url = w.BaseURL + path
	}
	res, bs, errs := gopay.NewHttpClient().Type(gopay.TypeXML).Post(url).SendString(generateXml(bm)).EndBytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %d", res.StatusCode)
	}
	if strings.Contains(string(bs), "HTML") || strings.Contains(string(bs), "html") {
		return nil, errors.New(string(bs))
	}
	return bs, nil
}

// Post请求、正式
func (w *Client) doWeChatPostProd(bm gopay.BodyMap, path string, tlsConfig *tls.Config) (bs []byte, err error) {
	var url = baseUrlCh + path
	w.mu.RLock()
	defer w.mu.RUnlock()
	bm.Set("appid", w.AppId)
	bm.Set("mch_id", w.MchId)

	if bm.Get("sign") == gopay.NULL {
		sign := getReleaseSign(w.ApiKey, bm.Get("sign_type"), bm)
		bm.Set("sign", sign)
	}

	httpClient := gopay.NewHttpClient()
	if w.IsProd && tlsConfig != nil {
		httpClient.SetTLSConfig(tlsConfig)
	}
	if w.BaseURL != gopay.NULL {
		url = w.BaseURL + path
	}
	res, bs, errs := httpClient.Type(gopay.TypeXML).Post(url).SendString(generateXml(bm)).EndBytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %d", res.StatusCode)
	}
	if strings.Contains(string(bs), "HTML") || strings.Contains(string(bs), "html") {
		return nil, errors.New(string(bs))
	}
	return bs, nil
}

// Get请求、正式
func (w *Client) doWeChatGetProd(bm gopay.BodyMap, path, signType string) (bs []byte, err error) {
	var url = baseUrlCh + path
	w.mu.RLock()
	defer w.mu.RUnlock()
	bm.Set("appid", w.AppId)
	bm.Set("mch_id", w.MchId)
	bm.Remove("sign")
	sign := getReleaseSign(w.ApiKey, signType, bm)
	bm.Set("sign", sign)

	if w.BaseURL != gopay.NULL {
		w.mu.RLock()
		url = w.BaseURL + path
		w.mu.RUnlock()
	}
	param := bm.EncodeGetParams()
	url = url + "?" + param
	res, bs, errs := gopay.NewHttpClient().Get(url).EndBytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %d", res.StatusCode)
	}
	if strings.Contains(string(bs), "HTML") || strings.Contains(string(bs), "html") {
		return nil, errors.New(string(bs))
	}
	return bs, nil
}
