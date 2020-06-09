package utils

import (
	"fmt"
	"reflect"
)

func printValue(prefix string, v reflect.Value, visited map[interface{}]bool) {

	fmt.Printf("%s: ", v.Type())

	// Drill down through pointers and interfaces to get a value we can print.
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.Kind() == reflect.Ptr {
			// Check for recursive data
			if visited[v.Interface()] {
				fmt.Println("visted")
				return
			}
			visited[v.Interface()] = true
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		fmt.Printf("%d elements\n", v.Len())
		for i := 0; i < v.Len(); i++ {
			fmt.Printf("%s%d: ", prefix, i)
			printValue(prefix+"   ", v.Index(i), visited)
		}
	case reflect.Struct:
		t := v.Type()
		fmt.Printf("%d fields\n", t.NumField())
		for i := 0; i < t.NumField(); i++ {
			fmt.Printf("%s%s: ", prefix, t.Field(i).Name)
			printValue(prefix+"   ", v.Field(i), visited)
		}
	case reflect.Invalid:
		fmt.Printf("nil\n")
	default:
		fmt.Printf("%v\n", v.Interface())
	}
}
