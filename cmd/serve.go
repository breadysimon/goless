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
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/breadysimon/goless/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type XX struct {
	ID      int    `gorm:"AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	Message string `gorm:"type:varchar(50);unique_index" json:"name" rest:"search"`
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long:  `.`,
	Run: func(cmd *cobra.Command, args []string) {
		s := rest.NewRestApi(&XX{}).
			Connect("sqlite3", ":memory:").
			Serve(viper.GetString("addr"), viper.GetInt("port"), "/api/v1/", true)
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)

		// Waiting for SIGINT (pkill -2)
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("addr", "a", "127.0.0.1", "bind address")
	serveCmd.Flags().IntP("port", "p", 8128, "bind port")
	viper.BindPFlags(serveCmd.Flags())
}
