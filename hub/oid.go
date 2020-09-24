package hub

import (
	"strconv"
	"time"

	"github.com/mirogta/vmhubctl/oids"
)

type OID string

var OIDs = map[oids.Name]OID{
	oids.SysUpTime:                  "1.3.6.1.2.1.1.3.0",
	oids.ArrisRouterHardwareVersion: OID(arrisRouterMIB + ".1.5.10.0"),
	oids.ArrisRouterFirmwareVersion: OID(arrisRouterMIB + ".1.5.11.0"),
	// oids.ArrisRouterAdminTimeout:    OID(arrisRouterMIB + ".1.5.2.0"),
	oids.ArrisRouterBssSSID: OID(arrisRouterMIB + ".1.3.22.1.2.10001"),
	// oids.ArrisRouterWPAPreSharedKey: OID(arrisRouterMIB + ".1.3.26.1.2.10001"),
	oids.ArrisRouterBssActive24GHz: OID(arrisRouterMIB + ".1.3.22.1.3.10001"),
	oids.ArrisRouterBssActive5GHz:  OID(arrisRouterMIB + ".1.3.22.1.3.10101"),
}

type FormatterFunc func(value string) (interface{}, string)

var OIDsFormatter = map[oids.Name]FormatterFunc{
	oids.SysUpTime:                 durationFormatter,
	oids.ArrisRouterBssActive24GHz: wifiStatusFormatter,
	oids.ArrisRouterBssActive5GHz:  wifiStatusFormatter,
}

func durationFormatter(value string) (interface{}, string) {
	// multiply *10 to get in milliseconds and divide by 1000 to get seconds
	// ... so divide by 100
	milliseconds, _ := strconv.Atoi(value)
	seconds := int64(milliseconds / 100.0)
	t := time.Unix(seconds, 0)
	return t, t.Format("15:04:05")
}

func wifiStatusFormatter(value string) (interface{}, string) {
	labels := map[string]string{
		"1": "enabled",
		"2": "disabled",
	}
	text := labels[value]
	boolValue := text == "enabled"
	return boolValue, text
}
