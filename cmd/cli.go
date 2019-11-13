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

	"github.com/breadysimon/goless/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "cli get|set|delete",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		baseUrl := fmt.Sprintf("http://%s:%d", viper.GetString("addr"), viper.GetInt("port"))
		c := rest.NewClient(baseUrl, "/api/v1/", "admin", "admin")
		switch args[0] {
		case "get":
			if viper.GetInt("id") != 0 {
				x := XX{}
				x.ID = viper.GetInt("id")
				c.Find(&x)
				fmt.Println(x)
			} else {
				xx := []XX{}
				c.Find(&xx)
				fmt.Println(xx)
			}

		case "set":
			x := XX{Message: viper.GetString("message")}
			id := viper.GetInt("id")
			if id != 0 {
				x.ID = id
				fmt.Println("err=", c.Update(&x).Error())
			} else {
				fmt.Println("err=", c.Create(&x).Error())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cliCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cliCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	cliCmd.Flags().StringP("addr", "a", "127.0.0.1", "bind address")
	cliCmd.Flags().IntP("port", "p", 8128, "bind port")
	cliCmd.Flags().IntP("id", "i", 0, "get by id")
	cliCmd.Flags().StringP("message", "m", "none", "set message")
	viper.BindPFlags(cliCmd.Flags())
}
