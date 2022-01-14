package selenium_cl

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"runtime"
	"strconv"
	"strings"
)

var service *selenium.Service

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
//				"--headless",
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

func ServiceStop() {
	service.Stop()
}



func Switch_frame(wd selenium.WebDriver, css_selector string) error {
	_, err := wd.FindElement(selenium.ByCSSSelector, css_selector)
	if err != nil {
//		panic(err)
		return nil
	}
	//	src, err := elem.GetAttribute("src")
	//	fmt.Println("!", src)

	frame_name := strings.Replace(css_selector, "#", "", -1)

	if err := wd.SwitchFrame(frame_name); err != nil {
		return nil
	} else {
		//		fmt.Println("Switching OK")
	}
	return nil
}

func Enter_data_to_field(wd selenium.WebDriver, css_selector string, data string) (selenium.WebElement, error) {
	elem, err := wd.FindElement(selenium.ByCSSSelector, css_selector)
	if err != nil {
		return nil, err
		//		panic(err)
	}
	if data != "" {
		err = elem.SendKeys(data)
	}

	if err != nil {
		return elem, err
		//		panic(err)
	}
	return elem, nil
}

func Click_button(wd selenium.WebDriver, css_selector string) (selenium.WebElement, error) {
	btn, err := wd.FindElement(selenium.ByCSSSelector, css_selector)
	if err != nil {
		panic(err)
	}
	if err := btn.Click(); err != nil {
		panic(err)
	}
	return btn, nil
}
