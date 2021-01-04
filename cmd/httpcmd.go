package cmd

import (
	"bufio"
	"fmt"
	ishell "gosploit/ishell"
	. "gosploit/termcolor"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tomnomnom/rawhttp"
)

var httpUrl string

// var httpUserAgent string
// var httpHeaders string
// var httpCookie string
// var httpProxy string
// var httpMethod string

func InitHTTP() {
	// name to use fore this module
	mName := "http"

	// add a color scheme for this module.
	ModColorSchemes[mName] = "yellow"

	// add name and description to available modules map
	Available_Modules[mName] = " http module (see " + mName + ".help)"

	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "url",
		Help: "Set/Get URL (http.url [URL])", // args need to be specified in help using [] so autohelp can spot them.
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				httpUrl = c.Args[0]
			} else {
				c.Println("\nhttp.url = " + httpUrl + "\n")
			}

		},
	})

	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "test",
		Help: "Set/Get URL (http.url [URL])", // args need to be specified in help using [] so autohelp can spot them.
		Func: func(c *ishell.Context) {
			a := parseFile("/tmp/test.req")
			fmt.Println(a)

		},
	})

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Auto generated help for each command in the module.
	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "help",
		Help: "Help for the DNS module",
		Func: func(c *ishell.Context) {
			c.Println(TagMsg(mName, "HELP", "DNS Commands:"))
			longestline := 0
			for _, b := range Shell.RootCmd().Children() { // gotta be a better way.  first loop to find the longest line for indent
				cmdA := strings.Split(b.Name, ".")
				if cmdA[0] == mName {
					argStr := ""
					re := regexp.MustCompile("(\\[.*\\])")
					match := re.FindStringSubmatch(b.Help)
					if match != nil {
						argStr = match[1] //strings.Join(match[:], " ")
					} else {
						argStr = ""
					}
					left := fmt.Sprintf("%s %s", ColorModName(mName)+"."+ColorText(strings.Join(cmdA[1:], "."), "yellow"), argStr)
					ll := len(left)
					if ll > longestline {
						longestline = ll
					}
				}
			}

			for _, b := range Shell.RootCmd().Children() {
				// auto list help so we dont have to create it.
				cmdA := strings.Split(b.Name, ".")
				if cmdA[0] == mName {
					argStr := ""
					re := regexp.MustCompile("(\\[.*\\])")
					match := re.FindStringSubmatch(b.Help)
					if match != nil {
						argStr = match[1]
					} else {
						argStr = ""
					}
					left := fmt.Sprintf("%s %s", ColorModName(mName)+"."+ColorText(strings.Join(cmdA[1:], "."), "yellow"), argStr)
					leftlen := len(left)
					indent := longestline + 5 - leftlen
					c.Println("         ", left, strings.Repeat(" ", indent), b.Help)
				}

			}
			c.Println("")
		},
	})
}

func parseFile(path string) rawhttp.Request {
	req := rawhttp.Request{}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ln := 0
	for scanner.Scan() {
		// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< may be able to split on double newline  headers[0] \n\n body[1:]  read all with (1:) for multipart, may need to join with \n\n
		ln++
		if ln == 1 { // first line
			ar := strings.Split(scanner.Text(), " ")
			req.Method = ar[0]

			pq := strings.Split(ar[1], "?")
			req.Path = pq[0]
			req.Query = pq[1]
			req.Proto = ar[2]

		}
		curln := scanner.Text()
		if strings.Contains(curln, "Host:") { //<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< may be able to just append each line to []Headers and let it deal with it.
			req.Hostname = strings.Split(curln, "Host: ")[1]
			fmt.Println(req.Hostname)
		}
		// fmt.Println(scanner.Text())
		//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< need to prompt for port/scheme
	}
	// fmt.Println(req.String())
	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }
	//
	// req.Port = "80"
	// // req.Fragment = "anchor"
	// req.EOL = "\r\n"
	// req.AutoSetContentLength()
	//
	// resp, err := rawhttp.Do(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf("< %s\n", resp.StatusLine())
	// for _, h := range resp.Headers() {
	// 	fmt.Printf("< %s\n", h)
	// }
	//
	// fmt.Printf("\n%s\n", resp.Body())
	return req
}

func testreq() {
	req, err := rawhttp.FromURL("POST", "https://httpbin.org")
	if err != nil {
		log.Fatal(err)
	}

	// automatically set the host header
	req.AutoSetHost()

	req.Method = "PUT"
	req.Hostname = "httpbin.org"
	req.Port = "443"
	req.Path = "/anything"
	req.Query = "one=1&two=2"
	req.Fragment = "anchor"
	req.Proto = "HTTP/1.1"
	req.EOL = "\r\n"

	req.AddHeader("Content-Type: application/x-www-form-urlencoded")

	req.Body = "username=AzureDiamond&password=hunter2"

	// automatically set the Content-Length header
	req.AutoSetContentLength()

	fmt.Printf("%s\n\n", req.String())

	resp, err := rawhttp.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("< %s\n", resp.StatusLine())
	for _, h := range resp.Headers() {
		fmt.Printf("< %s\n", h)
	}

	fmt.Printf("\n%s\n", resp.Body())
}
