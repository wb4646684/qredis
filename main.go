package main

import (
	"flag"
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	var redisHost string
	var keyName string
	var scanCount int64
	var action string
	flag.StringVar(&redisHost, "host", "127.0.0.1:6379", "-host 127.0.0.1:6379")
	flag.StringVar(&keyName, "key", "test", "-key test*")
	flag.Int64Var(&scanCount, "count", 1000, "-count 1000")
	flag.StringVar(&action, "action", "scan", "-action scan|delete")
	flag.Parse()

	redisDB := redis.NewClient(&redis.Options{
		Addr: redisHost,
	})
	if _, err := redisDB.Ping().Result(); err != nil {
		panic(err)
	}

	switch action {
	case "scan":
		fmt.Println(*ScanKey(redisHost, keyName, scanCount, redisDB))
	case "delete":
		fmt.Println(DeleteKey(ScanKey(redisHost, keyName, scanCount, redisDB), redisDB))
	}
	defer redisDB.Close()
}

func ScanKey(redisHost string, keyName string, scanCount int64, redisDB *redis.Client) *[]string {
	var cursor uint64 = 0
	var matchKey []string
	var matchKeyList []string

	for {
		matchKey, cursor = redisDB.Scan(cursor, keyName, scanCount).Val()
		matchKeyList = append(matchKeyList, matchKey...)
		if cursor == 0 {
			return &matchKeyList
		}
	}
}

func DeleteKey(keyNameList *[]string, redisDB *redis.Client) *redis.IntCmd {
	return redisDB.Del(*keyNameList...)
}
