package main

//*
import (
	"./src/api/digikey"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)
import "./src/storage/textFile"
//*/
/*
import "./src/storage/redis"
//*/

//import "./src/interface/mqtt"
//*
import "./src/encoding/base64"
import "./src/encoding/json_a"

//*/
import "./src/api/electronic_component"
//*/
import "./src/storage/postgresql"
//*/
import "./src/storage/excel"
import "./src/web"

import "./src/storage/redis"

var confMap map[string]string

func import_database() {

	bd := Excel.ReadSheet(confMap["excel_file_to_read"], 0)

	importFirstIdx := 2109
	importLastIdx := 2195
	for i := importFirstIdx; i < importLastIdx; i += 1 {
		curCell := bd[i]

		pn := curCell[2]
		cellIdx := curCell[0]
		// check if pn already exist
		componentInDb := Postgres.TableGetRow("part_numbers", "part_number, help_url", "WHERE part_number='"+pn + "'")
		helpUrl := componentInDb["help_url"]
		if len(componentInDb) == 0 {
			component, err := electronic_component.Get_info(pn)
			helpUrl = component["help_url"]
			if (err == nil) {
				fmt.Println("info: importing new Part Number", pn)

				if component["component_type"] == "Capacitors" {
					descrCap := component["description"].(string)
					descrCap += " "
					descrCap += component["tolerance"].(string)
					component["description"] = descrCap
				}
				Postgres.TableInsertRow("part_numbers", component)

			} else {
				fmt.Println("WARNING: unable to find ", pn,". CellIdx ", cellIdx, " is missed")
				continue
			}
		} else {
//			fmt.Println("INFO:", pn, " exist in database")
		}

		cellInDb := Postgres.TableGetRow("cells", "cellIdx", "WHERE cellIdx='"+cellIdx + "'")
		if len(cellInDb) != 0 {
//			fmt.Println("INFO:", cellIdx, " exist in database!")
			continue
		}

		// check if cell already exist
		cell := make(map[string]interface{})
		cell["cellIdx"] = cellIdx //curCell[0]
		cell["part_number"] = pn //curCell[2]
		cell["qty"] = curCell[5]

		if (len(curCell) > 6) && (curCell[6] != "") {
			cell["supplier_link"] = curCell[6]
		} else {
			cell["supplier_link"] = helpUrl
		}

		if len(curCell) > 7 {
			cell["notes"] = strings.Replace(curCell[7], "'", "", -1) //curCell[7]
		} else {
			cell["notes"] = ""
		}
		fmt.Println("info: insering new cell: ", cellIdx)
		Postgres.TableInsertRow("cells", cell)

		time.Sleep(450 * time.Millisecond)
	}

	Redis.Set("import_first_idx", importFirstIdx)
	Redis.Set("import_last_idx", importLastIdx)
}

func import_database_fix() {

	bd := Excel.ReadSheet(confMap["excel_file_to_read"], 0)

	importFirstIdx := 39
	importLastIdx := 2003
	for i := importFirstIdx; i < importLastIdx; i += 1 {
		curCell := bd[i]
		cur_cellIdx := curCell[0]
		cur_note := ""
		cur_note2 := ""
		cur_noteIsUrl := false
		cur_noteIsUrl2 := false
		res_note := ""
		res_url := ""
		if len(curCell) >= 7 {
			cur_note = curCell[6]
			cur_note = strings.Replace(cur_note, "'", "", -1)
			cur_noteIsUrl, _ = regexp.MatchString("^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;%=.]+$", cur_note)
			if !cur_noteIsUrl {
				if strings.Contains(cur_note, "http://www.chipmall.com/") {
					cur_noteIsUrl = true
				}
			}
		}
		if len(curCell) >= 8 {
			cur_note2 = curCell[7]
			cur_note2 = strings.Replace(cur_note2, "'", "", -1)
			cur_noteIsUrl2, _ = regexp.MatchString("^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;%=.]+$", cur_note2)
			if !cur_noteIsUrl2 {
				if strings.Contains(cur_note2, "http://www.chipmall.com/") {
					cur_noteIsUrl2 = true
				}
			}
		}

		if cur_noteIsUrl {
			res_url = cur_note
			res_note = cur_note2
		} else if cur_noteIsUrl2 {
			res_url = cur_note2
			res_note = cur_note
		} else {
			if cur_note != "" {
				res_note += cur_note
				if cur_note2 != "" {
					res_note += "; "
				}
			}
			if cur_note2 != "" {
				res_note += cur_note2
			}
		}
		Postgres.TableUpdateRow("cells", map[string]interface{}{ "supplier_link": res_url, "notes" : res_note}, "cellIdx = '" + cur_cellIdx + "'")

	}

	Redis.Set("import_first_idx", importFirstIdx)
	Redis.Set("import_last_idx", importLastIdx)
}


