package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
)

var globalValidate *validateCtx

type validateCtx struct {
	*validator.Validate
	trans ut.Translator
}

func ValidateStruct(param interface{}) (string, error) {
	err := globalValidate.Struct(param)
	if err != nil {
		invalidMsg := globalValidate.getError(err)
		return invalidMsg, err
	}
	return "", err
}

func init() {
	zhTranslator := zh.New()
	uni := ut.New(zhTranslator, zhTranslator)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	//通过自定义标签label来替换字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	_ = zh2.RegisterDefaultTranslations(validate, trans)
	globalValidate = &validateCtx{validate, trans}
}

func (vc *validateCtx) getError(errs error) string {
	var errStr string
	for _, err := range errs.(validator.ValidationErrors) {
		if vc.trans != nil {
			errStr = err.Translate(vc.trans)
		} else {
			errStr = err.Field() + "验证不符合" + err.Tag()
		}
		break
	}
	return errStr
}
