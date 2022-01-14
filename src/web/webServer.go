package Web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

import "../storage/postgresql"

import "../api/electronic_component"

func getWarehouseData(w http.ResponseWriter, r *http.Request) {
	message := getJsonData()
	w.Write(message)
}

func getNewPartParams(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {

	params := r.URL.Query()
	component_data := map[string]interface{}{}

	for key, values := range params {
		for _, val := range values {
			component_data[key] = val
		}
	}

	//check if part number field exist
	if _, ok := component_data["part_number"]; !ok {
//		jsonString, _ := json.Marshal(res)
//		w.Write(jsonString)
		err := fmt.Errorf("empty Part number field")
		return component_data, err
	}

	//replacing redundant spaces in part number to reduce duplicated part numbers
	re_leadclose_whtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := re_leadclose_whtsp.ReplaceAllString(component_data["part_number"].(string), "")
	final = re_inside_whtsp.ReplaceAllString(final, " ")
	fmt.Println(final)
	component_data["part_number"] = final

	fmt.Println(component_data)

	// check if part number already exist in DB
	componentInDb := Postgres.TableGetRow("part_numbers", "part_number, help_url", "WHERE part_number='"+component_data["part_number"].(string) + "'")
	if len(componentInDb) != 0 {
//		jsonString, _ := json.Marshal(res)
//		w.Write(jsonString)
		err := fmt.Errorf("this part number already exist")
		return component_data, err
	}
	return component_data, nil
}

func apiCreatePartNumberFast(w http.ResponseWriter, r *http.Request) {
	fmt.Println("apiCreatePartNumber")
	res := map[string]interface{}{
		"success" : "true"}

	newPartNumberParams, err := getNewPartParams(w, r)
	if err != nil {
		res["success"] = false
		res["error"] = err.Error()
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	fmt.Println("part number to find at Digikey: ", newPartNumberParams["part_number"])
	err = electronic_component.PartNumber_add( newPartNumberParams["part_number"].(string) )
	if err != nil {
		res["success"] = false
		res["error"] = err.Error()
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)
}

func getQueryParams(r *http.Request) map[string]interface{} {
	params := r.URL.Query()
	res := map[string]interface{}{}

	for key, values := range params {
		for _, val := range values {
			res[key] = val
		}
	}
	return res
}

func apiCreatePartNumber(w http.ResponseWriter, r *http.Request) {
	fmt.Println("apiCreatePartNumber")
	res := map[string]interface{}{
		"success" : "true"}
	// To-Do: following to be replaced by getNewPartParams(w,r)
//*
	params := r.URL.Query()
	component_data := map[string]interface{}{}

	for key, values := range params {
		for _, val := range values {
			component_data[key] = val
		}
	}

	//check if part number field exist
	if _, ok := component_data["part_number"]; !ok {
		res["success"] = false
		res["error"] = "empty Part number field"
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	//replacing redundant spaces in part number to reduce duplicated part numbers
	re_leadclose_whtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := re_leadclose_whtsp.ReplaceAllString(component_data["part_number"].(string), "")
	final = re_inside_whtsp.ReplaceAllString(final, " ")
	fmt.Println(final)
	component_data["part_number"] = final

	fmt.Println(component_data)

	// check if part number already exist in DB
	componentInDb := Postgres.TableGetRow("part_numbers", "part_number, help_url", "WHERE part_number='"+component_data["part_number"].(string) + "'")
	if len(componentInDb) != 0 {
		res["success"] = false
		res["error"] = "this Part number already exist in DB"
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}
//*/
	err := Postgres.TableInsertRow("part_numbers", component_data)
	if err != nil {
		res["success"] = false
		res["error"] = err.Error()
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)
}

func apiGetParcels(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("apiGetParcels")
	res := map[string]interface{}{
		"success" : "true"}

	res4 := Postgres.TableGetRows("parcel", "order_id, index_in_parcel, store, tracking_number, part_number, supplier_link, status, qty ", "")
	parcelsData, _ := json.Marshal(res4)
//	fmt.Println("parcels: ", string(parcelsData))
	res["data"] = string(parcelsData)
	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)
}

func apiCreateCell(w http.ResponseWriter, r *http.Request) {
	fmt.Println("apiCreateCell")
	res := map[string]interface{}{
		"success" : "true"}

	params := r.URL.Query()
	cell_data := map[string]interface{}{}

	for key, values := range params {
		for _, val := range values {
			cell_data[key] = val
		}
	}

	//check if cell idx field exist
	if _, ok := cell_data["cellIdx"]; !ok {
		res["success"] = false
		res["error"] = "empty Cell idx field"
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	//replacing redundant spaces in cell idx to reduce duplicated part numbers
	re_leadclose_whtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := re_leadclose_whtsp.ReplaceAllString(cell_data["cellIdx"].(string), "")
	final = re_inside_whtsp.ReplaceAllString(final, " ")
	fmt.Println(final)
	cell_data["cellIdx"] = final

	fmt.Println(cell_data)

	// check if cell idx already exist in DB
	componentInDb := Postgres.TableGetRow("cells", "cellIdx", "WHERE part_number='"+cell_data["cellIdx"].(string) + "'")
	if len(componentInDb) != 0 {
		res["success"] = false
		res["error"] = "this cell idx already exist in DB"
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	err := Postgres.TableInsertRow("cells", cell_data)
	if err != nil {
		res["success"] = false
		res["error"] = err.Error()
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)
}

var lastUploadedDeviceName = ""

func deviceUploadName(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("deviceUploadName")
	res := map[string]interface{}{
		"success" : "true"}

	params := r.URL.Query()
	//component_data := map[string]interface{}{}

	for key, values := range params {
		if (key == "deviceName") && (len(values) > 0){
			//fmt.Println("deviceName: ", values[0])
			lastUploadedDeviceName = values[0]
		}
	}

	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)

}

func deviceList(w http.ResponseWriter, r *http.Request)  {
	res := map[string]interface{}{
		"success" : "true"}

	deviceListQ := electronic_component.Device_getAll()
	deviceList, _ := json.Marshal(deviceListQ)
	//	fmt.Println("parcels: ", string(parcelsData))
	res["data"] = string(deviceList)

	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)
}

func deviceBom(w http.ResponseWriter, r *http.Request) {
	res := map[string]interface{}{
		"success" : "true"}

	params := getQueryParams(r)

	requested_device_id := ""
	if _, ok := params["device_id"]; !ok {
		res["success"] = false
		res["error"] = "device_id not exist in query"
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	requested_device_id = params["device_id"].(string)

	//deviceBomQ := electronic_component.DeviceBom_getByIdx(requested_device_id)
	deviceBomQ := electronic_component.DeviceBomStock_getByIdx(requested_device_id)

	// SELECT devicebom.part_number, devicebom.qty, sum(cells.qty) as warehouseStock FROM devicebom LEFT JOIN cells ON devicebom.part_number=cells.part_number GROUP BY devicebom.device_id, devicebom.qty, devicebom.part_number;
	deviceBom, _ := json.Marshal(deviceBomQ)
	res["data"] = string(deviceBom)

	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)
}

func deviceUpload(w http.ResponseWriter, r *http.Request) {

	var failed_PN bytes.Buffer
	failed_PN.WriteString("")
	res := map[string]interface{}{
		"success" : "true",
		"error" : ""}

	var deviceName = lastUploadedDeviceName
	lastUploadedDeviceName = ""
	//fmt.Println("deviceUpload ", deviceName)
	var Buf bytes.Buffer
	//file, header, err := r.FormFile("file")
	file, _, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//name := strings.Split(header.Filename, ".")
	//fmt.Printf("File name %s\n", name[0])
	io.Copy(&Buf, file)

	if deviceName == "" {
		fmt.Println("FAIL: something wrong. Empty file name")
		return
	}

	if electronic_component.Device_get(deviceName) == nil {
		err = electronic_component.Device_add(deviceName)
	}
	cur_device_id := electronic_component.Device_get(deviceName)["device_id"]
	fmt.Println("curDeviceIdx:", cur_device_id)

	contents := Buf.String()

	paramIndexes := make(map[string]int)
	for rowIdx, partNumber := range strings.Split(contents, "\n") {
		if rowIdx == 0 {
			for idx, paramName := range strings.Split(partNumber, ",") {
				//fmt.Println(idx, " ", paramName)
				if paramName == "Part Number" {
					paramIndexes["part_number"] = idx
				} else if paramName == "Quantity" {
					paramIndexes["qty"] = idx
				}
			}
			fmt.Println("paramIndexes: ", paramIndexes)
		} else {
			partNumberParam := strings.Split(partNumber, "\",\"")
			//fmt.Println(rowIdx, " pn: ", partNumber)
			//fmt.Println(rowIdx, len(partNumberParam))
			if len(partNumberParam) < 3 {
				continue
			}
			cur_part_number := strings.Replace(partNumberParam[paramIndexes["part_number"]], "\"", "", -1)
			cur_qty, _ := strconv.Atoi(partNumberParam[paramIndexes["qty"]])
			//fmt.Println(cur_part_number, cur_qty)
			bom_res := electronic_component.Devicebom_add(deviceName, cur_part_number, cur_qty)
			if bom_res != nil {
				res["success"] = false
				failed_PN.WriteString(cur_part_number + ", ")
			}
		}
	}

	Buf.Reset()

	failed_PN.WriteString("add them manually")
	if res["success"] == false {
		res["error"] = failed_PN.String()
	}
	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)

}

func apiCellUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("apiCellUpdate")
	fmt.Println(r.Form)

	vars := mux.Vars(r)
	req_idx := vars["cellidx"]
	req_action := r.FormValue("action")
	req_param := r.FormValue("param")
	req_val := r.FormValue("val")
	req_part_number := r.FormValue("part_number")
	fmt.Println(req_idx)
	fmt.Println(req_action)
	fmt.Println(req_param)
	fmt.Println(req_val)
	fmt.Println(req_part_number)

	res := map[string]interface{}{
		"success" : "true"}
	if req_idx == "undefined" {
		res["success"] = "false"
		res["error"] = "unknown id"
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	if req_param == "qty" || req_param == "notes" || req_param == "supplier_link" || req_param == "cellIdx" || req_param == "part_number" {
		err := Postgres.TableUpdateRow("cells", map[string]interface{}{req_param : req_val}, "cellidx = '" + req_idx + "'")
		if err != nil {
			res["success"] = false
			res["error"] = err
		}
	}
	if req_param == "description" || req_param == "ref_supplier" || req_param == "footprint" || req_param == "temperature" || req_param == "help_url"{
		err := Postgres.TableUpdateRow("part_numbers", map[string]interface{}{req_param : req_val}, "part_number = '" + req_part_number + "'")
		if err != nil {
			res["success"] = false
			res["error"] = err
		}
	}

	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)
}

