package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Config struct {
	Port   string
	KeyDir string
}

func requestHandler(config *Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slashSeperated := strings.Split(r.URL.Path[1:], "/")
		key := slashSeperated[0]
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", config.KeyDir, key))
		if err != nil {
			log.Printf("I %v %v %v %v %v", r.RemoteAddr, r.Method, r.URL.Path, 404, key)
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}
		stringData := strings.TrimSpace(string(data))
		reply := KVPair{Key: key, Value: stringData}
		log.Printf("I %v %v %v %v %v", r.RemoteAddr, r.Method, r.URL.Path, 200, key)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reply)
	})
}

func main() {
	viper.BindEnv("Port", "PORT")
	viper.BindEnv("KeyDir", "KEY_DIRECTORY")
	viper.SetDefault("Port", "8080")
	viper.SetDefault("KeyDir", "/keys")
	var Config Config
	viper.Unmarshal(&Config)

	http.HandleFunc("/", requestHandler(&Config))
	log.Printf("I Started on port %v reading keys from %v", Config.Port, Config.KeyDir)
	panic(http.ListenAndServe(":"+Config.Port, nil))
}
