package valid_helper

import (
	"gin_base/app/helper/exception_helper"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"sync"
)

var uni *ut.UniversalTranslator
var trans ut.Translator
var validate *validator.Validate
var once sync.Once

// 校验参数
// 参考使用：https://pkg.go.dev/github.com/go-playground/validator/v10
func Check(data interface{}) {
	once.Do(func() {
		uni = ut.New(zh.New())
		trans, _ = uni.GetTranslator("zh")

		validate = validator.New()
		//注册一个函数，获取struct tag里自定义的label作为字段名
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			return name
		})
		zh_translations.RegisterDefaultTranslations(validate, trans)
	})

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			exception_helper.CommonException(err.Translate(trans))
		}
	}
}
