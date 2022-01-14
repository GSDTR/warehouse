package electronic_component

import (
	"../../storage/postgresql"
	"fmt"
)

func InitDevice() {
	keys_order = []string{"device_id", "device_name"}

	cellInfo := make(map[string]string)
	cellInfo[keys_order[0]] = "serial"
	cellInfo[keys_order[1]] = "text"

	//	fKeys := []string{
	//		"FOREIGN KEY (part_number) REFERENCES part_numbers(part_number)",
	//	}

	Postgres.TableCreate_("device", keys_order, cellInfo, "device_id")
	//Postgres.TableInsertRow("parcel", map[string]interface{}{"order_id":"282547731", "index_in_parcel":1, "store":"digikey", "tracking_number":"012606077453968", "part_number":"MCIMX6Y0CVM05AA", "supplier_link":"https://www.digikey.com/product-detail/en/nxp-usa-inc/MCIMX6Y0CVM05AA/568-13074-ND/6556996", "status":"handled", "qty":5})
}

func InitDeviceBom() {
	keys_order = []string{"device_id", "part_number", "qty"}

	cellInfo := make(map[string]string)
	cellInfo[keys_order[0]] = "int"
	cellInfo[keys_order[1]] = "text"
	cellInfo[keys_order[2]] = "int"
	//cellInfo[keys_order[3]] = "text"
	//cellInfo[keys_order[4]] = "int"

	fKeys := []string{
		"FOREIGN KEY (part_number) REFERENCES part_numbers(part_number)",
		"FOREIGN KEY (device_id) REFERENCES device(device_id)",
	}

	Postgres.TableCreate__("deviceBom", keys_order, cellInfo, "", fKeys)
	//Postgres.TableInsertRow("parcel", map[string]interface{}{"order_id":"282547731", "index_in_parcel":1, "store":"digikey", "tracking_number":"012606077453968", "part_number":"MCIMX6Y0CVM05AA", "supplier_link":"https://www.digikey.com/product-detail/en/nxp-usa-inc/MCIMX6Y0CVM05AA/568-13074-ND/6556996", "status":"handled", "qty":5})
}

func AlterDeviceBom() {
	res := Postgres.AlterTable_AddColumn("deviceBom", "solderCell", "text", "REFERENCES cells (cellIdx)")
	if res != nil {
		fmt.Println("FAIL:", res)
	}
}

func AlterDeviceBom2() {
	res := Postgres.AlterTable_AddColumn("deviceBom", "solderSequence", "int", "")
	if res != nil {
		fmt.Println("FAIL:", res)
	}
}

func AlterDeviceBom3() {
	res := Postgres.AlterTable_AddConstraint("deviceBom", "solderCell", "FOREIGN KEY (solderCell) REFERENCES cells (cellIdx)")
	if res != nil {
		fmt.Println("FAIL:", res)
	}
}

// ALTER TABLE distributors ADD CONSTRAINT distfk FOREIGN KEY (address) REFERENCES addresses (address) MATCH FULL;

func Devicebom_add(deviceName string, partNumber string, qty int) error {

	newPartNumber := make(map[string]interface{})
	newPartNumber["device_id"] = Device_get(deviceName)["device_id"]

	if !PartNumber_is_exist(partNumber) {
		//fmt.Println("new part number: ", partNumber)
		res := PartNumber_add(partNumber)
		if res != nil {
			fmt.Println("FAIL. Unable to add new part number: ", res)
			return res
		}
	}
	newPartNumber["part_number"] = partNumber
	newPartNumber["qty"] = qty
	//fmt.Println("addToBom: ", newPartNumber)
	res := Postgres.TableInsertRow("deviceBom", newPartNumber)
	if res != nil {
		fmt.Println("FAIL: ", res)
		return res
	}
	return nil
}

func Device_add(device string) error {

	newDevice := make(map[string]interface{})
	newDevice["device_name"] = device

	res := Postgres.TableInsertRow("device", newDevice)
	if res != nil {
		fmt.Println("FAIL: ", res)
	}
	return nil
}

func Device_get(deviceName string) map[string]interface{} {
	device := Postgres.TableGetRow("device", "*", "WHERE device_name='"+deviceName+"'")
	if device == nil {
		fmt.Println(" no device", deviceName, " found")
	}
	return device
}

func Device_getAll() []map[string]interface{} {
	deviceList := Postgres.TableGetRows("device", "*", "")
	return deviceList
}

func DeviceBom_getByIdx(idx string) []map[string]interface{} {
	deviceBom := Postgres.TableGetRows("deviceBom", "part_number, qty", "where device_id="+idx)
	return deviceBom
}

func DeviceBomStock_getByIdx(idx string) []map[string]interface{} {
	//deviceBom := Postgres.TableGetRows("deviceBom", "devicebom.part_number, devicebom.qty, sum(cells.qty) as warehouseStock", "LEFT JOIN cells ON devicebom.part_number=cells.part_number WHERE devicebom.device_id=1 GROUP BY devicebom.device_id, devicebom.qty, devicebom.part_number")
	deviceBom := Postgres.TableGetRows("deviceBom", "devicebom.part_number, devicebom.qty, devicebom.soldercell, devicebom.soldersequence, part_numbers.description, sum(cells.qty) as instock, sum(cells.qty) / devicebom.qty as solderable", "LEFT JOIN cells ON devicebom.part_number=cells.part_number INNER JOIN part_numbers on devicebom.part_number=part_numbers.part_number WHERE devicebom.device_id=" + idx + " GROUP BY devicebom.device_id, devicebom.qty, devicebom.part_number, devicebom.soldercell, devicebom.soldersequence, part_numbers.description")
	return deviceBom
}

// SELECT devicebom.part_number, devicebom.qty, part_numbers.description, sum(cells.qty) as warehouseStock FROM devicebom LEFT JOIN cells ON devicebom.part_number=cells.part_number INNER JOIN part_numbers on devicebom.part_number=part_numbers.part_number WHERE devicebom.device_id=1 GROUP BY devicebom.device_id, devicebom.qty, devicebom.part_number, part_numbers.description;
