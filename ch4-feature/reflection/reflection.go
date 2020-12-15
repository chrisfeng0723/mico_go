/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2020/12/14 19:32
 */

package main

import (
	"fmt"
	"reflect"
)

type Person interface {
	sayHello(name string)
	Run() string
}

type Hero struct {
	Name  string
	Age   int
	Speed int
}

func (hero *Hero) SayHello(name string) {
	fmt.Println("Hello"+name, ",I am "+hero.Name)
}

func (hero *Hero) Run() string {
	fmt.Println("I am running at speed ", hero.Speed)
	return "Running"
}

func main(){
	//typeOfHero := reflect.TypeOf(Hero{})
	//fmt.Printf("Hero's type is %s,kind is %s\n",typeOfHero,typeOfHero.Kind())
	//
	//typeOfPtrHero :=reflect.TypeOf(&Hero{})
	//typeElem := typeOfPtrHero.Elem()
	//fmt.Println(typeElem)
	//
	//for i :=0;i<typeOfHero.NumField();i++{
	//	fmt.Printf("field name is %s,type is %s,kind is %s\n",
	//		typeOfHero.Field(i).Name,
	//		typeOfHero.Field(i).Type,
	//		typeOfHero.Field(i).Type.Kind(),
	//		)
	//}
	//
	//nameField,_ := typeOfHero.FieldByName("Name")
	//fmt.Printf("field' name is %s, type is %s, kind is %s\n", nameField.Name, nameField.Type, nameField.Type.Kind())
	name :="小明"
	valueOfName := reflect.ValueOf(&name)
	fmt.Println(valueOfName.Elem().CanAddr())
}
