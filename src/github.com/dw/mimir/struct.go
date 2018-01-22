package mimir

import (
	"net/url"
	"reflect"
)

const (
	structTagForm = "form"
)

// FormArgs for generating the http application/x-www-form-urlencoded payload
type FormArgs struct {
	AppID     string `form:"app_id"`
	StartDate string `form:"startdate"`
	EndDate   string `form:"enddate"`
	// IsIOS is for Android or iOS
	IsIOS   string `form:"isiOS"`
	Version string `form:"version"`
	Channel string `form:"channel"`
	//Type is for daily, weekly, monthly
	Type string `form:"type"`
}

func (fa *FormArgs) GenBody() []byte {
	params := url.Values{}
	v := reflect.ValueOf(fa).Elem()
	t := v.Type()
	b := v.NumField()
	for i := 0; i < b; i++ {
		field := t.Field(i)
		pk := field.Tag.Get(structTagForm)
		fv := v.Field(i)
		if fv.Kind() == reflect.String {
			pv := fv.String()
			if len(pv) > 0 {
				params.Set(pk, pv)
			}
		}
	}
	return []byte(params.Encode())
}
