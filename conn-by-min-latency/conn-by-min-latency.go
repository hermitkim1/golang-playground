package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	model "github.com/cloud-barista/cb-larva/poc-cb-net/pkg/cb-network/model"
	"github.com/cloud-barista/cb-larva/poc-cb-net/pkg/file"
	cblog "github.com/cloud-barista/cb-log"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// CBLogger represents a logger to show execution processes according to the logging level.
var CBLogger *logrus.Logger
var config model.Config

func init() {
	// Set cb-log
	env := os.Getenv("CBLOG_ROOT")
	if env != "" {
		// Load cb-log config from the environment variable path (default)
		fmt.Printf("CBLOG_ROOT: %v\n", env)
		CBLogger = cblog.GetLogger("cb-network")
	} else {

		// Load cb-log config from the current directory (usually for the production)
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exePath := filepath.Dir(ex)
		// fmt.Printf("exe path: %v\n", exePath)

		logConfPath := filepath.Join(exePath, "config", "log_conf.yaml")
		if file.Exists(logConfPath) {
			fmt.Printf("path of log_conf.yaml: %v\n", logConfPath)
			CBLogger = cblog.GetLoggerWithConfigPath("cb-network", logConfPath)

		} else {
			// Load cb-log config from the project directory (usually for development)
			logConfPath = filepath.Join(exePath, "..", "..", "config", "log_conf.yaml")
			if file.Exists(logConfPath) {
				fmt.Printf("path of log_conf.yaml: %v\n", logConfPath)
				CBLogger = cblog.GetLoggerWithConfigPath("cb-network", logConfPath)
			} else {
				err := errors.New("fail to load log_conf.yaml")
				panic(err)
			}
		}
		CBLogger.Debugf("Load %v", logConfPath)
	}
	// Load cb-network config from the current directory (usually for the production)
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	// fmt.Printf("exe path: %v\n", exePath)

	configPath := filepath.Join(exePath, "config", "config.yaml")
	if file.Exists(configPath) {
		fmt.Printf("path of config.yaml: %v\n", configPath)
		config, _ = model.LoadConfig(configPath)
	} else {
		// Load cb-network config from the project directory (usually for the development)
		configPath = filepath.Join(exePath, "..", "..", "config", "config.yaml")

		if file.Exists(configPath) {
			config, _ = model.LoadConfig(configPath)
		} else {
			err := errors.New("fail to load config.yaml")
			panic(err)
		}
	}

	CBLogger.Debugf("Load %v", configPath)

	fmt.Println("End......... init() of controller.go")
	fmt.Println("")
}

func main() {

	endpoints := config.ETCD.Endpoints
	ctx := context.TODO()

	// Performance evaluation
	// Target 0
	cli0, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoints[0]},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		CBLogger.Fatal(err)
	}

	defer func() {
		errClose := cli0.Close()
		if errClose != nil {
			CBLogger.Fatal("Can't close the etcd client", errClose)
		}
	}()

	//ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	stp00 := time.Now()
	cli0.Status(ctx, endpoints[0])
	fmt.Println("etcd Status took", time.Since(stp00))

	stp02 := time.Now()
	cli0.Get(ctx, "sample_key_0")
	fmt.Println("etcd Get took", time.Since(stp02))

	stp01 := time.Now()
	cli0.Put(ctx, "sample_key_0", "sample_value_0")
	fmt.Println("etcd Put took", time.Since(stp01))

	fmt.Println()

	// Target 1
	cli1, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoints[1]},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		CBLogger.Fatal(err)
	}

	defer func() {
		errClose := cli1.Close()
		if errClose != nil {
			CBLogger.Fatal("Can't close the etcd client", errClose)
		}
	}()

	//ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	stp10 := time.Now()
	cli1.Status(ctx, endpoints[1])
	fmt.Println("etcd Status took", time.Since(stp10))

	stp12 := time.Now()
	_, err = cli1.Get(ctx, "sample_key_1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("etcd Get took", time.Since(stp12))

	stp11 := time.Now()
	_, err = cli1.Put(ctx, "sample_key_1", "sample_value_1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("etcd Put took", time.Since(stp11))

	fmt.Println()

	// Target 2
	cli2, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoints[2]},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		CBLogger.Fatal(err)
	}

	defer func() {
		errClose := cli2.Close()
		if errClose != nil {
			CBLogger.Fatal("Can't close the etcd client", errClose)
		}
	}()

	//ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	stp20 := time.Now()
	cli2.Status(ctx, endpoints[2])
	fmt.Println("etcd Status took", time.Since(stp20))

	stp22 := time.Now()
	cli2.Get(ctx, "sample_key_2")
	fmt.Println("etcd Get took", time.Since(stp22))

	stp21 := time.Now()
	cli2.Put(ctx, "sample_key_2", "sample_value_2")
	fmt.Println("etcd Put took", time.Since(stp21))
	fmt.Println()

	// Check latency
	sizeOfEndpoints := len(endpoints)
	latencyList := make([]time.Duration, sizeOfEndpoints)

	fmt.Printf("Endpoints: %v\n", endpoints)

	for index, endpoint := range endpoints {
		fmt.Printf("Endpoint: %s\n", endpoint)

		client := resty.New()

		resp, err := client.R().
			Get(fmt.Sprintf("http://%s/health", endpoint))

		// Output print
		fmt.Printf("\nError: %v\n", err)
		fmt.Printf("Time: %v\n", resp.Time())
		fmt.Printf("Body: %v\n", resp)

		latencyList[index] = resp.Time()
		fmt.Println(latencyList[index])
	}

	// Sort by latency
	for i := 0; i < sizeOfEndpoints; i++ {
		for j := sizeOfEndpoints - 1; j >= i+1; j-- {
			if latencyList[j] < latencyList[j-1] {
				latencyList[j], latencyList[j-1] = latencyList[j-1], latencyList[j]
				endpoints[j], endpoints[j-1] = endpoints[j-1], endpoints[j]
			}
		}
	}

	fmt.Printf("Endpoints: %v\n", endpoints)
	fmt.Printf("Latency: %v\n", latencyList)

	// // etcd Section
	// etcdClient, err := clientv3.New(clientv3.Config{
	// 	Endpoints:   endpoints,
	// 	DialTimeout: 5 * time.Second,
	// })

	// if err != nil {
	// 	CBLogger.Fatal(err)
	// }

	// defer func() {
	// 	errClose := etcdClient.Close()
	// 	if errClose != nil {
	// 		CBLogger.Fatal("Can't close the etcd client", errClose)
	// 	}
	// }()

}
