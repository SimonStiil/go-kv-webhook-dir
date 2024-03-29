package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Config struct {
	Port    string
	KeyDir  string
	LogJson bool
}

func requestHandler(config *Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slashSeperated := strings.Split(r.URL.Path[1:], "/")
		key := slashSeperated[0]
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", config.KeyDir, key))
		if err != nil {
			if config.LogJson {
				json.NewEncoder(os.Stdout).Encode(
					JsonRequest{Timestamp: time.Now().Format(time.StampMilli), Type: "info", RemoteAddr: r.RemoteAddr, Method: r.Method, URLPath: r.URL.Path, Status: 404, Key: key})
			} else {
				log.Printf("I %v %v %v %v %v", r.RemoteAddr, r.Method, r.URL.Path, 404, key)
			}
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}
		stringData := strings.TrimSpace(string(data))
		reply := KVPair{Key: key, Value: stringData}

		if config.LogJson {
			json.NewEncoder(os.Stdout).Encode(
				JsonRequest{Timestamp: time.Now().Format(time.StampMilli), Type: "info", RemoteAddr: r.RemoteAddr, Method: r.Method, URLPath: r.URL.Path, Status: 200, Key: key})
		} else {
			log.Printf("I %v %v %v %v %v", r.RemoteAddr, r.Method, r.URL.Path, 200, key)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reply)
	})
}

type JsonRequest struct {
	Timestamp  string `json:"timestamp"`
	Type       string `json:"type"`
	RemoteAddr string `json:"remoteAddr"`
	Method     string `json:"message"`
	URLPath    string `json:"urlPath"`
	Status     int    `json:"status"`
	Key        string `json:"key"`
}

func main() {
	viper.BindEnv("Port", "PORT")
	viper.BindEnv("KeyDir", "KEY_DIRECTORY")
	viper.BindEnv("LogJson", "LOG_JSON")
	viper.SetDefault("Port", "8080")
	viper.SetDefault("KeyDir", "/keys")
	viper.SetDefault("LogJson", false)
	var Config Config
	viper.Unmarshal(&Config)

	http.HandleFunc("/", requestHandler(&Config))
	if Config.LogJson {
		log.Printf("I Started on port %v reading keys from %v", Config.Port, Config.KeyDir)
	} else {
		json.NewEncoder(os.Stdout).Encode(JsonStart{Timestamp: time.Now().Format(time.StampMilli), Type: "info", Message: "Server Started", Port: Config.Port, KeyDir: Config.KeyDir})
	}
	panic(http.ListenAndServe(":"+Config.Port, nil))
}

type JsonStart struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	Port      string `json:"port"`
	KeyDir    string `json:"keyDir"`
}
