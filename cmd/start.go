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
	"bufio"
	"database/sql"
	"fmt"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
		db, err := sql.Open("mysql", DBSource())
		if err != nil {
			panic(err)
		}

		// Open doesn't open a connection. Validate DSN data:
		err = db.Ping()
		if err != nil {
			panic(err)
		} else {
			fmt.Println("DB connection working")
		}

		//try to create table with for the host. If table already exists, continue.
		stmt, err := db.Prepare("CREATE TABLE " + Hostname() + "(Cpu float, Mem int(11), Uptime text, Time datetime);")
		if err != nil {
			log.Fatal(err)
		}
		_,err = stmt.Exec()
		if err != nil {
			if strings.HasSuffix(err.Error(), "already exists") {
				fmt.Println("Table exists in DB for " + Hostname())
			} else {
				log.Fatal(err)
			}
		} else {
			fmt.Println("table created for " + Hostname())
		}

		//constantly send metrics to DB
		for true {
			stmt, err := db.Prepare("INSERT INTO " + Hostname() + " VALUES (" + strconv.FormatFloat(getCpu(0), 'f', 6, 64) + "," + strconv.FormatUint(getMemory(1), 10) + "," + "'" + getUptime() + "'" + "," + "CURRENT_TIMESTAMP" + ");")
			if err != nil {
				log.Fatal(err)
			}

			_, err = stmt.Exec()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("updated stats successfully")
			}

			time.Sleep(time.Duration(5) * time.Second)
		}
	},
}

func getMemory(i int) uint64 {
	mem, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	if i == 0 {
		return mem.Total / 1024 / 1024
	} else if i == 1 {
		return mem.Used / 1024 / 1024
	} else if i == 2 {
		return mem.Cached / 1024 / 1024
	} else if i == 3 {
		return mem.Free / 1024 / 1024
	}
	return 0
}

func getCpu(i int) float64 {
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	total := float64(after.Total - before.Total)
	if i == 0 {
		return float64(after.User-before.User) / total * 100
	} else if i == 1 {
		return float64(after.System-before.System) / total * 100
	} else if i == 2 {
		return float64(after.Idle-before.Idle) / total * 100
	}
	return 0
}

func getUptime() string {
	out, err := exec.Command("uptime", "-p").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return ""
	}
	output := fmt.Sprintf("%s", out)
	output = strings.TrimSuffix(output, "\n")
	return output
}

func DBSource() string {
	file, err := os.Open(".login.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	var user string
	var pass string
	var ip string
	for scanner.Scan() {
		if i == 0 {
			user = scanner.Text()
		} else if i == 1 {
			pass = scanner.Text()
		} else if i == 2 {
			ip = scanner.Text()
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	output := user + ":" + pass + "@tcp(" + ip + ")/stats"

	return output
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
