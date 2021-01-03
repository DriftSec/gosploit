package termcolor

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// foreground colors
var forecolors = map[string]color.Attribute{
	"red":     color.FgRed,
	"black":   color.FgBlack,
	"blue":    color.FgBlue,
	"white":   color.FgWhite,
	"cyan":    color.FgCyan,
	"yellow":  color.FgYellow,
	"magenta": color.FgMagenta,
	"green":   color.FgGreen,
}

// bacground colors
var backcolors = map[string]color.Attribute{
	"red":     color.BgRed,
	"black":   color.BgBlack,
	"blue":    color.BgBlue,
	"white":   color.BgWhite,
	"cyan":    color.BgCyan,
	"yellow":  color.BgYellow,
	"magenta": color.BgMagenta,
	"green":   color.BgGreen,
}

// color schemes for modules
var ModColorSchemes = map[string]string{
	"core":  "green",
	"other": "cyan",
}

// color schemes for msg types
var PriorityColorSchemes = map[string]string{
	"error":   "red",
	"warning": "yellow",
	"success": "green",
	"help":    "yellow",
	"info":    "white",
	"other":   "white",
}

// maps certain text colors to the background color so it contrasts
var textcolorscemes = map[string]string{
	"red":     "black",
	"black":   "black",
	"blue":    "white", // blue background white text
	"white":   "black",
	"cyan":    "black",
	"yellow":  "black",
	"magenta": "white",
	"green":   "White",
}

func StackedPrompt() string {
	prompt := ""
	promptcore := color.New(color.FgWhite).Add(color.BgGreen).SprintFunc()
	prompt = prompt + promptcore(" GoSploit ") + " » "
	return prompt
}

func Tag(text string, bcolor string) string {

	tagf := color.New(forecolors[textcolorscemes[bcolor]]).Add(backcolors[bcolor]).SprintFunc()
	return tagf(text)
}

func ColorModName(mod string) string {
	modcolorscheme := ModColorSchemes[strings.ToLower(mod)]
	if modcolorscheme == "" {
		modcolorscheme = ModColorSchemes["other"]
	}

	colf := color.New(forecolors[modcolorscheme]).SprintFunc()
	return colf(mod)

}

func ColorText(text string, txtcolor string) string {
	colf := color.New(forecolors[txtcolor]).SprintFunc()
	return colf(text)

}

// returns the color coded msg so it can be used elsewhere.
func TagMsg(msgmod string, msgtype string, msg string) string {

	msgmodbase := strings.Split(msgmod, ".")[0]
	modcolorscheme := ModColorSchemes[strings.ToLower(msgmodbase)]
	if modcolorscheme == "" {
		modcolorscheme = ModColorSchemes["other"]
	}
	modtextcolorscheme := textcolorscemes[modcolorscheme]

	typecolorscheme := PriorityColorSchemes[strings.ToLower(msgtype)]
	if typecolorscheme == "" {
		typecolorscheme = PriorityColorSchemes["other"]
	}
	typetextcolorscheme := textcolorscemes[typecolorscheme]

	modarrowcolor := color.New(forecolors["black"]).Add(backcolors[modcolorscheme]).SprintFunc()
	modcolor := color.New(forecolors[modtextcolorscheme]).Add(backcolors[modcolorscheme]).SprintFunc()
	joinarrowcolor := color.New(forecolors[modcolorscheme]).Add(backcolors[typecolorscheme]).SprintFunc()
	typecolor := color.New(forecolors[typetextcolorscheme]).Add(backcolors[typecolorscheme]).SprintFunc()
	endarrotcolor := color.New(forecolors[typecolorscheme]).Add(backcolors["black"]).SprintFunc()
	msgmod = strings.ToUpper(msgmod)
	msgtype = strings.ToUpper(msgtype)
	return fmt.Sprintf("\n%s%s%s%s%s %s", modarrowcolor("\ue0b0"), modcolor(msgmod), joinarrowcolor("\ue0b0"), typecolor(msgtype), endarrotcolor("\ue0b0"), msg)
}
