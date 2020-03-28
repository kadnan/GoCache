package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"regexp"
	"strings"
)

/*
	This function will process the command and execute relevant LRU methods
*/
func processCommand(c string, gc Cache) string {
	var result string
	var er error
	var pattern = ""

	if strings.Index(strings.ToLower(c), "get") > -1 {
		pattern = `(?i)GET\s+(\w+)`
	} else if strings.Index(strings.ToLower(c), "set") > -1 {
		pattern = `(?i)SET\s(\w+)\s(\w+)`
	}
	// TODO: Based on results call Cache Methods

	var re = regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(c, -1)

	if len(matches) > 0 {
		// For GET command
		if len(matches[0]) == 2 {
			key := matches[0][1]
			result, er = gc.get(key)
			if er != nil {
				result = er.Error()
			}

		} else if len(matches[0]) == 3 {
			key := matches[0][1]
			val := matches[0][2]
			result = gc.set(key, val)
			//gc.print()
		}
	}
	return result
}

func validateCommand(cmd string) bool {
	var msg bool
	//Split the command to make sure it is not more than 2 or 3
	cmdArray := strings.Split(cmd, " ")

	if len(cmdArray) > 3 {
		msg = false

	} else if (len(cmdArray) == 2) && (strings.TrimSpace(strings.ToLower(cmdArray[0])) != "get") {
		msg = false

	} else if len(cmdArray) == 3 && strings.ToLower(cmdArray[0]) != "set" {
		msg = false

	} else if len(cmdArray) == 3 && strings.ToLower(cmdArray[0]) == "set" {
		msg = true

	} else if len(cmdArray) == 2 && strings.ToLower(cmdArray[0]) == "get" {
		msg = true

	}

	return msg
}

func handleConnection(c net.Conn, gc Cache) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if validateCommand(temp) {
			outcome := processCommand(temp, gc)
			c.Write([]byte(string(outcome + "\n")))

		} else {
			c.Write([]byte(string("Invalid Command\n")))
		}
		//fmt.Println(z)

		if temp == "\\q" {
			break
		}
	}
	c.Close()
}

func main() {

	portArg := flag.String("port", "9000", "GoCached Server Port")
	capacityArg := flag.Int("capacity", 5, "Capacity of the cache")
	flag.Parse()

	PORT := ":" + *portArg
	fmt.Printf("Launching server at port %s with the capacity %d \n", *portArg, *capacityArg)
	//os.Exit(0)
	var gCache = Cache{*capacityArg}

	// listen on all interfaces
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
		go handleConnection(c, gCache)
	}
}
