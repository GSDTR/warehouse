package electronic_component

import (
	"../../encoding/json_a"
	"../../storage/postgresql"
	"../digikey"
	"fmt"
	"regexp"
)

var keys_order []string

func Get_keys_order() []string {
	return keys_order
}
func Init_() error {
	keys_order = []string{"part_number", "description", "footprint", "temperature", "manufacturer",
		"ref_supplier", "help_url", "status", "component_type", "family_type", "component_series", "mfg_footpriint", "detailed_description",
		"price_1_pcs", "price_10_pcs", "price_100_pcs", "price_1000_pcs", "price_5000_pcs", "datasheet_url"}

	component, err := Get_info("STM32L051C8T6")
	if(err != nil) {
		return err
	}
	Postgres.TableCreateBySnapshot("part_numbers", Get_keys_order(), component, "part_number")
	return nil
	//	fmt.Println(keys_order)
}

func Init() error {
	keys_order = []string{"part_number", "description", "footprint", "temperature", "manufacturer",
		"ref_supplier", "help_url", "status", "component_type", "family_type", "component_series", "mfg_footpriint", "detailed_description",
		"price_1_pcs", "price_10_pcs", "price_100_pcs", "price_1000_pcs", "price_5000_pcs", "datasheet_url"}

	cellInfo := make(map[string]string)
	cellInfo[keys_order[0]] = "text"
	cellInfo[keys_order[1]] = "text"
	cellInfo[keys_order[2]] = "int"
	cellInfo[keys_order[3]] = "text"
	cellInfo[keys_order[4]] = "text"
	cellInfo[keys_order[5]] = "text"
	cellInfo[keys_order[6]] = "text"
	cellInfo[keys_order[7]] = "text"
	cellInfo[keys_order[8]] = "text"
	cellInfo[keys_order[9]] = "text"
	cellInfo[keys_order[10]] = "text"
	cellInfo[keys_order[11]] = "text"
	cellInfo[keys_order[12]] = "text"
	cellInfo[keys_order[13]] = "real"
	cellInfo[keys_order[14]] = "real"
	cellInfo[keys_order[15]] = "text"
	cellInfo[keys_order[16]] = "real"
	cellInfo[keys_order[17]] = "text"
	cellInfo[keys_order[18]] = "text"

	Postgres.TableCreate_("part_numbers", keys_order, cellInfo, "part_number")
	return nil

}

func Get_info(part_number string) (map[string]interface{}, error) {
//	var component electronic_component;
	electronic_component := make(map[string]interface{})

	part_number_json, err := Digikey.Get_part_number(part_number)
	if( err != nil) {
		return nil, err
	}
//	fmt.Println("pn: ", part_number_json)
	electronic_component["part_number"] = json_a.Get(part_number_json, "PartDetails.ManufacturerPartNumber")
	electronic_component["detailed_description"] = json_a.Get( part_number_json, "PartDetails.DetailedDescription")
	electronic_component["description"] = json_a.Get( part_number_json, "PartDetails.ProductDescription")
	electronic_component["help_url"] = "https://www.digikey.com" + json_a.Get(part_number_json, "PartDetails.PartUrl").(string)
	electronic_component["manufacturer"] = json_a.Get( part_number_json, "PartDetails.ManufacturerName.Text")
	electronic_component["component_type"] = json_a.Get( part_number_json, "PartDetails.Category.Text")
	electronic_component["component_series"] = json_a.GetMap( part_number_json, "PartDetails.Series")["Value"]
	electronic_component["family_type"] = json_a.Get( part_number_json, "PartDetails.Family.Text")
	electronic_component["datasheet_url"] = json_a.Get( part_number_json, "PartDetails.PrimaryDatasheet")
	electronic_component["ref_supplier"] = " "

	electronic_component["price_1_pcs"] = "0.0"
	electronic_component["price_10_pcs"] = "0.0"
	electronic_component["price_100_pcs"] = "0.0"
	electronic_component["price_1000_pcs"] = "0.0"
	electronic_component["price_5000_pcs"] = "0.0"
	tmp1 := json_a.GetArr( part_number_json, "PartDetails.StandardPricing")
//	fmt.Println("\tprices: ", tmp1)
	for _, elem := range tmp1 {
//		fmt.Println("\t\t", elem.(map[string]interface{})["BreakQuantity"] )
		switch  elem.(map[string]interface{})["BreakQuantity"].(float64) {
		case 1:	electronic_component["price_1_pcs"] = elem.(map[string]interface{})["UnitPrice"]
		case 10: electronic_component["price_10_pcs"] = elem.(map[string]interface{})["UnitPrice"]
		case 100: electronic_component["price_100_pcs"] = elem.(map[string]interface{})["UnitPrice"]
		case 1000: electronic_component["price_1000_pcs"] = elem.(map[string]interface{})["UnitPrice"]
		case 5000: electronic_component["price_5000_pcs"] = elem.(map[string]interface{})["UnitPrice"]
		default:
			//fmt.Println("\t\tprice ? pcs: ", elem.(map[string]interface{})["BreakQuantity"] )
		}
	}

	tmp2 := json_a.GetArr( part_number_json, "PartDetails.Parameters")
	for _, elem := range tmp2 {
//		fmt.Println("\t", elem.(map[string]interface{})["Parameter"], " = ", elem.(map[string]interface{})["Value"])
		switch elem.(map[string]interface{})["Parameter"] {
		case "Operating Temperature":	electronic_component["temperature"] = elem.(map[string]interface{})["Value"]
		case "Supplier Device Package": electronic_component["mfg_footpriint"] = elem.(map[string]interface{})["Value"]
		case "Package / Case": electronic_component["footprint"] = elem.(map[string]interface{})["Value"]
		case "Part Status": electronic_component["status"] = elem.(map[string]interface{})["Value"]
		case "Tolerance": electronic_component["tolerance"] = elem.(map[string]interface{})["Value"]
		case "Vgs(th) (Max) @ Id":
			r, _ := regexp.Compile("[0-9\\.]+V")
			match := r.FindString( elem.(map[string]interface{})["Value"].(string) )
			descr := electronic_component["description"].(string) + " Vgs(th)" + match
			electronic_component["description"] = descr
//			fmt.Println("descr: ", electronic_component["description"])
		}
	}
//	fmt.Println("component: ", electronic_component)
	return electronic_component, nil
}

func PartNumber_add(pn string) error {
	component, err := Get_info(pn)
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
		fmt.Println("WARNING: unable to find ", pn)
		err := fmt.Errorf("WARNING: unable to find %s", pn)
		return err
	}
	return nil
}

func PartNumber_get(pn string) map[string]interface{} {
	part_number := Postgres.TableGetRow("part_numbers", "*", "WHERE part_number='"+pn+"'")
	return part_number
}

func PartNumber_is_exist(pn string) bool {
	part_number := Postgres.TableGetRow("part_numbers", "part_number", "WHERE part_number='"+pn+"'")
	if part_number == nil {
		return false
	} else {
		return true
	}
}