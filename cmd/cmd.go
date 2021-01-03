package cmd

import (
	"os"

	ishell "gosploit/ishell"
	. "gosploit/termcolor"
)

var Shell *ishell.Shell = ishell.New()
var Available_Modules = map[string]string{}

func ResetPrompt() {
	Shell.Printf("\n%s", StackedPrompt())
}

func Init() {
	// shell := ishell.New()
	// Load(shell)
	InitDNS()
	InitTest()

	// display info.
	Shell.Println("Sample Interactive Shell")
	Shell.SetPrompt(StackedPrompt())
	//Consider the unicode characters supported by the users font
	//shell.SetMultiChoicePrompt(" >>"," - ")
	//shell.SetChecklistOptions("[ ] ","[X] ")

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Core Commands <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	// Custom Help.
	Shell.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"?", "usage"},
		Help:    "Help for Core Functions",
		Func: func(c *ishell.Context) {
			c.Println(TagMsg("Core", "HELP", "\n"))
			c.Println("         help      Show this crap")
			c.Println("         clear     Clear the screen")
			c.Println("         exit      Exit GoSploit")
			c.Println("         modules   List available modules")
			// c.Println("         ")
			// c.Println("         ")
			c.Println("")

		},
	})

	// Custom Help.
	Shell.AddCmd(&ishell.Cmd{
		Name: "modules",
		// Aliases: []string{"?", "usage"},
		Help: "List available modules",
		Func: func(c *ishell.Context) {
			c.Println(TagMsg("Core", "Available modules", ""))
			c.Println("")
			for moda, _ := range Available_Modules {
				c.Println(ColorModName(moda) + "          " + Available_Modules[moda])
			}
			c.Println("")
		},
	})

	//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	// when started with "exit" as first argument, assume non-interactive execution
	if len(os.Args) > 1 && os.Args[1] == "exit" {
		Shell.Process(os.Args[2:]...)
	} else {
		// start shell
		Shell.Run()
		// teardown
		Shell.Close()
	}
}
