package winterSocket

import (
	"errors"
	"fmt"
	"reflect"
)

var route map[string]reflect.Value

func init() {
	route = make(map[string]reflect.Value)
}

func PutRoute(method string, fun any) error {
	if _, isOk := route[method]; isOk {
		return errors.New(fmt.Sprintf("route %s repetition", method))
	}
	route[method] = reflect.ValueOf(fun)
	return nil
}

func GetRoute(method string) (reflect.Value, bool) {
	v, is := route[method]
	return v, is
}
