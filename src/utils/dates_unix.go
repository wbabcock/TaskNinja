//go:build !windows
// +build !windows

package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetLocaleDateFormat() (string, error) {
	cmd := exec.Command("locale", "-k", "LC_TIME")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	var locale string
	// Parse the output to find the `d_fmt` value
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "d_fmt") {
			locale = strings.ReplaceAll(strings.TrimPrefix(line, "d_fmt="), "\"", "")
			locale = strings.ReplaceAll(locale, "y", "Y")
		}
	}

	// Infer date format based on the locale
	// For simplicity, we assume common formats for certain locales
	switch locale {
	case "%m/%d/%Y":
		return "01/02/2006", nil
	case "%d/%m/%Y":
		return "02/01/2006", nil
	case "%Y/%d/%m":
		return "2006/02/01", nil
	case "%Y/%m/%d":
		return "2006/01/02", nil
	default:
		return "", fmt.Errorf("unknown date format")
	}
}
