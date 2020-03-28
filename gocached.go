// Package main implements an LRU cache daemon in golang.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"regexp"
	"strings"

	log "github.com/golang/glog"
	gcache "github.com/morrowc/GoCache/gCache"
)

var (
	portArg     = flag.String("port", "9000", "GoCached Server Port")
	capacityArg = flag.Int("capacity", 5, "Capacity of the cache")
)

// processCommand will process the command and execute relevant LRU methods
func processCommand(c string, gc *gcache.Cache) (string, error) {
	var result string

	pattern := ""
	if strings.Index(strings.ToLower(c), "get") > -1 {
		pattern = `(?i)GET\s+(\w+)`
	} else if strings.Index(strings.ToLower(c), "set") > -1 {
		pattern = `(?i)SET\s(\w+)\s(\w+)`
	}
	var re = regexp.MustCompile(pattern)
	// TODO(kadnan): Based on results call Cache Methods

	matches := re.FindAllStringSubmatch(c, -1)

	var err error
	switch {
	case len(matches[0]) == 2:
		key := matches[0][1]
		result, err = gc.get(key)
		if err != nil {
			return "", err
		}
	case len(matches[0]) == 3:
		key := matches[0][1]
		val := matches[0][2]
		result = gc.set(key, val)
	}
	return result, nil
}

func validateCommand(cmd string) bool {
	var msg bool
	//Split the command to make sure it is not more than 2 or 3
	cmdArray := strings.Split(cmd, " ")
	// cmd must be at least 2 elements (get/set and what to get/set) and not more than 3 elements.
	if len(cmdArray) < 2 || len(cmdArray) > 3 {
		return false
	}
	gs := strings.ToLower(cmdArray[0])

	// default to return false, only test for truth and return true.
	switch {
	case len(cmdArray) == 3 && gs == "set":
		return true
	case len(cmdArray) == 2 && gs == "get":
		return true
	}
	return false
}

func handleConnection(c net.Conn, gc *Cache) {
	defer c.Close()
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Fatalf("error building string reader: %v", err)
		}

		temp := strings.TrimSpace(string(netData))
		result := []byte("Invalid Command\n")
		if validateCommand(temp) {
			outcome, err := processCommand(temp, gc)
			if err != nil {
				log.Fatalf("error processing cache command: %v", err)
			}
			result = []byte(outcome + "\n")
		}
		if temp == "\\q" {
			return
		}
		c.Write(result)
	}
}

func main() {
	flag.Parse()

	PORT := ":" + *portArg
	fmt.Printf("Launching server at port %s with the capacity %d \n", *portArg, *capacityArg)
	gCache := gcache.New(*capacityArg)

	// Listen on all IPv4 interfaces.
	connection, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()

	for {
		c, err := connection.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c, &gCache)
	}
}