func cellComponentsSoldered(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cellComponentsSoldered")
	req_cellIdx := r.FormValue("cellIdx")
	req_qty := r.FormValue("qty")

	res := map[string]interface{}{
		"success" : "true"}
	soldered_qty, _ := strconv.Atoi(req_qty)
	fmt.Println("req_cellIdx: ", req_cellIdx, " req_qty: ", soldered_qty)

	qty := electronic_component.Cell_getQty(req_cellIdx)
	fmt.Println("qty in stock: ", qty)
	if (qty < soldered_qty) {

		updateRes := electronic_component.Cell_setQty(req_cellIdx, 0)
		if (updateRes != nil) {
			res["success"] = false
			res["error"] = "not enough qty in cell" + updateRes.Error()
			jsonString, _ := json.Marshal(res)
			w.Write(jsonString)
			return
		}

		res["success"] = false
		res["error"] = "not enough qty in cell. Left " + strconv.Itoa(soldered_qty-qty)
		res["actualQty"] = strconv.Itoa(qty)
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	updateRes := electronic_component.Cell_setQty(req_cellIdx, qty - soldered_qty)
	if (updateRes != nil) {
		res["success"] = false
		res["error"] = updateRes
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}
	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)
}

func deviceBomUpdate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	req_idx := vars["device_id"]
	//req_action := r.FormValue("action")
	req_param := r.FormValue("param")
	req_val := r.FormValue("val")
	req_part_number := r.FormValue("part_number")

	res := map[string]interface{}{
		"success" : "true"}
	if req_idx == "undefined" {
		res["success"] = "false"
		res["error"] = "unknown device_id"
		jsonString, _ := json.Marshal(res)
		w.Write(jsonString)
		return
	}

	if req_param == "perDevice" {
		req_param = "qty"
	}
	fmt.Println("deviceBomUpdate")
	if req_param == "soldersequence" || req_param == "soldercell" || req_param == "qty" || req_param == "part_number" {
		fmt.Println("allowed")
		err := Postgres.TableUpdateRow("devicebom", map[string]interface{}{req_param : req_val}, "device_id = '" + req_idx + "'" + " AND part_number='" + req_part_number + "'")
		if err != nil {
			res["success"] = false
			res["error"] = err
		} else {
			res["data"] = "req_param " + req_param + " updated "
		}
	} else {
		fmt.Println("not allowed")
		res["success"] = false
		res["error"] = "fields allowed to update: part_number, perDevice, soldercell, soldersequence"
	}

	//if req_param == "description" || req_param == "ref_supplier" || req_param == "footprint" || req_param == "temperature" || req_param == "help_url"{
	//	err := Postgres.TableUpdateRow("part_numbers", map[string]interface{}{req_param : req_val}, "part_number = '" + req_part_number + "'")
	//	if err != nil {
	//		res["success"] = false
	//		res["error"] = err
	//	}
	//}

	jsonString, _ := json.Marshal(res)
	w.Write(jsonString)
}

