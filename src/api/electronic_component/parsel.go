package electronic_component

import (
	"../../storage/postgresql"
)

func InitParsel() {
	keys_order = []string{"id", "order_id", "index_in_parcel", "store", "tracking_number", "part_number", "supplier_link", "status", "qty"}

	cellInfo := make(map[string]string)
	cellInfo[keys_order[0]] = "serial"
	cellInfo[keys_order[1]] = "text"
	cellInfo[keys_order[2]] = "int"
	cellInfo[keys_order[3]] = "text"
	cellInfo[keys_order[4]] = "text"
	cellInfo[keys_order[5]] = "text"
	cellInfo[keys_order[6]] = "text"
	cellInfo[keys_order[7]] = "text"
	cellInfo[keys_order[8]] = "int"

	fKeys := []string{
		"FOREIGN KEY (part_number) REFERENCES part_numbers(part_number)",
	}

	Postgres.TableCreate__("parcel", keys_order, cellInfo, "", fKeys)
	//Postgres.TableInsertRow("parcel", map[string]interface{}{"order_id":"282547731", "index_in_parcel":1, "store":"digikey", "tracking_number":"012606077453968", "part_number":"MCIMX6Y0CVM05AA", "supplier_link":"https://www.digikey.com/product-detail/en/nxp-usa-inc/MCIMX6Y0CVM05AA/568-13074-ND/6556996", "status":"handled", "qty":5})

}
