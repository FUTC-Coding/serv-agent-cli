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
	_ "github.com/go-sql-driver/mysql"
	"github.com/mackerelio/go-osstat/network"
	cpu2 "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
		stmt, err := db.Prepare("CREATE TABLE " + Hostname() + "(CpuPerc float, MemTotal int(11), " +
			"MemUsed int(11), MemCached int(11), MemFree int(11), RxBytes int(11), TxBytes int(11), " +
			"DiskUsed int(11), DiskFree int(11), DiskRead int(11), DiskWrite int(11), Uptime text, Time datetime);")
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
			stmt, err := db.Prepare("INSERT INTO " + Hostname() + " VALUES (" +
				strconv.FormatFloat(getCpu(), 'f', 6, 64) + "," +
				strconv.FormatUint(getMemory(0), 10) + "," +
				strconv.FormatUint(getMemory(1), 10) + "," +
				strconv.FormatUint(getMemory(2), 10) + "," +
				strconv.FormatUint(getMemory(3), 10) + "," +
				strconv.FormatUint(getNetwork(0), 10) + "," +
				strconv.FormatUint(getNetwork(1), 10) + "," +
				strconv.FormatUint(getDisk(0), 10) + "," +
				strconv.FormatUint(getDisk(1), 10) + "," +
				strconv.FormatUint(getDiskIO(0), 10) + "," +
				strconv.FormatUint(getDiskIO(1), 10) + "," +
				"'" + getUptime() + "'" + "," + "CURRENT_TIMESTAMP" + ");")
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
	v, e := mem.VirtualMemory()
	if e != nil {
		log.Fatal("couldn't get Memory stats")
	}

	if i == 0 {
		return v.Total / 1024 / 1024
	} else if i == 1 {
		return v.Used / 1024 / 1024
	} else if i == 2 {
		return v.Cached / 1024 / 1024
	} else if i == 3 {
		return v.Free / 1024 / 1024
	}
	return 0
}

func getCpu() float64 {
	c,_ := cpu2.Percent(time.Second * 1, false)
	return c[0]
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

func getNetwork(i int) (uint64){
	before, err := network.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := network.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	if i == 0 {
		return after[0].RxBytes - before[0].RxBytes
	}
	if i == 1 {
		return after[0].TxBytes - before[0].TxBytes
	}

	return 0
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

func getDisk(i int) uint64 {
	d,e := disk.Usage("/")
	if e != nil {
		log.Fatal("couldn't get Disk Info")
	}
	if i == 0 {
		return d.Used/1024/1024
	}
	if i == 1 {
		return d.Free/1024/1024
	}
	return 0
}

func getDiskIO(i int) uint64 {
	d := disk.IOCountersStat{}
	if i == 0 {
		return d.ReadBytes
	}
	if i == 1 {
		return d.WriteBytes
	}
	return 0
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
