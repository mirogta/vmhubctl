package oids

type Name int

const (
	SysUpTime Name = iota
	ArrisRouterHardwareVersion
	ArrisRouterFirmwareVersion
	ArrisRouterAdminTimeout
	ArrisRouterBssSSID
	ArrisRouterWPAPreSharedKey
	ArrisRouterBssActive24GHz
	ArrisRouterBssActive5GHz
)
