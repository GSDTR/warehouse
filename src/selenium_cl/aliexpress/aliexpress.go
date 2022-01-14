package Aliexpress

import (
	"fmt"
	"github.com/tebeka/selenium"
	"regexp"
	"strconv"
	"strings"
	"time"
)

import "../../api/email"

import "../../selenium_cl"

type configuration struct {
	aliexpress_login string
	aliexpress_password string
}
var config configuration

func Init(conf map[string]string) {
	config.aliexpress_login = "almaz_1c@mail.ru"// conf["aliexpress_login"]
	config.aliexpress_password = "Reutov20,100%"//conf["aliexpress_password"]
}


//var driver = "chrome"

func verification() string {
	msgList := email.Read_messages("almaz_1c@mail.ru", "Reutov20,100%", 5)
	code := ""
	for _, msg := range msgList {
		if( msg["subject"] == "Verification Code From Alibaba Group") {
			r, _ := regexp.Compile(">[0-9]+<\\/span>")
			code = r.FindString(msg["body"])
			code = strings.Replace(code, ">", "", -1)
			code = strings.Replace(code, "</span", "", -1)
		}
	}
	return code
}


func GetOrdersList(wd selenium.WebDriver, amount int) ([]map[string]interface{}, selenium.WebDriver, error) {
	if amount == -1 {
		amount = 99999
	}
	cnt := 0
	resLinks := make([]map[string]interface{}, 1, 1)

	if err := wd.Get("https://trade.aliexpress.com/orderList.htm"); err != nil {
		//if err = wd.Get("https://m.ru.aliexpress.com/login.html"); err != nil {
		panic(err)
	}

	for true {
		elemOrderTable, err := wd.FindElement(selenium.ByCSSSelector, "#buyer-ordertable")
		if err != nil {
			return nil, wd, err
		}

		elem_wrappers, err := elemOrderTable.FindElements(selenium.ByCSSSelector, ".order-item-wraper")
		if err != nil {
			fmt.Println("1")
			return nil, wd, err
		}

		for _, elem := range elem_wrappers {
			elemOrderDetail, err := elem.FindElement(selenium.ByCSSSelector, "a[class='view-detail-link']")
			if err != nil {
				fmt.Println("2")
				return nil, wd, err
			}
			orderLink, _ := elemOrderDetail.GetAttribute("href")
			fmt.Println(orderLink)

			elemPartNumbers, err := elem.FindElements(selenium.ByXPATH, ".//tr[@class='order-body']")
//			elemPartNumberDetails, err := elem.FindElements(selenium.ByXPATH, ".//p[@class='product-title']//a[@class='baobei-name']")
			if err != nil {
				fmt.Println("3")
				return nil, wd, err
			}
			for _, elemPn := range elemPartNumbers {
				pn, _ := elemPn.FindElement(selenium.ByXPATH, ".//p[@class='product-title']//a[@class='baobei-name']")
				pnLink, _ := pn.GetAttribute("href")
				fmt.Println(pnLink)

				pnSub := ""
				elemPnSub, err := elemPn.FindElement(selenium.ByXPATH, ".//p[@class='product-property']//span[@class='val']")
				if err == nil {
					pnSub, _ = elemPnSub.Text()
					fmt.Println(pnSub)
				}

				pnPriceAmount := ""
				pnPrice := 0.0
				pnAmount := 0
				elemPnPrice, err := elemPn.FindElement(selenium.ByXPATH, ".//p[@class='product-amount']")
				if err == nil {
					pnPriceAmount, _ = elemPnPrice.Text()
					fmt.Println(pnPriceAmount)

					r, _ := regexp.Compile("\\$ [0-9.]*")
					tmp := r.FindString(pnPriceAmount)
					tmp = strings.Replace(tmp, "$ ", "", -1)
					pnPrice, _ = strconv.ParseFloat(tmp, 64)
					fmt.Println(pnPrice)

					r, _ = regexp.Compile("X[0-9]+")
					tmp = r.FindString(pnPriceAmount)
					tmp = strings.Replace(tmp, "X", "", -1)
					pnAmount, _ = strconv.Atoi(tmp)
					fmt.Println(pnAmount)
				}

				curLink := make(map[string]interface{})
				curLink["orderLink"] = orderLink
				curLink["partNumberLink"] = pnLink
				curLink["partNumberSub"] = pnSub
				curLink["partNumberPrice"] = pnPrice
				curLink["partNumberAmount"] = pnAmount
				resLinks = append(resLinks, curLink)
			}
			fmt.Println("\r\n")
			cnt++
			if cnt >= amount {
				fmt.Println("4")
				return resLinks, wd, nil
			}
		}

		elemsNextBtn, err := wd.FindElements(selenium.ByCSSSelector, "a[class='ui-pagination-next ui-goto-page']")
		if err != nil {
			fmt.Println("5")
			return resLinks, wd, nil
		}
		if len(elemsNextBtn) == 2 {
			elemsNextBtn[1].Click()
		} else {
			fmt.Println("6")
			return resLinks, wd, nil
		}
	}
	return resLinks, wd, nil
}


