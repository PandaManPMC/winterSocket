package winterSocket

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// requiredParamsReflect 检查结构体的必传参数等
// 支持 required ，字符串类型支持 min、max 字符长度，会去除首位空格。
func requiredParamsReflect(data reflect.Value, typeof reflect.Type, params map[string]interface{}) (string, bool) {
	dr := data
	drE := reflect.Indirect(dr)

	dt := typeof
	dtE := typeof

	if !data.IsValid() {
		dtE = dt.Elem()
	}

	for i := 0; i < dtE.NumField(); i++ {
		field := drE.Field(i)
		typ := dtE.Field(i)

		json_ := typ.Tag.Get("json")
		if "" == json_ {
			continue
		}
		val, isOk := params[json_]
		//if isOk {
		//	field.Set(reflect.ValueOf(val))
		//}

		required := typ.Tag.Get("required")
		if "" == required {
			continue
		}

		// 必传参数，字符串参数核验长度，其它类型只看是否有值
		if !isOk {
			return fmt.Sprintf("%s is a required parameter", json_), false
		}

		// 核实字符串长度
		if "string" != field.Type().String() {
			continue
		}

		max_ := typ.Tag.Get("max")
		if "" != max_ {
			// 最长
			maxLen, _ := strconv.Atoi(max_)
			s := strings.TrimSpace(val.(string))
			r := []rune(s)
			if maxLen < len(r) {
				return fmt.Sprintf("%s up to %s, your length is %d", json_, max_, len(r)), false
			}
		}

		min_ := typ.Tag.Get("min")
		if "" != min_ {
			// 最短
			minLen, _ := strconv.Atoi(min_)
			s := strings.TrimSpace(val.(string))
			r := []rune(s)
			if minLen > len(r) {
				return fmt.Sprintf("%s minimum %s, your length is %d", json_, min_, len(r)), false
			}
		}
	}
	return "", true
}
