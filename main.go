package main

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	// "regexp"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "bt-cli [args]",
	Short: "A CLI tool for managing bluetooth connections",
	// TODO: implement argument
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Bluetooth CLI tool executed with argument:", args[0])
		if args[0] == "help" {
			cmd.Help()
		}

		if args[0] == "connect" {
			// TODO: get a list of available bt devices and show menu to select one using bluetoothctl
			cmd := exec.Command("bluetoothctl", "devices")

			var out bytes.Buffer
			cmd.Stdout = &out

			err := cmd.Run()

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

			fmt.Printf("You selected: %s\n", choice)

			// for _, line := range split {
			// 	// match mac address based on regex
			// 	var regex = `((?:[0-9A-Fa-f]{2}[:-]){5}(?:[0-9A-Fa-f]{2}))\s(\S.*)`
			// 	if matched, _ := regexp.MatchString(regex, line); matched {
			// 		re := regexp.MustCompile(regex)
			// 		matches := re.FindStringSubmatch(line)
			// 		if len(matches) < 2 {
			// 			continue
			// 		}

			// 		var macAddress = matches[1]
			// 		var name = matches[2]

			// 		// print them
			// 		fmt.Printf("MAC Address: %s, Name: %s\n", macAddress, name)
			// 		// show a selection menu via cobra
			// 	}
			// }
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
