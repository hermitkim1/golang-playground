package main

import (
	"context"
	"errors"
	"fmt"
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
	fmt.Println("Start......... init()")
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

	fmt.Println("End......... init()")
	fmt.Println("")
}

func main() {

	endpoints := config.ETCD.Endpoints
	ctx := context.TODO()

	// Pre-test
	fmt.Println("############################")
	fmt.Println("## Pre-test               ##")
	fmt.Println("############################")

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		CBLogger.Fatal(err)
	}

	defer func() {
		errClose := etcdClient.Close()
		if errClose != nil {
			CBLogger.Fatal("Can't close the etcd client", errClose)
		}
	}()

	for i := 1; i <= 3; i++ {
		t1 := time.Now()
		resp, _ := etcdClient.Get(ctx, "sample_key_2")
		fmt.Println("MemberId: ", resp.Header.MemberId)
		fmt.Println("etcd Get took", time.Since(t1))

		t2 := time.Now()
		etcdClient.Put(ctx, "sample_key_2", "sample_value_2")
		fmt.Println("etcd Put took", time.Since(t2))
		fmt.Println()

		time.Sleep(500 * time.Millisecond)
	}

	// Performance evaluation
	fmt.Println("############################")
	fmt.Println("## Performance evaluation ##")
	fmt.Println("############################")

	// Target 0
	fmt.Println(endpoints[0])
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
	resp, _ := cli0.Get(ctx, "sample_key_0")
	fmt.Println("MemberId: ", resp.Header.MemberId)
	fmt.Println("etcd Get took", time.Since(stp02))

	stp01 := time.Now()
	cli0.Put(ctx, "sample_key_0", "sample_value_0")
	fmt.Println("etcd Put took", time.Since(stp01))

	fmt.Println()

	// Target 1
	fmt.Println(endpoints[1])
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
	resp, _ = cli1.Get(ctx, "sample_key_1")
	fmt.Println("MemberId: ", resp.Header.MemberId)
	fmt.Println("etcd Get took", time.Since(stp12))

	stp11 := time.Now()
	cli1.Put(ctx, "sample_key_1", "sample_value_1")
	fmt.Println("etcd Put took", time.Since(stp11))

	fmt.Println()

	// Target 2
	fmt.Println(endpoints[2])
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
	resp, _ = cli2.Get(ctx, "sample_key_2")
	fmt.Println("MemberId: ", resp.Header.MemberId)
	fmt.Println("etcd Get took", time.Since(stp22))

	stp21 := time.Now()
	cli2.Put(ctx, "sample_key_2", "sample_value_2")
	fmt.Println("etcd Put took", time.Since(stp21))
	fmt.Println()

	// Check latency
	fmt.Println("############################")
	fmt.Println("## Latency checking       ##")
	fmt.Println("############################")
	sizeOfEndpoints := len(endpoints)
	latencyList := make([]time.Duration, sizeOfEndpoints)

	fmt.Printf("Endpoints: %v\n", endpoints)

	for index, endpoint := range endpoints {
		fmt.Printf("Endpoint: %s\n", endpoint)

		client := resty.New()

		trial := 50
		timeList := make([]int64, trial)

		for i := 0; i < trial; i++ {
			resp, _ := client.R().
				Get(fmt.Sprintf("http://%s/health", endpoint))

			// Output print
			// fmt.Printf("Error: %v\n", err)
			// fmt.Printf("Response time (%d): %v\n", i, resp.Time())
			// fmt.Printf("Body: %v\n", resp)
			timeList[i] = int64(resp.Time() / time.Millisecond)
		}

		var total int64 = 0
		for _, value := range timeList {
			total += value
		}

		avgLatency := time.Duration(total/int64(trial)) * time.Millisecond

		latencyList[index] = avgLatency
		fmt.Printf("Average response time: %v\n", latencyList[index])
		fmt.Println()
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

	// Post-test
	fmt.Println("############################")
	fmt.Println("## Post-test              ##")
	fmt.Println("############################")
	etcdClient2, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		CBLogger.Fatal(err)
	}

	defer func() {
		errClose := etcdClient2.Close()
		if errClose != nil {
			CBLogger.Fatal("Can't close the etcd client", errClose)
		}
	}()

	for i := 1; i <= 3; i++ {
		t1 := time.Now()
		resp, _ := etcdClient2.Get(ctx, "sample_key_2")
		fmt.Println("MemberId: ", resp.Header.MemberId)
		fmt.Println("etcd Get took", time.Since(t1))

		t2 := time.Now()
		etcdClient2.Put(ctx, "sample_key_2", "sample_value_2")
		fmt.Println("etcd Put took", time.Since(t2))
		fmt.Println()

		time.Sleep(500 * time.Millisecond)
	}

}
