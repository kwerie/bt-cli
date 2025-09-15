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

var connectCommand = &cobra.Command{
	Use:                "connect",
	Short:              "Shows a menu of available bluetooth devices to connect to",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: filter out already connected devices
		devices := exec.Command("bluetoothctl", "devices")

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

		var regex = `((?:[0-9A-Fa-f]{2}[:-]){5}(?:[0-9A-Fa-f]{2}))\s(\S.*)`

		if matched, _ := regexp.MatchString(regex, choice); matched {
			re := regexp.MustCompile(regex)
			matches := re.FindStringSubmatch(choice)
			if len(matches) < 2 {
				return
			}
			var macAddress = matches[1]
			var name = matches[2]

			fmt.Printf("Connecting to %s \n", name)

			exec.Command("bluetoothctl", "connect", macAddress).Run()
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
