package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/flan6/ddns/config"
	"github.com/flan6/ddns/godaddy"
)

func main() {
	const configPath = "config/config.json"

	if err := LoadEnv(".env"); err != nil {
		fmt.Println(err)
		return
	}

	go config.WatchConfig(configPath)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dns := godaddy.NewGoDaddy(os.Getenv("api_key"), os.Getenv("api_secret"))

	domains, err := dns.Domains()
	if err != nil {
		logger.Error("Failed to get domains", err)
		return
	}

	logger.Info("Domains:", "domains", domains)
	logger.Info("Starting to update records....")
	for {
		resp, err := http.Get("https://api.ipify.org")
		if err != nil {
			logger.Error("Failed to get IP:", err)
			return
		}

		ip, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Failed to read IP:", err)
			return
		}
		resp.Body.Close()

		records, err := dns.RecordsByType("qual.work", "A", "@")
		if err != nil {
			logger.Error("Failed to get records", err)
			return
		}

		for i := range records {
			records[i].Data = string(ip)
		}

		if err := dns.SetRecordsByType("qual.work", "A", "@", records); err != nil {
			logger.Error("Failed to set records", err)
			return
		}

		logger.Info("Updated records", "ip", string(ip), "records", records)
		time.Sleep(time.Minute * time.Duration(config.GetConfig().UpdateInterval))
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
