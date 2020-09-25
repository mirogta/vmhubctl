package hub

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/mirogta/vmhubctl/oids"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/spf13/viper"
)

const defaultUserName = "admin"

var (
	automationID = "_n=12345"
	// NOTE: the second _=<value> param should be incremental for each request?
	adminRequestID = fmt.Sprintf("%s&_=1600201184287", automationID)

	// MIB details - VM Superhub model is 'ARRIS TG2492LG-85 Router'
	// see https://mibs.observium.org/mib/ARRIS-ROUTER-DEVICE-MIB/
	arrisRouterMIB = "1.3.6.1.4.1.4115.1.20.1"
	log            = logrus.New()
)

type Hub struct {
	routerIP          string
	baseURL           string
	base64Credentials string
	credentialCookie  *http.Cookie
	logLevel          logrus.Level
}

func NewHub() *Hub {
	hub := &Hub{
		logLevel: logrus.InfoLevel,
	}
	hub.initFromConfig()
	return hub
}

func init() {
	log.SetLevel(logrus.InfoLevel)
}

func (hub *Hub) initFromConfig() {

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	logLevel, ok := viper.Get("logLevel").(string)
	if ok {
		level, _ := logrus.ParseLevel(logLevel)
		hub.logLevel = level
		log.SetLevel(level)
	}

	logFile, ok := viper.Get("log").(string)
	if ok {
		rotateFileHook, _ := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   logFile,
			MaxSize:    5,
			MaxBackups: 7,
			MaxAge:     7,
			Level:      hub.logLevel,
			Formatter:  &logrus.TextFormatter{},
		})
		log.AddHook(rotateFileHook)
		log.Infof("Log: %s", logFile)
	}

	routerIP, ok := viper.Get("router_ip").(string)
	if ok == false {
		log.Fatalf("router_ip not set")
	}
	hub.routerIP = routerIP
	hub.baseURL = fmt.Sprintf("http://%s", routerIP)
	log.Infof("Router URL: %s", hub.baseURL)

	password, _ := viper.Get("password").(string)
	cred := fmt.Sprintf("%s:%s", defaultUserName, password)
	hub.base64Credentials = b64.URLEncoding.EncodeToString([]byte(cred))
}

func (hub *Hub) Login() error {
	url := fmt.Sprintf("%s/login?arg=%s&%s", hub.baseURL, hub.base64Credentials, adminRequestID)

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	var credential string
	if resp.StatusCode == http.StatusOK {
		log.Info("LOGGED IN")
		credential = getResponseBody(resp)
	}

	hub.credentialCookie = newCredCookie(hub.routerIP, credential)

	if credential == "" {
		return errors.New("Unable to log in")
	}

	return nil
}

func (hub Hub) Logout() {
	client := &http.Client{}
	logoutURL := fmt.Sprintf("%s/logout?%s", hub.baseURL, adminRequestID)
	req, err := http.NewRequest("GET", logoutURL, nil)
	if hub.credentialCookie != nil {
		req.AddCookie(hub.credentialCookie)
	}
	_, err = client.Do(req)
	if err != nil {
		log.Fatalf("LOGOUT ERROR: %v", err)
	}
	log.Info("LOGGED OUT")
}

func (hub Hub) ListAll() {
	// 	key := arrisRouterMIB + ".1.5.20"
	// 	value := hub.snmpWalk(arrisRouterMIB + key)
	// 	for k, v := range value {
	// 		fmt.Printf("%-36s:\t%s\n", k, v)
	// 	}

	for k, _ := range OIDs {
		oid := k
		title := strcase.ToLowerCamel(oid.String())
		_, formattedString := hub.snmpGetByOID(oid)
		fmt.Printf("%-30s:\t%s\n", title, formattedString)
	}
}

func (hub Hub) GetValue(oidName oids.Name) interface{} {
	value, _ := hub.snmpGetByOID(oidName)
	return value
}

func (hub Hub) SetWifi24GHzEnabled(enabled bool) {
	currentValue := hub.GetValue(oids.ArrisRouterBssActive24GHz)
	log.Info("Changing Wifi 2.4Ghz enabled to:", enabled)
	if currentValue.(bool) == enabled {
		log.Info("NOOP")
		return
	}

	key := arrisRouterMIB + ".1.3.22.1.3.10001"
	values := map[bool]string{
		true:  "1;2",
		false: "2;2",
	}
	value := values[enabled]
	response := hub.snmpSet(key, value)
	log.Debug(response)

	hub.applyChange()
}

func (hub Hub) SetWifi5GHzEnabled(enabled bool) {
	currentValue := hub.GetValue(oids.ArrisRouterBssActive5GHz)
	log.Info("Changing Wifi 5Ghz enabled to:", enabled)
	if currentValue.(bool) == enabled {
		log.Info("NOOP")
		return
	}

	key := arrisRouterMIB + ".1.3.22.1.3.10101"
	values := map[bool]string{
		true:  "1;2",
		false: "2;2",
	}
	value := values[enabled]
	response := hub.snmpSet(key, value)
	log.Debug(response)

	hub.applyChange()
}

func (hub Hub) Reboot() {
	title := "arrisRouterICtrlInitiateReboot"
	key := arrisRouterMIB + ".1.10.5.1"
	value := hub.snmpSet(key, "")
	fmt.Printf("%-30s:\t%s\n", title, value)
}

func (hub Hub) applyChange() {
	// run get before set
	key := arrisRouterMIB + ".1.9.0"
	response := hub.snmpGet(key)
	log.Debug(response)

	value := "1;2"
	response = hub.snmpSet(key, value)
	log.Debug(response)
}

func (hub Hub) snmpWalk(oids string) map[string]string {
	client := &http.Client{}
	url := fmt.Sprintf("%s/walk?oids=%s;%s", hub.baseURL, oids, automationID)
	log.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	req.AddCookie(hub.credentialCookie)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	result := getResponseBody(resp)

	byt := strings.Replace(result, "Error in OID formatting!", "", 1)
	var dat map[string]string
	if err := json.Unmarshal([]byte(byt), &dat); err != nil {
		panic(err)
	}

	return dat
}

func (hub Hub) snmpGetByOID(oidName oids.Name) (interface{}, string) {
	key := OIDs[oidName]
	value := hub.snmpGet(string(key))
	formatter := OIDsFormatter[oidName]
	if formatter == nil {
		return value, value
	}
	return formatter(value)
}

func (hub Hub) snmpGet(oid string) string {
	client := &http.Client{}
	url := fmt.Sprintf("%s/snmpGet?oids=%s;%s", hub.baseURL, oid, automationID)
	log.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	req.AddCookie(hub.credentialCookie)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	result := getResponseBody(resp)

	byt := strings.Replace(result, "Error in OID formatting!", "", 1)
	var dat map[string]string
	if err := json.Unmarshal([]byte(byt), &dat); err != nil {
		panic(err)
	}

	return dat[oid]
}

func (hub Hub) snmpSet(oid string, value string) string {
	client := &http.Client{}
	url := fmt.Sprintf("%s/snmpSet?oid=%s=%s;&%s", hub.baseURL, oid, value, automationID)
	log.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	req.AddCookie(hub.credentialCookie)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	result := getResponseBody(resp)

	byt := strings.Replace(result, "Error in OID formatting!", "", 1)
	return byt
}

func getResponseBody(resp *http.Response) string {
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyBytes)
}

func newCredCookie(domain string, credential string) *http.Cookie {
	return &http.Cookie{
		Name:   "credential",
		Domain: domain,
		Path:   "/",
		Value:  credential,
		// RawExpires: "Session",
	}
}
