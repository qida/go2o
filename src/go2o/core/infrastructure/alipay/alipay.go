/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-11 22:09
 * description : alipay golang sdk in "github.com/ascoders/alipay"
 * history :
 */
package alipay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var (
	conf_partner string //合作者ID
	conf_key     string //合作者私钥
	conf_seller  string //网站卖家邮箱地址
)

// 成功通知函数
type SuccessNotifyFunc func(r *http.Request, orderNo string)

// 配置接口
func Configure(partner, key, seller string) error {
	if partner == "" || key == "" || seller == "" {
		return errors.New("miss some profile.")
	}
	conf_partner = partner
	conf_key = key
	conf_seller = seller
	return nil
}

// 按照支付宝规则生成sign
func alipaySign(param interface{}) string {
	//解析为字节数组
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return ""
	}
	//重组字符串
	var sign string
	oldString := string(paramBytes)
	//为保证签名前特殊字符串没有被转码，这里解码一次
	oldString = strings.Replace(oldString, `\u003c`, "<", -1)
	oldString = strings.Replace(oldString, `\u003e`, ">", -1)
	//去除特殊标点
	oldString = strings.Replace(oldString, "\"", "", -1)
	oldString = strings.Replace(oldString, "{", "", -1)
	oldString = strings.Replace(oldString, "}", "", -1)
	paramArray := strings.Split(oldString, ",")
	for _, v := range paramArray {
		detail := strings.SplitN(v, ":", 2)
		//排除sign和sign_type
		if detail[0] != "sign" && detail[0] != "sign_type" {
			//total_fee转化为2位小数
			if detail[0] == "total_fee" {
				number, _ := strconv.ParseFloat(detail[1], 32)
				detail[1] = strconv.FormatFloat(number, byte('f'), 2, 64)
			}
			if sign == "" {
				sign = detail[0] + "=" + detail[1]
			} else {
				sign += "&" + detail[0] + "=" + detail[1]
			}
		}
	}
	//追加密钥
	sign += conf_key
	//md5加密
	m := md5.New()
	m.Write([]byte(sign))
	sign = hex.EncodeToString(m.Sum(nil))
	return sign
}

type alipayParameters struct {
	InputCharset string  `json:"_input_charset"` //网站编码
	Body         string  `json:"body"`           //订单描述
	NotifyUrl    string  `json:"notify_url"`     //异步通知页面
	OutTradeNo   string  `json:"out_trade_no"`   //订单唯一id
	Partner      string  `json:"partner"`        //合作者身份ID
	PaymentType  uint8   `json:"payment_type"`   //支付类型 1：商品购买
	ReturnUrl    string  `json:"return_url"`     //回调url
	SellerEmail  string  `json:"seller_email"`   //卖家支付宝邮箱
	Service      string  `json:"service"`        //接口名称
	Subject      string  `json:"subject"`        //商品名称
	TotalFee     float32 `json:"total_fee"`      //总价
	Sign         string  `json:"sign"`           //签名，生成签名时忽略
	SignType     string  `json:"sign_type"`      //签名类型，生成签名时忽略
}

// 生成支付宝即时到帐的表单参数
// @params string 订单唯一id
// @params int 价格
// @params int 获得代金券的数量
// @params string 充值账户的名称
// @params string 充值描述
func CreatePaymentGateWay(orderNo string, fee float32, subject string, body string, returnUrl, notifyUrl string) string {
	//实例化参数
	param := &alipayParameters{}
	param.InputCharset = "utf-8"
	param.Body = body
	param.OutTradeNo = orderNo
	param.Partner = conf_partner
	param.PaymentType = 1
	param.ReturnUrl = returnUrl
	param.NotifyUrl = notifyUrl
	param.SellerEmail = conf_seller
	param.Service = "create_direct_pay_by_user"
	param.Subject = subject
	param.TotalFee = fee
	//生成签名
	sign := alipaySign(param)
	//追加参数
	param.Sign = sign
	param.SignType = "MD5"
	//生成自动提交form
	return `
		<form id="alipaysubmit" name="alipaysubmit" action="https://mapi.alipay.com/gateway.do?_input_charset=utf-8" method="get" style='display:none;'>
			<input type="hidden" name="_input_charset" value="` + param.InputCharset + `">
			<input type="hidden" name="body" value="` + param.Body + `">
			<input type="hidden" name="notify_url" value="` + param.NotifyUrl + `">
			<input type="hidden" name="out_trade_no" value="` + param.OutTradeNo + `">
			<input type="hidden" name="partner" value="` + param.Partner + `">
			<input type="hidden" name="payment_type" value="` + strconv.Itoa(int(param.PaymentType)) + `">
			<input type="hidden" name="return_url" value="` + param.ReturnUrl + `">
			<input type="hidden" name="seller_email" value="` + param.SellerEmail + `">
			<input type="hidden" name="service" value="` + param.Service + `">
			<input type="hidden" name="subject" value="` + param.Subject + `">
			<input type="hidden" name="total_fee" value="` + strconv.FormatFloat(float64(param.TotalFee), byte('f'), 2, 32) + `">
			<input type="hidden" name="sign" value="` + param.Sign + `">
			<input type="hidden" name="sign_type" value="` + param.SignType + `">
		</form>
		<script>
			document.forms['alipaysubmit'].submit();
		</script>
	`
}

