package protobuf

import (
	"fmt"
	"reflect"
	"regexp"
)

func printFields(prefix string, t reflect.Type) {
	if t.Kind() != reflect.Struct {
		return
	}
	//	re := regexp.MustCompile(",([0-9]+)")
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Println(prefix, f.Name, f.Type )
		if f.Type.Kind() == reflect.Struct {
			//			fmt.Println(reflect.New(f.Type))
			printFields(fmt.Sprintf("  %v", prefix), f.Type)
		} else if f.Type.Kind() == reflect.Ptr {
			printFields(fmt.Sprintf("  %v", prefix), f.Type.Elem())
		}
	}
}

func ReflectPrintStructProtobuf(s interface{}) {
	printFieldsProtobuf("", reflect.ValueOf(s).Type())
}

func printFieldsProtobuf(prefix string, t reflect.Type) {
	if t.Kind() != reflect.Struct {
		return
	}
	re := regexp.MustCompile(",([0-9]+)")
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		regRes := re.FindStringSubmatch(f.Tag.Get("protobuf"))
		if len(regRes) > 1 {
			fmt.Println(prefix, f.Name, f.Type, regRes[1] )
		} else {
			//			fmt.Println(prefix, f.Name, f.Type )
		}
		if f.Type.Kind() == reflect.Struct {
			//			fmt.Println(reflect.New(f.Type))
			printFields(fmt.Sprintf("  %v", prefix), f.Type)
		} else if f.Type.Kind() == reflect.Ptr {
			printFields(fmt.Sprintf("  %v", prefix), f.Type.Elem())
		}
	}
}

func reflectInitStuct(prefix string, rp reflect.Value) {
	//	fmt.Println(prefix, "rp:", rp)
	if rp.Kind() != reflect.Struct {
		return
	}
	for i:=0; i<rp.NumField(); i++ {
		if rp.Field(i).Kind() == reflect.Ptr {
			rp.Field(i).Set(reflect.New(rp.Field(i).Type().Elem()))
			fmt.Println( prefix, rp.Field(i).Type() )
			reflectInitStuct(fmt.Sprintf("	%v", prefix), reflect.Indirect(rp.Field(i)))
		} else {
			fmt.Println( prefix, rp.Field(i).Type() )
		}
	}
}