func WebServer() {
	fs := http.FileServer(http.Dir("./ui/build"))
/*
	http.Handle("/", http.StripPrefix("/", fs))
	http.HandleFunc("/warehouse", getWarehouseData)
	http.ListenAndServe(":4900", nil)
//	http.ListenAndServeTLS(":4900", "/etc/letsencrypt/fullchain.pem", "/etc/letsencrypt/privkey.pem", router )
//*/

	fmt.Println("web")
	router := mux.NewRouter()
	router.PathPrefix("/admin/").Handler(http.StripPrefix("/admin/", fs ))
	router.HandleFunc("/warehouse", getWarehouseData).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/part_numbers/create", apiCreatePartNumber).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/part_numbers/createFast", apiCreatePartNumberFast).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/cells/create", apiCreateCell).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/cell/{cellidx}", apiCellUpdate).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/parcels", apiGetParcels).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/deviceUpload", deviceUpload).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/deviceUploadName", deviceUploadName).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/deviceList", deviceList).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/deviceBom", deviceBom).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/deviceBomUpdate/{device_id}", deviceBomUpdate).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/cellComponentsSoldered", cellComponentsSoldered).Methods("GET", "OPTIONS")
	http.ListenAndServe(":4900", router)
}

// api/v1/cellComponentsSoldered