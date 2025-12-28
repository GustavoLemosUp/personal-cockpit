package main

const (
	AppVersion  = "1.0.0"
	AppName     = "Personal Cockpit"
	BuildNumber = "368"
	CurrentYear = "2025"
)

// GetFullVersion retorna versão completa
func GetFullVersion() string {
	if BuildNumber != "" {
		return AppVersion + "b" + BuildNumber
	}
	return AppVersion
}
