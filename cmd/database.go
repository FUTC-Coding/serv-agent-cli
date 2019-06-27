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
	"log"
	"os"

	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("database called")
		fmt.Println("username: " + user)
		fmt.Println("password: " + password)
		fmt.Println("ip and port: " + ip)
		writeCredentials()
	},
}

var user string
var password string
var ip string

func writeCredentials() {
	f, err := os.Create(".login.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Fprint(f, user + "\n" + password + "\n" + ip)
}

func init() {
	rootCmd.AddCommand(databaseCmd)
	databaseCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "DB username")
	databaseCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "DB password")
	databaseCmd.PersistentFlags().StringVarP(&ip, "ip", "i", "", "-i [ip]:[port]")
	databaseCmd.MarkFlagRequired("user")
	databaseCmd.MarkFlagRequired("password")
	databaseCmd.MarkFlagRequired("ip")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// databaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// databaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}