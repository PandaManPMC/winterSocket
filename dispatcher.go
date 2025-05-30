package winterSocket

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"reflect"
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
