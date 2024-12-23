package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/skywind3000/tcpdup/forward"
)

func start(listen string, target string, input string, output string) int {
	if listen == "" || target == "" {
		return -1
	}
	logger := log.Default()
	logger.Printf("Service starting:\n")
	logger.Printf("config: listen %s\n", listen)
	logger.Printf("config: target %s\n", target)
	logger.Printf("config: input %s\n", input)
	logger.Printf("config: output %s\n", output)
	service := forward.NewTcpForward()
	service.SetLogger(logger)
	service.SetInput(input)
	service.SetOutput(output)
	listenAddr := forward.AddressResolve(listen)
	if listenAddr == nil {
		logger.Printf("invalid listen address: %s\n", listen)
		return -1
	}
	targetAddr := forward.AddressResolve(target)
	if targetAddr == nil {
		logger.Printf("invalid target address: %s\n", target)
		return -1
	}
	service.Open(listenAddr, targetAddr)
	service.Wait()
	return 0
}

func main() {
	listen := flag.String("listen", "", "local address, eg: 0.0.0.0:8080")
	target := flag.String("target", "", "destination address, eg: 8.8.8.8:8080")
	input := flag.String("input", "", "input duplication address, eg: 127.0.0.1:8081")
	output := flag.String("output", "", "output duplication address, eg: 127.0.0.1:8082")
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Printf("Usage of %s:\n", os.Args[0])
		order := []string{"listen", "target", "input", "output"}
		for _, name := range order {
			flag := flagSet.Lookup(name)
			fmt.Printf("-%s\n", flag.Name)
			fmt.Printf("  %s\n", flag.Usage)
		}
	}
	if false {
		start("0.0.0.0:8080", "127.0.0.1:8000", "127.0.0.1:8081", "127.0.0.1:8082")
		return
	}
	flag.Parse()
	if listen == nil || target == nil {
		flag.Usage()
		return
	}
	if *listen == "" || *target == "" {
		flag.Usage()
		return
	}
	i := ""
	o := ""
	if input != nil {
		i = *input
	}
	if output != nil {
		o = *output
	}
	start(*listen, *target, i, o)
}
