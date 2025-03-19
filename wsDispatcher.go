package winterSocket

import (
	"encoding/json"
	"github.com/gobwas/ws/wsutil"
	"net"
	"reflect"
)

func (that *WsServer) wsJSONDispatcher(cmd *Cmd, conn *WsConn, jsonDataByte []byte) (bool, []byte) {
	funVal, isOk := GetRoute(cmd.Cmd)

	if !isOk {
		if nil != that.tracking {
			that.tracking.Dispatcher404(conn, cmd, jsonDataByte)
		}
		return false, nil
	}

	refMtdType := funVal.Type()
	numIn := refMtdType.NumIn()
	methodParams := make([]reflect.Value, numIn)
	for i := 0; i < refMtdType.NumIn(); i++ {
		inType := refMtdType.In(i)
		switch inType.String() {
		case "*net.Conn":
			methodParams[i] = reflect.ValueOf(conn.Conn)
		case "*winterSocket.WsConn":
			methodParams[i] = reflect.ValueOf(conn)
		case "*winterSocket.Cmd":
			methodParams[i] = reflect.ValueOf(cmd)
		default:
			obj := reflect.New(inType)
			if err := json.Unmarshal(jsonDataByte, obj.Interface()); nil != err {
				pError("Dispatcher to json Unmarshal data failure [obj]", err)
				if nil != that.tracking {
					that.tracking.ParameterUnmarshalError(conn, cmd, jsonDataByte)
				} else {
					pError("ParameterUnmarshalError", nil)
				}
				break
			}
			methodParams[i] = obj.Elem()
			mp := make(map[string]interface{})
			if err := json.Unmarshal(jsonDataByte, &mp); nil != err {
				pError("Dispatcher to json Unmarshal data failure [mp]", err)
				if nil != that.tracking {
					that.tracking.ParameterUnmarshalError(conn, cmd, jsonDataByte)
				} else {
					pError("ParameterUnmarshalError", nil)
				}
				return false, nil
			}
			msg, isOk := requiredParamsReflect(obj, inType, mp)
			if !isOk {
				if nil != that.tracking {
					that.tracking.ParameterError(conn, msg)
				} else {
					pError("ParameterError", nil)
				}
				return false, nil
			}
		}
	}

	var res []byte
	result := funVal.Call(methodParams)
	if nil != result && 0 < len(result) {
		rsu := result[0]
		rsu = reflect.Indirect(rsu)
		if !rsu.IsValid() {
			return true, nil
		}
		switch rsu.Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice:
			res, _ = json.Marshal(rsu.Interface())
			_ = that.WriteText(res, conn.Conn)
		default:
			res = []byte(rsu.String())
			_ = that.WriteText(res, conn.Conn)
		}
	}

	return true, res
}

func (that *WsServer) WriteText(buff []byte, conn *net.Conn) error {
	if e := wsutil.WriteServerText(*conn, buff); nil != e {
		pError("winterSocket WriteBuff", e)
		return e
	}
	return nil
}
