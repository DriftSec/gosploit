package cmd

import (
	"fmt"
	ishell "gosploit/ishell"
	. "gosploit/termcolor"
	"regexp"
	"strings"
)

// test.go is a template for new modules.

// InitNAME() must be called in cmd.go to initialize
func InitTest() {
	// name to use fore this module
	mName := "test"

	// add a color scheme for this module.
	ModColorSchemes[mName] = "red"

	// add name and description to available modules map
	Available_Modules[mName] = " Test module (see " + mName + ".help)"

	// multiple choice
	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "test1",
		Help: "test1 blah blah help",
		Func: func(c *ishell.Context) {
			c.Println("\ntest1\n")
		},
	})

	Shell.AddCmd(&ishell.Cmd{
		Name: mName + "." + "var",
		Help: "Set/Get a var (test.var [int])", // args need to be specified in help using [] so autohelp can spot them.
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				c.Println("args =" + c.Args[0])
			} else {
				c.Println("\nvar Port = \n")
			}

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
