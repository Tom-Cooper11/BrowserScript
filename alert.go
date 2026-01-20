package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"github.com/fatih/color"
)

import "alert/createAlert"

type Argue struct {
	Auth       string
	Method    string
}

func main() {
	c := color.New(color.FgBlue)
	c.Println("    ALERT CREATE TOOLS")
	c1 := color.New(color.FgRed)
	c1.Println("         V1.0.0")

	var argue Argue
	flag.StringVar(&argue.Auth, "a", "", "cookie authorization.")
	flag.StringVar(&argue.Method, "m", "", "which function to use.")
	flag.Parse()

	if argue.Auth == "" || argue.Method == "" {
		c2 := color.New(color.FgRed)
		c2.Println("Please enter the correct parameters.")
		c2.Println("Example: alert -a 'cookie authorization' -m 'a'")
		return
	}

	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}
	execDir := filepath.Dir(execPath)

	switch argue.Method {
		case "a":
		createAlert.Alert_Follow(argue.Auth, execDir)
		case "b":
		createAlert.Alert_Judge(argue.Auth, execDir)
		case "c":
		createAlert.Alert_Task(argue.Auth, execDir)
		case "d":
		createAlert.Alert_Save(argue.Auth, execDir)
		default:
		c3 := color.New(color.FgRed)
		c3.Println("Please enter the correct method parameter.")
		c3.Println("Method options: a, b, c, d")
	}
}