func Autorization(wd selenium.WebDriver) (selenium.WebDriver, error){
//	service, wd, err := InitWebDriver()
/*-----------------------------------------------------------------------------*/

	if err := wd.Get("https://login.aliexpress.com/buyer_ru.htm"); err != nil {
	//if err = wd.Get("https://m.ru.aliexpress.com/login.html"); err != nil {
		panic(err)
	}

//	var elem selenium.WebElement

	err := selenium_cl.Switch_frame(wd , "#alibaba-login-box")
	if err == nil {	fmt.Println("aliexpress frame switched to login window")}

	_, err = selenium_cl.Enter_data_to_field(wd, "#fm-login-id", "almaz_1c@mail.ru")
	if err == nil {	fmt.Println("aliexpress login entered")}

	_, err = selenium_cl.Enter_data_to_field(wd, "#fm-login-password", "Reutov20,100%")
	if err == nil {	fmt.Println("aliexpress password entered")}

	_, err = selenium_cl.Click_button(wd, "#fm-login-submit")
	if err == nil {	fmt.Println("aliexpress login btn clicked")}

	time.Sleep(5000 * time.Millisecond)
	_, err = selenium_cl.Enter_data_to_field(wd, "#search-key", "hello world")
	if err == nil {
		fmt.Println("aliexpress search field found")
		return wd, nil
	} else {
		fmt.Println("aliexpress search field NOT found")
	}

	elemJcode, err := selenium_cl.Enter_data_to_field(wd, "#J_Checkcode", "")
	if err == nil {
		fmt.Println("J_Checkcode field found!")
		time.Sleep(10000 * time.Millisecond)
		code := verification()
		fmt.Println("verification code:", code)
		err = elemJcode.SendKeys(code)
		if err != nil {
			fmt.Println("unable to enter verification code")
			return wd, err
		}
		_, err := selenium_cl.Click_button(wd, "button[type='submit']")
		if err == nil {
			fmt.Println("aliexpress verification code btn clicked")
		} else {
			fmt.Println("unable to click verification code btn")
			return wd, err
		}
	}

	time.Sleep(5000 * time.Millisecond)
	_, err = selenium_cl.Enter_data_to_field(wd, "#search-key", "hello world")
	if err == nil {
		fmt.Println("aliexpress search field FOUND!")
		return wd, nil
	} else {
		fmt.Println("aliexpress search field NOT found")
	}

	time.Sleep(15000 * time.Millisecond)
	fmt.Println("\r\n\r\n\r\n")
	fmt.Println(wd.PageSource())
	fmt.Println("\r\n\r\n\r\n")
	return wd, nil
}


//var service *selenium.Service
//var wd selenium.WebDriver

/*
func InitWebDriver() (selenium.WebDriver, error) {
	const (
		// These paths will be different on your system.
		seleniumPath    = "vendor/selenium-server-standalone-3.14.0.jar"
		geckoDriverPath = "vendor/geckodriver-v0.23.0-linux64"
		port            = 4444
	)

	//opts := []selenium.ServiceOption{
	//	selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	//	selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	//	selenium.Output(os.Stderr),            // Output debug information to STDERR.
	//}
	////        selenium.SetDebug(true)
	//selenium.SetDebug(false) // to reduce amount of logs
	//service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	//if err != nil {
	//	panic(err) // panic is used only as an example and is not otherwise recommended.
	//}
	//defer service.Stop()
	//
	//// Connect to the WebDriver instance running locally.
	//caps := selenium.Capabilities{"browserName": "firefox"}
	//wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	//if err != nil {
	//	panic(err)
	//}
	//defer wd.Quit()


	var err error
	var opts []selenium.ServiceOption
	if runtime.GOOS == "windows" {
		service, err = selenium.NewChromeDriverService("vendor\\chromedriver.exe",
			port, opts...)
	} else {
		service, err = selenium.NewChromeDriverService("vendor/chromedriver-linux64-2.42",
			port, opts...)
	}
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	var chrCaps chrome.Capabilities
	if runtime.GOOS == "windows" {
		chrCaps = chrome.Capabilities{
			Args: []string{
				"--headless",
				"--no-sandbox",
			},
			Path: "vendor\\GoogleChromePortable\\App\\Chrome-bin\\chrome.exe",
		}
	} else {
		chrCaps = chrome.Capabilities{
			Args: []string{
				"--headless",
				"--no-sandbox",
			},
			Path: "vendor/chrome-linux/chrome",
		}
	}
	caps.AddChrome(chrCaps)

	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(port)+"/wd/hub")
	if err != nil {
		panic(err)
		return nil, err
	}
	//	defer wd.Quit()
	return wd, nil
}
//*/
