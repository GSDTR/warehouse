package digikey

import (
	"fmt"
	"github.com/tebeka/selenium"
	"time"
)

import "../../selenium_cl"

type configuration struct {
	promelec_login string
	promelec_password string
}
var config configuration

func Init(conf map[string]string) {
	config.promelec_login = "almaz1c"// conf["aliexpress_login"]
	config.promelec_password = "Reutov20,100%"//conf["aliexpress_password"]
}

func Autorization(wd selenium.WebDriver) (selenium.WebDriver, error) {
	if err := wd.Get("https://www.digikey.com/MyDigiKey/Login"); err != nil {
		//if err = wd.Get("https://m.ru.aliexpress.com/login.html"); err != nil {
		fmt.Println("ps:", err)
		panic(err)
	}
	time.Sleep(15000 * time.Millisecond)

	ss, _ := wd.PageSource()
	fmt.Println("ps:", ss)
	err := selenium_cl.Switch_frame(wd, "#frmLogin")
	if err != nil {
		fmt.Println("digikey frame switch failed: ", err)
		ss, _ := wd.PageSource()
		fmt.Println("ps:", ss)
		return nil, err
	}
	fmt.Println("digikey frame switched to login window")

	_, err = selenium_cl.Enter_data_to_field(wd, "input[name='username']", config.promelec_login)
	if err != nil {
		fmt.Println("digikey enter username failed: ", err)
		ss, _ := wd.PageSource()
		fmt.Println("ps:", ss)
		return nil, err
	}
	fmt.Println("promelec login entered")


	_, err = selenium_cl.Enter_data_to_field(wd, "input[name='Password']", config.promelec_password)
	if err != nil {
		fmt.Println("digikey enter password failed: ", err)
		return nil, err
	}
	fmt.Println("promelec password entered")


	_, err = selenium_cl.Click_button(wd, "input[id='btnPostLogin']")
	if err != nil {
		fmt.Println("digikey autorization failed: ", err)
		return nil, err
	}
	fmt.Println("digikey autorization completed")

	return wd, nil
}

func GetOrdersList(wd selenium.WebDriver, amount int) ([]map[string]interface{}, selenium.WebDriver, error) {
	if amount == -1 {
		amount = 99999
	}
	//	cnt := 0
	//	resLinks := make([]map[string]interface{}, 1, 1)

	if err := wd.Get("https://www.promelec.ru/cabinet/listorders/"); err != nil {
		//if err = wd.Get("https://m.ru.aliexpress.com/login.html"); err != nil {
		panic(err)
	}

	elem_expand_btns, _ := wd.FindElements(selenium.ByCSSSelector, "a[class='f-bu f-bu-default']")
	for _, elem_btn := range elem_expand_btns {
		//fmt.Println("elem_btn: ", elem_btn)
		elem_btn.Click()
		time.Sleep(1000 * time.Millisecond)
	}

	elem_pn_lists, _ := wd.FindElements(selenium.ByCSSSelector, "tr[class='hideline order_active_tr show']")
	for _, elem_pn_list := range elem_pn_lists {
		elem_table, err := elem_pn_list.FindElement(selenium.ByCSSSelector, "table[class='f-table-zebra table_show']")
		if err != nil {	return nil, wd, err	}
		pn_rows, err := elem_table.FindElements(selenium.ByCSSSelector, "tr")
		for _, pn_row := range pn_rows {
			fmt.Println(pn_row.Text())
		}
	}

	time.Sleep(5000 * time.Millisecond)
	return nil, nil, nil
}