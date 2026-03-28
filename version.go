package main

import "personal-cockpit/config"

// GetFullVersion returns the full version string from config.
func GetFullVersion() string {
	return config.MustLoad().FullVersion()
}