func CreateTradeByUser() {
	/*
			 ////////////////////////////////////////////请求参数////////////////////////////////////////////

		        //支付类型
		        string payment_type = "1";
		        //必填，不能修改
		        //服务器异步通知页面路径
		        string notify_url = "http://商户网关地址/create_partner_trade_by_buyer-CSHARP-UTF-8/notify_url.aspx";
		        //需http://格式的完整路径，不能加?id=123这类自定义参数
		        //页面跳转同步通知页面路径
		        string return_url = "http://商户网关地址/create_partner_trade_by_buyer-CSHARP-UTF-8/return_url.aspx";
		        //需http://格式的完整路径，不能加?id=123这类自定义参数，不能写成http://localhost/
		        //卖家支付宝帐户
		        string seller_email = WIDseller_email.Text.Trim();
		        //必填
		        //商户订单号
		        string out_trade_no = WIDout_trade_no.Text.Trim();
		        //商户网站订单系统中唯一订单号，必填
		        //订单名称
		        string subject = WIDsubject.Text.Trim();
		        //必填
		        //付款金额
		        string price = WIDprice.Text.Trim();
		        //必填
		        //商品数量
		        string quantity = "1";
		        //必填，建议默认为1，不改变值，把一次交易看成是一次下订单而非购买一件商品
		        //物流费用
		        string logistics_fee = "0.00";
		        //必填，即运费
		        //物流类型
		        string logistics_type = "EXPRESS";
		        //必填，三个值可选：EXPRESS（快递）、POST（平邮）、EMS（EMS）
		        //物流支付方式
		        string logistics_payment = "SELLER_PAY";
		        //必填，两个值可选：SELLER_PAY（卖家承担运费）、BUYER_PAY（买家承担运费）
		        //订单描述
		        string body = WIDbody.Text.Trim();
		        //商品展示地址
		        string show_url = WIDshow_url.Text.Trim();
		        //需以http://开头的完整路径，如：http://www.商户网站.com/myorder.html
		        //收货人姓名
		        string receive_name = WIDreceive_name.Text.Trim();
		        //如：张三
		        //收货人地址
		        string receive_address = WIDreceive_address.Text.Trim();
		        //如：XX省XXX市XXX区XXX路XXX小区XXX栋XXX单元XXX号
		        //收货人邮编
		        string receive_zip = WIDreceive_zip.Text.Trim();
		        //如：123456
		        //收货人电话号码
		        string receive_phone = WIDreceive_phone.Text.Trim();
		        //如：0571-88158090
		        //收货人手机号码
		        string receive_mobile = WIDreceive_mobile.Text.Trim();
		        //如：13312341234


		        ////////////////////////////////////////////////////////////////////////////////////////////////

		        //把请求参数打包成数组
		        SortedDictionary<string, string> sParaTemp = new SortedDictionary<string, string>();
		        sParaTemp.Add("partner", Config.Partner);
		        sParaTemp.Add("_input_charset", Config.Input_charset.ToLower());
		        sParaTemp.Add("service", "create_partner_trade_by_buyer");
		        sParaTemp.Add("payment_type", payment_type);
		        sParaTemp.Add("notify_url", notify_url);
		        sParaTemp.Add("return_url", return_url);
		        sParaTemp.Add("seller_email", seller_email);
		        sParaTemp.Add("out_trade_no", out_trade_no);
		        sParaTemp.Add("subject", subject);
		        sParaTemp.Add("price", price);
		        sParaTemp.Add("quantity", quantity);
		        sParaTemp.Add("logistics_fee", logistics_fee);
		        sParaTemp.Add("logistics_type", logistics_type);
		        sParaTemp.Add("logistics_payment", logistics_payment);
		        sParaTemp.Add("body", body);
		        sParaTemp.Add("show_url", show_url);
		        sParaTemp.Add("receive_name", receive_name);
		        sParaTemp.Add("receive_address", receive_address);
		        sParaTemp.Add("receive_zip", receive_zip);
		        sParaTemp.Add("receive_phone", receive_phone);
		        sParaTemp.Add("receive_mobile", receive_mobile);

		        //建立请求
		        string sHtmlText = Submit.BuildRequest(sParaTemp, "get", "确认");
		        Response.Write(sHtmlText);

	*/
}

