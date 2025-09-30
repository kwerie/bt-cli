package main

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var deviceRegex = `((?:[0-9A-Fa-f]{2}[:-]){5}(?:[0-9A-Fa-f]{2}))\s(\S.*)`

var connectCommand = &cobra.Command{
	Use:                "connect",
	Short:              "Shows a menu of available bluetooth devices to connect to",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		connectedDevices := exec.Command("bluetoothctl", "devices", "Connected")
		var connectedDevicesOut bytes.Buffer
		connectedDevices.Stdout = &connectedDevicesOut

		if err := connectedDevices.Run(); err != nil {
			fmt.Println("Error executing bluetoothctl:", err)
			return
		}

		var connected = connectedDevicesOut.String()
		var connectedList = strings.Split(connected, "\n")

		devices := exec.Command("bluetoothctl", "devices")
		var availableDevicesOut bytes.Buffer
		devices.Stdout = &availableDevicesOut

		if err := devices.Run(); err != nil {
			fmt.Println("Error executing bluetoothctl:", err)
			return
		}

		fmt.Println("Available Bluetooth devices:")
		var output = availableDevicesOut.String()
		var split = strings.Split(output, "\n")

		// filter out connected devices from available devices
		// TODO: this can be improved by a lot
		var filtered []string
		for _, device := range split {
			skip := false
			for _, connectedDevice := range connectedList {
				if device == connectedDevice && len(device) > 0 {
					skip = true
					break
				}
			}
			if !skip && len(device) > 0 {
				filtered = append(filtered, device)
			}
		}
		split = filtered
		var choice string
		prompt := &survey.Select{
			Message: "Which device do you want to connect to?",
			Options: split,
		}

		if err := survey.AskOne(prompt, &choice); err != nil {
			fmt.Println("Error selecting device:", err)
			return
		}

		if len(choice) == 0 {
			return
		}

		if matched, _ := regexp.MatchString(deviceRegex, choice); matched {
			re := regexp.MustCompile(deviceRegex)
			matches := re.FindStringSubmatch(choice)
			if len(matches) < 2 {
				return
			}
			var macAddress = matches[1]
			//var name = matches[2]

			//fmt.Printf("Connecting to %s \n", name)

			exec.Command("bluetoothctl", "connect", macAddress).Run()
		}
	},
}

var disconnectCommand = &cobra.Command{
	Use:                "disconnect",
	Short:              "Shows a menu of available bluetooth devices to disconnect from",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		devices := exec.Command("bluetoothctl", "devices", "Connected")
		// to get currently connected devices: bluetoothctl devices Connected

		var out bytes.Buffer
		devices.Stdout = &out

		err := devices.Run()

		if err != nil {
			fmt.Println("Error executing bluetoothctl:", err)
			return
		}

		fmt.Println("Available Bluetooth devices:")
		var output = out.String()
		var split = strings.Split(output, "\n")
		var choice string
		prompt := &survey.Select{
			Message: "Which device do you want to disconnect from?",
			Options: split,
		}

		if err := survey.AskOne(prompt, &choice); err != nil {
			fmt.Println("Error selecting device:", err)
			return
		}

		if len(choice) == 0 {
			return
		}

		if matched, _ := regexp.MatchString(deviceRegex, choice); matched {
			re := regexp.MustCompile(deviceRegex)
			matches := re.FindStringSubmatch(choice)
			if len(matches) < 2 {
				return
			}
			var macAddress = matches[1]
			//var name = matches[2]

			//fmt.Printf("Connecting to %s \n", name)

			exec.Command("bluetoothctl", "disconnect", macAddress).Run()
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "bt-cli [args]",
	Short: "A CLI tool for managing bluetooth connections",
	Args:  cobra.ExactArgs(1),
}

func main() {
	rootCmd.AddCommand(connectCommand)
	rootCmd.AddCommand(disconnectCommand)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
