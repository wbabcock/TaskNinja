//go:build windows
// +build windows

package utils

import (
	"golang.org/x/sys/windows/registry"
)

func GetLocaleDateFormat() (string, error) {
	// Open registry key for locale settings
	key, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\\International`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	// Query the "sShortDate" value, which indicates the date format
	dateFormat, _, err := key.GetStringValue("sShortDate")
	if err != nil {
		return "", err
	}

	return dateFormat, nil

}