// 被动接收支付宝同步跳转的页面
func ReturnFunc(r *http.Request, sf SuccessNotifyFunc) error {
	//	type Params struct {
	//		Body        string `form:"body" json:"body"`                 //描述
	//		BuyerEmail  string `form:"buyer_email" json:"buyer_email"`   //买家账号
	//		BuyerId     string `form:"buyer_id" json:"buyer_id"`         //买家ID
	//		Exterface   string `form:"exterface" json:"exterface"`       //接口名称
	//		IsSuccess   string `form:"is_success" json:"is_success"`     //交易是否成功
	//		NotifyId    string `form:"notify_id" json:"notify_id"`       //通知校验id
	//		NotifyTime  string `form:"notify_time" json:"notify_time"`   //校验时间
	//		NotifyType  string `form:"notify_type" json:"notify_type"`   //校验类型
	//		OutTradeNo  string `form:"out_trade_no" json:"out_trade_no"` //在网站中唯一id
	//		PaymentType uint8  `form:"payment_type" json:"payment_type"` //支付类型
	//		SellerEmail string `form:"seller_email" json:"seller_email"` //卖家账号
	//		SellerId    string `form:"seller_id" json:"seller_id"`       //卖家id
	//		Subject     string `form:"subject" json:"subject"`           //商品名称
	//		TotalFee    string `form:"total_fee" json:"total_fee"`       //总价
	//		TradeNo     string `form:"trade_no" json:"trade_no"`         //支付宝交易号
	//		TradeStatus string `form:"trade_status" json:"trade_status"` //交易状态 TRADE_FINISHED或TRADE_SUCCESS表示交易成功
	//		Sign        string `form:"sign" json:"sign"`                 //签名
	//		SignType    string `form:"sign_type" json:"sign_type"`       //签名类型
	//	}

	vals := r.URL.Query()

	//生成签名
	sign := alipaySign(vals)

	// 校验签名
	if sign == vals.Get("sign") {
		tradeStat := vals.Get("trade_status")
		//判断订单是否已完成
		if tradeStat == "TRADE_FINISHED" || tradeStat == "TRADE_SUCCESS" {
			if sf != nil {
				sf(r, vals.Get("out_trade_no"))
			}
			return nil
		}
		return errors.New(tradeStat) //交易未完成
	}
	return errors.New("incorret sign") //签名失败
}

// 被动接收支付宝异步通知
func NotifyFunc(r *http.Request, sf SuccessNotifyFunc) string {
	r.ParseForm()
	vals := r.Form
	//生成签名
	sign := alipaySign(vals)

	// 校验签名
	if sign == vals.Get("sign") {
		tradeStat := vals.Get("trade_status")
		//判断订单是否已完成
		if tradeStat == "TRADE_FINISHED" || tradeStat == "TRADE_SUCCESS" {
			if sf != nil {
				sf(r, vals.Get("out_trade_no"))
				return "success"
			} else {
				return "inqure function to procesing."
			}
		}
		return tradeStat //交易未完成
	}
	return "incorret sign" //签名失败
}
