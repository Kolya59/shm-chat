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
	"log"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kolya59/shm-chat/pkg/client"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start chat client",
	Run: func(cmd *cobra.Command, args []string) {
		readerID, err := strconv.Atoi(cmd.Flag("reader id").Value.String())
		if err != nil {
			log.Println("Invalid args")
			return
		}
		if readerID == 0 {
			log.Println("Invalid reader id")
			return
		}

		writerID, err := strconv.Atoi(cmd.Flag("writer id").Value.String())
		if err != nil {
			log.Println("Invalid args")
			return
		}
		if writerID == 0 {
			log.Println("Invalid writer id")
			return
		}

		client.StartClient(readerID, writerID)
	},
	Args: cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().IntP("reader id", "r", 0, "Reader segment id")
	clientCmd.PersistentFlags().IntP("writer id", "w", 0, "Writer segment id")
}
