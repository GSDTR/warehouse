package Digikey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

import "../../common"

type configuration struct {
	digikey_client_id string
	digikey_client_secret string
	digikey_redirect_uri string
	digikey_redirect_uri_short string
	digikey_username string
	digikey_password string
	digikey_autorization_code string
	digikey_access_token string
	digikey_refresh_token string
	digikey_token_expires_in int
}
var config configuration

func Init(conf map[string]string) {
	config.digikey_client_id = conf["digikey_client_id"]
	config.digikey_redirect_uri = conf["digikey_redirect_uri"]
	config.digikey_username = conf["digikey_username"]
	config.digikey_password = conf["digikey_password"]
	config.digikey_client_secret = conf["digikey_client_secret"]
	config.digikey_autorization_code = conf["digikey_autorization_code"]
	config.digikey_access_token = conf["digikey_access_token"]

//	fmt.Println("redirect_uri:", config.digikey_redirect_uri)
	r, _ := regexp.Compile("\\/[a-zA-Z0-9]*\\/$")
	config.digikey_redirect_uri_short = r.FindString(config.digikey_redirect_uri)
//	fmt.Println("short_url: ", config.digikey_redirect_uri_short)
//	fmt.Println("digikey_autorization_code: ", config.digikey_autorization_code)
//	fmt.Println("digikey_client_id: ", config.digikey_client_id)
//	fmt.Println("digikey_access_token: ", config.digikey_access_token)
}

func redirectServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	fmt.Printf("request: %s %s\r\n", req.Host, req.URL.Path)

	code, ok := req.URL.Query()["code"]

	if !ok || len(code[0]) < 1 {
		return
	} else {
		fmt.Print("=========================================code: ")
		fmt.Println(code[0])
		config.digikey_autorization_code = code[0]
	}
}

func Run_redirect_URI_server() {
	fmt.Println("============================Digikey redirect URI server started")
	http.HandleFunc(config.digikey_redirect_uri_short, redirectServer)
	err := http.ListenAndServeTLS(":8085", "fullchain.pem", "privkey.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func Autorization_get_code() {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	const (
		// These paths will be different on your system.
		seleniumPath    = "vendor/selenium-server-standalone-3.14.0.jar"
		geckoDriverPath = "vendor/geckodriver-v0.23.0-linux64"
		port            = 4444
	)
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	//        selenium.SetDebug(true)
	selenium.SetDebug(false) // to reduce amount of logs
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
//	caps := selenium.Capabilities{"browserName": "firefox"}
	caps := selenium.Capabilities{"browserName": "firefox", "acceptInsecureCerts":true }

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// Navigate to the simple playground interface.
	if err := wd.Get("https://sso.digikey.com/as/authorization.oauth2?response_type=code&client_id=" + config.digikey_client_id + "&redirect_uri=" + config.digikey_redirect_uri); err != nil {
		panic(err)
	}
	
	elem, err := wd.FindElement(selenium.ByCSSSelector, "#username")
	if err != nil {
		panic(err)
	}

	// Remove the boilerplate code already in the text box.
	if err := elem.Clear(); err != nil {
		panic(err)
	}

	// Enter some new code in text box.
	err = elem.SendKeys(config.digikey_username)
	if err != nil {
		panic(err)
	}

	elem, err = wd.FindElement(selenium.ByCSSSelector, "#password")
	if err != nil {
		panic(err)
	}

	// Remove the boilerplate password already in the text box.
	if err := elem.Clear(); err != nil {
		panic(err)
	}

	// Enter new valid in text box.
	err = elem.SendKeys(config.digikey_password)
	if err != nil {
		panic(err)
	}

	// Click the run button.
	btn, err := wd.FindElement(selenium.ByXPATH, "/html/body/div/div[2]/div/form/div[6]/a[2]")
	if err != nil {
		panic(err)
	}
	if err := btn.Click(); err != nil {
		panic(err)
	}

	// Wait for the program to finish running and get the output.
	btn, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div[2]/form/div/div[3]/a[2]")
	if err != nil {
		log.Println(err)
	} else {
		if err := btn.Click(); err != nil {
			panic(err)
		}
	}
}

