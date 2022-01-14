package json_a

import (
	"encoding/json"
	"fmt"
	"strconv"
)

/*
Result of following function can be parsed in following way:

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}

 */
func RawByteArrToJson(data []byte) map[string]interface{}{

	var f interface{}
	_ = json.Unmarshal( data, &f)
	res := f.(map[string]interface{})
	return res
}



func Get_(json interface{}, fields ...string) interface{} {
	res := json
	for _, v := range fields {
		res = res.(map[string]interface{})[v]
	}
	return res
}

func DecodePath(path string) [][]string {
	var fields [][]string
	//dotIdx := 0
	curSubstrPtr := 0
	leftBracketIdx := 0
	rightBracketIdx := 0
	for i, rune := range path {
		if rune == '.' {
			//fields = append(fields, path[curSubstrPtr:i])
			curArr := []string{"map", path[curSubstrPtr:i]}
			fields = append(fields, curArr)
			curSubstrPtr = i + 1
		} else if rune == '['{
			leftBracketIdx = i
			if curSubstrPtr != leftBracketIdx {
				curArr := []string{"map", path[curSubstrPtr:leftBracketIdx]}
				fields = append(fields, curArr)
			}
		} else if rune == ']'{
			rightBracketIdx = i
			//fields = append(fields, path[curSubstrPtr:leftBracketIdx])
			//fields = append(fields, path[leftBracketIdx+1:rightBracketIdx])
			curArr := []string{"arr", path[leftBracketIdx+1:rightBracketIdx]}
			fields = append(fields, curArr)
			curSubstrPtr = i + 1
		}
	}
	if curSubstrPtr != len(path) {
		curArr := []string{"map", path[curSubstrPtr:]}
		fields = append(fields, curArr)
	}
//	fmt.Println("fields: ", fields)
	return fields
}

func Get(json interface{}, path string) interface{} {
	res := json
	fields := DecodePath(path)
	for _, v := range fields {
		if v[0] == "map" {
			res = res.(map[string]interface{})[v[1]]
		} else {
			i, _ := strconv.Atoi(v[1])
			tmp := res.([]interface{})
			res = tmp[i]
		}
	}
//	fmt.Println(res)
	return res
}

func GetArr(json interface{}, path string) []interface{} {
	res := json
	fields := DecodePath(path)
	for _, v := range fields {
		if v[0] == "map" {
			res = res.(map[string]interface{})[v[1]]
		} else {
			i, _ := strconv.Atoi(v[1])
			tmp := res.([]interface{})
			res = tmp[i]
		}
	}
//	fmt.Println(res)
	return res.([]interface{})
}

func GetMap(json interface{}, path string) map[string]interface{} {
	res := json
	fields := DecodePath(path)
	for _, v := range fields {
		if v[0] == "map" {
			res = res.(map[string]interface{})[v[1]]
		} else {
			i, _ := strconv.Atoi(v[1])
			tmp := res.([]interface{})
			res = tmp[i]
		}
	}
//	fmt.Println(res)
	return res.(map[string]interface{})
}


//https://stackoverflow.com/questions/52466447/loop-through-json-keys-and-values-and-same-time-replacing-specify-matched-value
func VerifyJSON(bv interface{}) {
	var dumpJSON func(v interface{}, kn string)
	dumpJSON = func( v interface{}, kn string) {
		iterMap := func(suff string, x map[string]interface{}, root string) {
			suff += "\t"
			var knf string
			if root == "root" {
				//knf = "%v/%v"
				knf = "%v %v %v"
			} else {
				//knf = "%v/%v"
				knf = "%v %v %v"
			}
			for k, v := range x {
				switch vv := v.(type) {
				case map[string]interface{}:
					//fmt.Printf("%s => (map[string]interface{}) ...\n", fmt.Sprintf(knf, root, k))
					//fmt.Printf("%s => \n", fmt.Sprintf(knf, root, k) )
					fmt.Printf("%s => \n", fmt.Sprintf(knf, suff, root, k) )
					//fmt.Println(suff, k, " => " )
				case []interface{}:
					//fmt.Printf("%s => ([]interface{}) ...\n", fmt.Sprintf(knf, root, k))
					//fmt.Printf("%s => \n", fmt.Sprintf(knf, root, k))
					fmt.Printf("%s => \n", fmt.Sprintf(knf, suff, root, k))
					//fmt.Println(suff, k, " => " )
				default:
					//fmt.Printf("%s => %v\n", fmt.Sprintf(knf, root, k), vv)
					fmt.Printf("%s => %v\n", fmt.Sprintf(knf, suff, root, k), vv)
					//fmt.Println(suff, k, " => ", vv)
					x[k] = "rgk"
				}
				dumpJSON(v, fmt.Sprintf(knf, suff, root, ""))
			}
		}
		iterSlice := func(suff string, x []interface{}, root string) {
			suff += "\t"
			var knf string
			if root == "root" {
				//knf = "%v/%v"
				knf = "%v %v %v"
			} else {
				//knf = "%v/%v"
				knf = "%v %v %v"
			}
			for k, v := range x {
				switch vv := v.(type) {
				case map[string]interface{}:
					//fmt.Printf("%s => (map[string]interface{}) ...\n", fmt.Sprintf(knf, root, k))
					//fmt.Printf("%s => \n", fmt.Sprintf(knf, root, k) )
					fmt.Printf("%s => \n", fmt.Sprintf(knf, suff, root, k) )
					//fmt.Println(suff, k, " => " )
				case []interface{}:
					//fmt.Printf("%s => ([]interface{}) ...\n", fmt.Sprintf(knf, root, k))
					//fmt.Printf("%s => \n", fmt.Sprintf(knf, root, k))
					fmt.Printf("%s => \n", fmt.Sprintf(knf, suff, root, k))
					//fmt.Println(suff, k, " => " )
				default:
					//fmt.Printf("%s => %v\n", fmt.Sprintf(knf, root, k), vv)
					fmt.Printf("%s => %v\n", fmt.Sprintf(knf, suff, root, k), vv)
					//fmt.Println(suff, k, " => ", vv)
					x[k] = "rg"
				}
				dumpJSON(v, fmt.Sprintf(knf, suff, root, ""))
			}
		}

		switch vv := v.(type) {
		case map[string]interface{}:
			iterMap("", vv, kn)
		case []interface{}:
			iterSlice("", vv, kn)
		default:
		}
	}
	dumpJSON( bv, " ")
}
