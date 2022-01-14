package electronic_component

import (
	"../../storage/postgresql"
	"strconv"
)

func InitHolder() {
	keys_order = []string{"cellIdx", "part_number", "qty", "supplier_link", "notes"}

	cellInfo := make(map[string]string)
	cellInfo[keys_order[0]] = "text"
	cellInfo[keys_order[1]] = "text"
	cellInfo[keys_order[2]] = "int"
	cellInfo[keys_order[3]] = "text"
	cellInfo[keys_order[4]] = "text"

	Postgres.TableCreate_("cells", keys_order, cellInfo, "cellIdx")
}

//func addCell(cellInfo map[string]interface{}) error {
//
//}

func Cell_getQty(cellIdx string) int {
	res := Postgres.TableGetRow("cells", "qty", "WHERE cellIdx='"+cellIdx+"'")
	return int(res["qty"].(int64))
}

func Cell_setQty(cellIdx string, solderedQty int) error {
	qty := strconv.Itoa(solderedQty)
	err := Postgres.TableUpdateRow("cells", map[string]interface{}{"qty" : qty}, "cellidx = '" + cellIdx + "'")
	return err
}