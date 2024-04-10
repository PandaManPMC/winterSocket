package winterSocket

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"reflect"
	"strconv"
	"strings"
)

func Dispatcher(method Method, conn *websocket.Conn, jsonDataStr string) bool {
	funVal, isOk := GetRoute(method.Method)

	if !isOk {
		if nil != tracking {
			tracking.Dispatcher404(conn)
		}
		return false
	}
	jsonBuff := []byte(jsonDataStr)

	refMtdType := funVal.Type()
	numIn := refMtdType.NumIn()
	methodParams := make([]reflect.Value, numIn)
	for i := 0; i < refMtdType.NumIn(); i++ {
		inType := refMtdType.In(i)
		switch inType.String() {
		case "*websocket.Conn":
			methodParams[i] = reflect.ValueOf(conn)
		default:
			obj := reflect.New(inType)
			if err := json.Unmarshal(jsonBuff, obj.Interface()); nil != err {
				pError("Dispatcher to json Unmarshal data failure [obj]", err)
				if nil != tracking {
					tracking.ParameterUnmarshalError(conn)
				} else {
					pError("ParameterUnmarshalError", nil)
				}
				break
			}
			methodParams[i] = obj.Elem()
			mp := make(map[string]interface{})
			if err := json.Unmarshal(jsonBuff, &mp); nil != err {
				pError("Dispatcher to json Unmarshal data failure [mp]", err)
				if nil != tracking {
					tracking.ParameterUnmarshalError(conn)
				} else {
					pError("ParameterUnmarshalError", nil)
				}
				return false
			}
			msg, isOk := requiredParamsReflect(obj, inType, mp)
			if !isOk {
				if nil != tracking {
					tracking.ParameterError(conn, msg)
				} else {
					pError("ParameterError", nil)
				}
				return false
			}
		}
	}

	result := funVal.Call(methodParams)
	if nil != result && 0 < len(result) {
		rsu := result[0]
		rsu = reflect.Indirect(rsu)
		if !rsu.IsValid() {
			return true
		}
		switch rsu.Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice:
			marshalData, _ := json.Marshal(rsu.Interface())
			WriteBuff(marshalData, conn)
		default:
			WriteBuff([]byte(rsu.String()), conn)
		}
	}

	return true
}

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

		json := typ.Tag.Get("json")
		if "" == json {
			continue
		}
		val, isOk := params[json]
		//if isOk {
		//	field.Set(reflect.ValueOf(val))
		//}

		required := typ.Tag.Get("required")
		if "" == required {
			continue
		}

		// 必传参数，字符串参数核验长度，其它类型只看是否有值
		if !isOk {
			return fmt.Sprintf("%s is a required parameter", json), false
		}

		// 核实字符串长度
		if "string" != field.Type().String() {
			continue
		}

		max := typ.Tag.Get("max")
		if "" != max {
			// 最长
			maxLen, _ := strconv.Atoi(max)
			s := strings.TrimSpace(val.(string))
			r := []rune(s)
			if maxLen < len(r) {
				return fmt.Sprintf("%s up to %s, your length is %d", json, max, len(r)), false
			}
		}

		min := typ.Tag.Get("min")
		if "" != min {
			// 最短
			minLen, _ := strconv.Atoi(min)
			s := strings.TrimSpace(val.(string))
			r := []rune(s)
			if minLen > len(r) {
				return fmt.Sprintf("%s minimum %s, your length is %d", json, max, len(r)), false
			}
		}
	}
	return "", true
}
