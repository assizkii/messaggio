/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/assizkii/messaggio/internal/adapters/storages"
	"github.com/assizkii/messaggio/internal/queue"
	"github.com/assizkii/messaggio/internal/utils"
	"github.com/streadway/amqp"

	"github.com/spf13/cobra"
)

// consumerServeCmd represents the consumerServe command
var consumerServeCmd = &cobra.Command{
	Use:   "consumer",
	Short: "A consumer server",
	Long:  `Run a consumer queue`,
	Run: func(cmd *cobra.Command, args []string) {
		appConf := utils.GetAppConfig()

		dsn := "host=" + appConf.DbHost + " user=" + appConf.DbUser + " password=" + appConf.DbPassword + " dbname=" + appConf.DbName + "  sslmode=disable"
		storage := storages.New(dsn)

		conn, err := amqp.Dial(appConf.AmpqHost)
		queue.HandleError(err, "Can't connect to AMQP")

		defer conn.Close()

		amqpChannel, err := conn.Channel()
		queue.HandleError(err, "Can't create a amqpChannel")
		defer amqpChannel.Close()

		queue.RunConsumer(amqpChannel, storage)

	},
}

func init() {
	rootCmd.AddCommand(consumerServeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumerServeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumerServeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
