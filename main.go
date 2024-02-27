package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/flan6/ddns/config"
	"github.com/flan6/ddns/godaddy"
)

func main() {
	if err := LoadEnv(".env"); err != nil {
		fmt.Println(err)
		return
	}

	go config.WatchConfig("config/config.json")

	conf := config.GetConfig()
	dns := godaddy.NewGoDaddy(os.Getenv("api_key"), os.Getenv("api_secret"))

	domains, err := dns.Domains()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(domains)

	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Println("Error getting IP:", err)
		return
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading IP:", err)
		return
	}

	for {
		records, err := dns.RecordsByType("qual.work", "A", "@")
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := range records {
			records[i].Data = string(ip)
		}

		if err := dns.SetRecordsByType("qual.work", "A", "@", records); err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(time.Minute * time.Duration(conf.UpdateInterval))
	}
}

func LoadEnv(envPath string) error {
	file, err := os.Open(envPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if parts := strings.SplitN(scanner.Text(), "=", 2); len(parts) == 2 {
			if err := os.Setenv(parts[0], parts[1]); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}
