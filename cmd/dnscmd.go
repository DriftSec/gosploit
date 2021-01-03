package cmd

import (
	b64 "encoding/base64"
	"fmt"
	ishell "gosploit/ishell"
	. "gosploit/termcolor"
	"regexp"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

var dnsNet string = "udp"
var dnsPort int = 54
var dnsMod string = "dns"

var records = map[string]string{
	"test.service.": "192.168.0.2",
}

func InitDNS() {
	// name to use fore this module
	mName := dnsMod

	// add a color scheme for this module.
	ModColorSchemes[mName] = "blue"

	// add name and description to available modules map
	Available_Modules[mName] = " DNS server module (see " + mName + ".help)"

	// new message types can be added and color coded.
	PriorityColorSchemes["exfil"] = "magenta"

	// multiple choice
	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "net",
		Help: "protocol for the server to use (interactive)",
		Func: func(c *ishell.Context) {
			choice := c.MultiChoice([]string{"udp", "tcp"}, "")
			if choice == 1 {
				dnsNet = "tcp"
			} else {
				dnsNet = "udp"
			}
			c.Println("\nUsing ", dnsNet, "\n")
		},
	})

	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "port",
		Help: "Set/Get the port for the server to use (dns.port [int])", // args need to be specified in help using [] so autohelp can spot them.
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				p, err := strconv.Atoi(c.Args[0])
				if err == nil {
					dnsPort = p
					c.Println("\nDNS Port = "+strconv.Itoa(dnsPort), "\n")
				} else {
					c.Println(TagMsg(mName, "ERROR", "Not a valid port !!!"))
					fmt.Println(err)
				}
			} else {
				c.Println("\nDNS Port = "+strconv.Itoa(dnsPort), "\n")
			}

		},
	})

	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "start",
		Help: "Start the DNS server",
		Func: func(c *ishell.Context) {
			Shell.Printf("%s%d\n\n", TagMsg(dnsMod, "INFO", "Starting the DNS server on port "), dnsPort)
			go start()

		},
	})

	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "stop",
		Help: "Stop the DNS server",
		Func: func(c *ishell.Context) {

			Shell.Printf("%s%d\n\n", TagMsg(dnsMod, "INFO", "Stoping the DNS server."))
		},
	})

	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "records",
		Help: "Edit DNS records (interactive)",
		Func: func(c *ishell.Context) {

			c.Println("do some dns stuff")
		},
	})

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Auto generated help for each command in the module.
	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "help",
		Help: "Help for the DNS module",
		Func: func(c *ishell.Context) {
			c.Println(TagMsg(mName, "HELP", "DNS Commands:\n"))
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

func parseQuery(m *dns.Msg) {
	// fmt.Println(m.Question[0].Name)
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:

			ResetPrompt() // have to reset the prompt because its async (need to find a better way.)

			ip := records[q.Name]
			if ip != "" {
				Shell.Printf("\n%s%s%s\n", TagMsg(dnsMod, "QUERY", "Query for "), q.Name, " > "+ip)
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			} else {
				subs := strings.Split(q.Name, ".")
				sub := subs[0]
				Shell.Println(sub)
				sDec, err := b64.StdEncoding.DecodeString(sub)
				if err == nil {
					Shell.Printf("\n%s%s\n", TagMsg(dnsMod, "EXFIL", ""), sDec)
					ResetPrompt()
				} else {
					Shell.Printf("\n%s%s%s\n", TagMsg(dnsMod, "QUERY", "Query for "), q.Name, " > FORWARDING")
					ResetPrompt()
				}

			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	// Shell.Println("a> ", r.Question[0].Name)
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}
	w.WriteMsg(m)
}

func start() {
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	port := dnsPort
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: dnsNet}
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		Shell.Printf("\n%s%s\n\n", TagMsg(dnsMod, "ERROR", "Error starting the server: "), err.Error())
		ResetPrompt()
		// log.Fatalf("Failed to start server: %s\n ", err.Error())
	}

}