func main() {

/*
	Aliexpress.Init(confMap)
	wd, _ := selenium_cl.InitWebDriver()
	wd, err := Promelec.Autorization(wd)
	Promelec.GetOrdersList(wd, 30)
	fmt.Println(wd, err)
	/*
	wd, err := Aliexpress.Autorization(wd)
	if err != nil {
		fmt.Println("Aliexpress autorization failed")
	}
	resLinks, wd, err := Aliexpress.GetOrdersList(wd, 30)
	fmt.Println("resLinks: ", resLinks)
//*/
//	selenium_cl.ServiceStop()

//*/

	confMap = TextFile.Read_key_value_map("conf/conf.txt", " ")
	Digikey.Init(confMap)

	//*
	go Digikey.Run_redirect_URI_server()
	Digikey.Autorization_get_code()
	Digikey.Access_get_token()
	//*/

	//ress, _ := electronic_component.Get_info("GRM1555C1H2R2WA01D")
	electronic_component.Get_info("GRM1555C1H2R2WA01D")
	Postgres.Init(confMap)
	Postgres.Connect()

	//*
	err := Redis.Connect()
	if err != nil {
		log.Println("Redis connection error")
//		panic(err)
	} else {
		log.Println("Redis connected")
	}
	//*/

	electronic_component.Init()
	electronic_component.InitHolder()
	electronic_component.InitParsel()
	electronic_component.InitDevice()
	electronic_component.InitDeviceBom()
	//electronic_component.AlterDeviceBom()
	//electronic_component.AlterDeviceBom2()
	//electronic_component.AlterDeviceBom3()

/*
	digi_selenium.Init(confMap)
	wd, _ := selenium_cl.InitWebDriver()
	wd, err = digi_selenium.Autorization(wd)
//	digi_selenium.GetOrdersList(wd, 30)
	fmt.Println(wd, err)
	selenium_cl.ServiceStop()
//*/
	//import_database()
	//import_database_fix()

//	_/*res :*/= Postgres.TableGetRow("cells", "cellIdx, part_number, qty", "WHERE part_number='RC0603JR-070RL'")

//	_/*res3 :*/= Postgres.TableGetRows("part_numbers", "part_number, description, footprint, temperature, help_url", "WHERE part_number='RC0603JR-070RL'")

	res4 := Postgres.TableGetRows("cells INNER JOIN part_numbers", "cells.cellIdx, cells.part_number, part_numbers.description, part_numbers.ref_supplier, part_numbers.footprint, part_numbers.temperature, part_numbers.help_url, cells.qty, cells.supplier_link, cells.notes", "ON part_numbers.part_number=cells.part_number and cellidx='1A13'")
//	fmt.Println("before Marshalling: ", res4)
	warehouseData, _ := json.Marshal(res4)
//	fmt.Println("after Marshalling", string(warehouseData))
	Web.SetJsonData(warehouseData)

	go Web.WebServer()

	//*/

	//	Digikey.Get_order_history()
	//	Digikey.Get_order_status()
//*/

/*
	Mqtt.Connect(confMap["mqtt_server_addr"])
	Mqtt.Subscribe("gateway/+/stats")
	Mqtt.SubscribeClbk("gateway/aa555a0000000005/rx", mqtt_clbk)
	Mqtt.SubscribeClbk("gateway/aa555a0000000005/tx", mqtt_clbk_tx)
	// application/1/device/343934354d377316/rx
	Mqtt.SubscribeClbk("application/+/device/+/rx", mqtt_clbk_app)

//	data := []byte(`{"applicationID":"1","applicationName":"electric_meter","deviceName":"electric_meter_m303","devEUI":"343934354d377316","rxInfo":[{"gatewayID":"aa555a0000000005","name":"GW-01_dev","rssi":-99,"loRaSNR":-15.5,"location":{"latitude":0,"longitude":0,"altitude":0}}],"txInfo":{"frequency":868300000,"dr":0},"adr":true,"fCnt":1,"fPort":2,"data":"MDEyMzQ1Njc4OUFCQ0RFRg=="}`)
//	res := json.RawByteArrToJson(data)
//	fmt.Println(res["data"])

//	decoded := base64.Decode(res["data"].(string))
//	for _, n := range(decoded) {
//		fmt.Printf("%02X ", n)
//	}
	fmt.Println("")
//*/
	for {
		time.Sleep(1 * time.Second)
		res4 := Postgres.TableGetRows("cells INNER JOIN part_numbers", "cells.cellIdx, cells.part_number, part_numbers.description, part_numbers.ref_supplier, part_numbers.footprint, part_numbers.temperature, part_numbers.help_url, cells.qty, cells.supplier_link, cells.notes", "ON part_numbers.part_number=cells.part_number")
		warehouseData, _ := json.Marshal(res4)
		Web.SetJsonData(warehouseData)
	}
}

func mqtt_clbk(topic string, message []byte)  {
	fmt.Printf("Mqtt RX. Topic: %s. Message: %s\r\n", topic, message)
}

func mqtt_clbk_tx(topic string, message []byte)  {
	fmt.Printf("Mqtt TX. Topic: %s. Message: %s\r\n", topic, message)
}

func mqtt_clbk_app(topic string, message []byte)  {
	fmt.Printf("Mqtt App. Topic: %s. Message: %s\r\n", topic, message)

	res := json_a.RawByteArrToJson(message)
	decoded := base64.Decode(res["data"].(string))
	for _, n := range(decoded) {
		fmt.Printf("%02X ", n)
	}
	fmt.Println("")

}