func  Access_get_token() {

	request_url := "https://sso.digikey.com/as/token.oauth2"
	form := url.Values{
		"code": {config.digikey_autorization_code},
		"client_id":  {config.digikey_client_id},
		"client_secret":  {config.digikey_client_secret},
		"redirect_uri":     {config.digikey_redirect_uri},
		"grant_type":     {"authorization_code"},
	}

	body := bytes.NewBufferString(form.Encode())
	rsp, err := http.Post(request_url, "application/x-www-form-urlencoded", body)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	body_byte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Print( "got tokens: ")
	fmt.Println(string(body_byte))

//	body_byte := []byte("{\"access_token\":\"Rgn8ZEE4JLTVCkLSFNF7urgne41h\",\"refresh_token\":\"aOtz1TWxJdIozKB0uAGUoeq2ykK30erLLl04T4pftG\",\"token_type\":\"Bearer\",\"expires_in\":86400}")
	var tokens_json interface{}
	err = json.Unmarshal(body_byte, &tokens_json)
	if err != nil {
		log.Println(err)
	}
	fmt.Print("tokens json: ")
	fmt.Println(tokens_json)
	config.digikey_access_token = tokens_json.(map[string]interface{})["access_token"].(string)
	config.digikey_refresh_token = tokens_json.(map[string]interface{})["refresh_token"].(string)
	token_expires_in := tokens_json.(map[string]interface{})["expires_in"].(float64)
//	fmt.Println(token_expires_in)
	config.digikey_token_expires_in = int(token_expires_in)
//	fmt.Println("access tokern: ", config.digikey_access_token)
//	fmt.Println("refresh_token: ", config.digikey_refresh_token)
//	fmt.Println("token_expires_in: ", config.digikey_token_expires_in)
}


var flag = 0

func Get_part_number(part_number string) (interface{}, error) {
//	fmt.Println("digikey_client_id: ",  config.digikey_client_id)
//	fmt.Println("digikey_access_token: ",  config.digikey_access_token)

	url := "https://api.digikey.com/services/partsearch/v2/partdetails"

	payload := strings.NewReader("{\"Part\":\"" + part_number + "\",\"IncludeAllAssociatedProducts\":\"false\",\"IncludeAllForUseWithProducts\":\"false\"}")

	req, err := http.NewRequest("POST", url, payload)

	if( err != nil) {
		fmt.Println("ERR:", err)
		return nil, err
	}
	req.Header.Add("x-ibm-client-id", config.digikey_client_id) // client id
	req.Header.Add("x-digikey-locale-site", "US")
	req.Header.Add("x-digikey-locale-language", "en")
	req.Header.Add("x-digikey-locale-currency", "USD")
	req.Header.Add("x-digikey-locale-shiptocountry", "RU")
	//      req.Header.Add("x-digikey-customer-id", "REPLACE_THIS_VALUE")
	//      req.Header.Add("x-digikey-partner-id", "REPLACE_THIS_VALUE")
	req.Header.Add("authorization", config.digikey_access_token) //access token
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if( res.Status != "200 OK") {
		fmt.Println("ERR: ", string(body) )
		return nil, Common.UserError(res.Status)
	}

	//fmt.Println(string(body))
	var part_number_json interface{}
	err = json.Unmarshal(body, &part_number_json)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//fmt.Println(part_number_json)
	return part_number_json, nil
	//verifyJSON(part_number_json)

}


func Get_order_history() {
/*
	fmt.Println("digikey_client_id: ",  config.digikey_client_id)
	fmt.Println("digikey_access_token: ",  config.digikey_access_token)

	url := "https://api.digikey.com/services/orderhistory/v1/customersalesorderhistory/9988036/2016-01-01/2018-12-12/undefined"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-ibm-client-id", "e9a31104-81f2-4110-918b-4e558068f42b") //config.digikey_client_id) // client id
	req.Header.Add("authorization", "pzbuUtcsEy5J41LJffAUJXaiRy6w") //config.digikey_access_token) //access token

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	*/

	url := "https://api.digikey.com/services/orderhistory/v1/rootsalesorderhistory/9988036/2016-01-01/2018-12-12/undefined"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-ibm-client-id", "e9a31104-81f2-4110-918b-4e558068f42b")
	req.Header.Add("authorization", "pzbuUtcsEy5J41LJffAUJXaiRy6w")
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}


func Get_order_status() {
	/*
		fmt.Println("digikey_client_id: ",  config.digikey_client_id)
		fmt.Println("digikey_access_token: ",  config.digikey_access_token)

		url := "https://api.digikey.com/services/orderhistory/v1/customersalesorderhistory/9988036/2016-01-01/2018-12-12/undefined"

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("x-ibm-client-id", "e9a31104-81f2-4110-918b-4e558068f42b") //config.digikey_client_id) // client id
		req.Header.Add("authorization", "pzbuUtcsEy5J41LJffAUJXaiRy6w") //config.digikey_access_token) //access token

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		fmt.Println(res)
		fmt.Println(string(body))
		*/

	url := "https://api.digikey.com/services/orderStatus/v2/orderStatus/9988036/51967160?rootAccount=false"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-ibm-client-id", "e9a31104-81f2-4110-918b-4e558068f42b")
	req.Header.Add("authorization", "pzbuUtcsEy5J41LJffAUJXaiRy6w")
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
