package winterSocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"reflect"
	"testing"
)

type Dog struct {
	Name  string `json:"name" required:"true"`
	Age   int    `json:"age" required:"true"`
	Hobby string `json:"hobby" required:"true"`
}

type Animal struct {
}

func (that *Animal) Eat(conn *websocket.Conn, dog Dog) {
	println("吃东西")
	fmt.Println(dog)
}

func TestRoute(t *testing.T) {
	dog := Dog{
		Name: "大黄",
		Age:  0,
	}

	ani := new(Animal)
	ani.Eat(nil, dog)

	t.Log("----")

	if e := PutRoute("eat", ani.Eat); nil != e {
		panic(e)
	}

	t.Log("----")

	eatVal, _ := GetRoute("eat")

	//eatVal := reflect.ValueOf(eat)
	//eatVal.Call(nil) // panic: reflect: Call with too few input arguments

	t.Log("----")
	dogVal := reflect.ValueOf(dog)
	vas := make([]reflect.Value, 2)
	vas[0] = reflect.ValueOf(new(websocket.Conn))
	vas[1] = dogVal
	eatVal.Call(vas)

	t.Log("---- 读取方法参数")

	buf, _ := json.Marshal(dog)

	// 读取方法参数
	refMtdType := eatVal.Type()
	numIn := refMtdType.NumIn()
	methodParams := make([]reflect.Value, numIn)
	for i := 0; i < refMtdType.NumIn(); i++ {
		inType := refMtdType.In(i)
		t.Log(inType.String())
		switch inType.String() {
		case "*websocket.Conn":
			c := new(websocket.Conn)
			methodParams[i] = reflect.ValueOf(c)
		default:
			obj := reflect.New(inType)
			if err := json.Unmarshal(buf, obj.Interface()); nil != err {
				t.Log("requestToData", err)
				t.Log(errors.New("request to json Unmarshal data failure"))
			}
			methodParams[i] = obj.Elem()
			t.Log("存一个 map 用于 required 校验")
			mp := make(map[string]interface{})
			json.Unmarshal(buf, &mp)
			delete(mp, "hobby")

			t.Log(mp)
			t.Log(obj.Elem())
		}
	}

	t.Log(len(methodParams))
	eatVal.Call(methodParams)
}
