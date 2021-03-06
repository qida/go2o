/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"errors"
	"fmt"
	"go2o/core/domain/interface/partner"
	"go2o/core/infrastructure"
	"go2o/share/variable"
	"time"
)

var _ partner.IPartner = new(Partner)

type Partner struct {
	value    *partner.ValuePartner
	saleConf *partner.SaleConf
	siteConf *partner.SiteConf
	rep      partner.IPartnerRep
	shops    []partner.IShop
	host     string
}

func NewPartner(v *partner.ValuePartner, rep partner.IPartnerRep) (partner.IPartner, error) {

	var err error

	if v == nil {
		fmt.Println("-----------", v)
		err = errors.New("101:no such partner")
		return nil, err
	}
	if time.Now().Unix() > v.ExpiresTime {
		err = errors.New("103: partner is expires")
	}

	return &Partner{
		value: v,
		rep:   rep,
	}, err
}

func (this *Partner) clearShopCache() {
	this.shops = nil
}

func (this *Partner) GetAggregateRootId() int {
	return this.value.Id
}
func (this *Partner) GetValue() partner.ValuePartner {
	return *this.value
}

func (this *Partner) SetValue(v *partner.ValuePartner) error {
	tv := this.value
	if v.Id == tv.Id {
		tv.Name = v.Name
		tv.Address = v.Address
		if v.LastLoginTime > 0 {
			tv.LastLoginTime = v.LastLoginTime
		}
		if v.LoginTime > 0 {
			tv.LoginTime = v.LoginTime
		}

		if len(v.Logo) != 0 {
			tv.Logo = v.Logo
		}
		tv.Phone = v.Phone
		tv.Pwd = v.Pwd
		tv.Secret = v.Secret
		tv.UpdateTime = time.Now().Unix()

		fmt.Printf("------------%#v\n", this.value)
	}
	return nil
}

// 保存
func (this *Partner) Save() (int, error) {
	var id int = this.GetAggregateRootId()
	var err error

	if id > 0 {
		return this.rep.SavePartner(this.value)
	}
	id, err = this.rep.SavePartner(this.value)

	this.rep.InitPartner(this.value.Id)

	return id, err
}

// 获取商户的域名
func (this *Partner) GetMajorHost() string {
	if len(this.host) == 0 {
		host := this.rep.GetPartnerMajorHost(this.GetAggregateRootId())
		if len(host) == 0 {
			host = fmt.Sprintf("%s.%s", this.value.Usr, infrastructure.GetContext().
				Config().GetString(variable.ServerDomain))
		}

		this.host = host
	}
	return this.host
}

// 获取销售配置
func (this *Partner) GetSaleConf() partner.SaleConf {
	if this.saleConf == nil {

		//10%分成
		//0.2,         #上级
		//0.1,         #上上级
		//0.8          #消费者自己

		this.saleConf = this.rep.GetSaleConf(
			this.GetAggregateRootId())
	}
	return *this.saleConf
}

// 保存销售配置
func (this *Partner) SaveSaleConf(v *partner.SaleConf) error {
	this.saleConf = v
	this.saleConf.PartnerId = this.value.Id
	return this.rep.SaveSaleConf(this.saleConf)
}

// 获取站点配置
func (this *Partner) GetSiteConf() partner.SiteConf {
	if this.siteConf == nil {
		this.siteConf = this.rep.GetSiteConf(this.GetAggregateRootId())
	}
	return *this.siteConf
}

// 保存站点配置
func (this *Partner) SaveSiteConf(v *partner.SiteConf) error {
	this.siteConf = v
	this.saleConf.PartnerId = this.value.Id
	return this.rep.SaveSiteConf(this.siteConf)
}

// 新建商店
func (this *Partner) CreateShop(v *partner.ValueShop) partner.IShop {
	v.CreateTime = time.Now().Unix()
	v.PartnerId = this.GetAggregateRootId()
	return newShop(this, v, this.rep)
}

// 获取所有商店
func (this *Partner) GetShops() []partner.IShop {
	if this.shops == nil {
		shops := this.rep.GetShopsOfPartner(this.GetAggregateRootId())
		this.shops = make([]partner.IShop, len(shops))
		for i, v := range shops {
			this.shops[i] = this.CreateShop(v)
		}
	}

	return this.shops
}

// 获取商店
func (this *Partner) GetShop(shopId int) partner.IShop {
	//	v := this.rep.GetValueShop(this.GetAggregateRootId(), shopId)
	//	if v == nil {
	//		return nil
	//	}
	//	return this.CreateShop(v)
	shops := this.GetShops()

	for _, v := range shops {
		if v.GetValue().Id == shopId {
			return v
		}
	}
	return nil
}

// 删除门店
func (this *Partner) DeleteShop(shopId int) error {
	//todo : 检测订单数量
	return this.rep.DeleteShop(this.GetAggregateRootId(), shopId)
}
