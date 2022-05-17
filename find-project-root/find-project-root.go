package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// CBLogger represents a logger to show execution processes according to the logging level.
// var CBLogger *logrus.Logger

func init() {
	// fmt.Println("Start......... init() of networking-rule.go")
	// ex, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }
	// exePath := filepath.Dir(ex)
	// fmt.Printf("exePath: %v\n", exePath)

	// // Load cb-log config from the current directory (usually for the production)
	// logConfPath := filepath.Join(exePath, "config", "log_conf.yaml")
	// fmt.Printf("logConfPath: %v\n", logConfPath)
	// if !file.Exists(logConfPath) {
	// 	fmt.Printf("not exist - %v\n", logConfPath)
	// 	// Load cb-log config from the project directory (usually for development)
	// 	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	// 	fmt.Printf("projectRoot: %v\n", string(path))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	projectPath := strings.TrimSpace(string(path))
	// 	logConfPath = filepath.Join(projectPath, "poc-cb-net", "config", "log_conf.yaml")
	// }
	// CBLogger = cblog.GetLoggerWithConfigPath("cb-network", logConfPath)
	// CBLogger.Debugf("Load %v", logConfPath)
	// fmt.Println("End......... init() of networking-rule.go")

}

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func PrintMyPath() {
	fmt.Println(basepath)
}

func main() {

	PrintMyPath()

}
