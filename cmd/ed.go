/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/breadysimon/goless/edit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// edCmd represents the ed command
var edCmd = &cobra.Command{
	Use:   "ed filename",
	Short: "edit text",
	Long:  `edit text.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		script := viper.GetString("script")
		yaml := viper.GetString("yaml")
		inplace := viper.GetBool("inplace")
		output := viper.GetString("output")
		if output == "" {
			if inplace {
				output = file
			}
		}
		if yaml != "" {
			if err := edit.RunYaml([]byte(yaml), file, output); err != nil {
				fmt.Printf("ERROR: %v", err)
			}
		} else if script != "" {
			if err := edit.RunScript(script, file, output); err != nil {
				fmt.Printf("ERROR: %v", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(edCmd)

	edCmd.Flags().StringP("script", "s", "-", "script file, default stdin")
	edCmd.Flags().StringP("yaml", "y", "", "yaml script")
	edCmd.Flags().StringP("output", "o", "", "output file")
	edCmd.Flags().BoolP("inplace", "i", false, "edit and write back to then file")
	viper.BindPFlags(edCmd.Flags())
}
