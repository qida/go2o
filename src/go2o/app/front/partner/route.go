/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/app/session"
	"net/http"
	"net/url"
)

var routes *web.RouteMap = new(web.RouteMap)

//处理请求
func Handle(w http.ResponseWriter, r *http.Request) {
	routes.Handle(w, r)
}

func chkLogin(r *http.Request) (b bool, partnerId int) {
	//todo:仅仅做了id的检测，没有判断有效性
	i, err := session.LSession.GetPartnerIdFromCookie(r)
	return err == nil, i
}

func redirect(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<script>window.parent.location.href='/login?return_url=" +
		url.QueryEscape(r.URL.String()) + "'</script>"))
}

//注册路由
func RegisterRoutes(context app.Context) {
	mc := &mainC{Context: context} //入口控制器
	sc := &shopC{context}          //商家门店控制器
	gc := &goodsC{context}         //食谱控制器
	lc := &loginC{Context: context}
	cc := &commC{Context: context}
	oc := &orderC{Context: context}
	cat_c := &categoryC{Context: context}
	conf_c := &configC{Context: context}
	prom_c := &promC{Context: context}

	routes.Add("^/comm/", func(w http.ResponseWriter, r *http.Request) {
		mvc.Handle(cc, w, r, true)
	})

	routes.Add("^/pt/shop/",
		func(w http.ResponseWriter, r *http.Request) {
			if b, id := chkLogin(r); b {
				mvc.Handle(sc, w, r, true, id)
			} else {
				redirect(w, r)
			}
		})

	routes.Add("^/pt/category/",
		func(w http.ResponseWriter, r *http.Request) {
			if b, id := chkLogin(r); b {
				mvc.Handle(cat_c, w, r, true, id)
			} else {
				redirect(w, r)
			}
		})

	routes.Add("^/pt/order/",
		func(w http.ResponseWriter, r *http.Request) {
			if b, id := chkLogin(r); b {
				mvc.Handle(oc, w, r, true, id)
			} else {
				redirect(w, r)
			}
		})

	routes.Add("^/pt/goods/",
		func(w http.ResponseWriter, r *http.Request) {
			if b, id := chkLogin(r); b {
				mvc.Handle(gc, w, r, true, id)
			} else {
				redirect(w, r)
			}
		})

	routes.Add("^/pt/prom/",
		func(w http.ResponseWriter, r *http.Request) {
			if b, id := chkLogin(r); b {
				mvc.Handle(prom_c, w, r, true, id)
			} else {
				redirect(w, r)
			}
		})

	routes.Add("^/export/getExportData$", func(w http.ResponseWriter, r *http.Request) {
		if b, id := chkLogin(r); b {
			GetExportData(w, r, id)
		} else {
			redirect(w, r)
		}
	})

	routes.Add("^/pt/conf/",
		func(w http.ResponseWriter, r *http.Request) {
			if b, id := chkLogin(r); b {
				mvc.Handle(conf_c, w, r, true, id)
			} else {
				redirect(w, r)
			}
		})

	routes.Add("^/login$", func(w http.ResponseWriter, r *http.Request) {
		mvc.Handle(lc, w, r, true)
	})

	routes.Add("^/", func(w http.ResponseWriter, r *http.Request) {
		if b, id := chkLogin(r); b {
			mvc.Handle(mc, w, r, true, id)
		} else {
			redirect(w, r)
		}
	})
}
