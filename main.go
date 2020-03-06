package main

import (
	"fmt"
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	clear map[string]func() 		// map for storing screen clear funcs
	input = bufio.NewScanner(os.Stdin)

	greetingTxt = `
 ****************************************
 ************* NoteMap v1.0 *************
 ****************************************
`
	menuTxt = `
 Main Menu
	
 Select an Option
 [1] New Map
 [2] Open Map
 [3] Quit
	
 |--> `
)

func init() {
	clear = make(map[string]func())
	
    clear["linux"] = func() { 
        cmd := exec.Command("clear") //Linux
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
    clear["windows"] = func() {
        cmd := exec.Command("cmd", "/c", "cls") //Windows
        cmd.Stdout = os.Stdout
        cmd.Run()
	}

	LoadMaps()
	LoadSettings()
}

func main() {
	fmt.Println(greetingTxt)

	for {
		fmt.Print(menuTxt)

		choice, err := GetInput()
		if err != nil {
			fmt.Printf("\n\terror @menu - \n", os.Stderr, err)
		}

		if choice == "3" {
			fmt.Println("\n\tQuiting... ")
			for done := false; !done; {
				fmt.Print("Save all Changes? [y/n]: ")
	
				choice, err = GetInput()
				choice = stirngs.ToUpper(choice)
				if err != nil {
					fmt.Printf("\nFailed to read save option selection\n\t%v reading standard input: %v", os.Stderr, err)
				}
	
				if choice == "Y" {
					Save()
					done = true
				} else if choice == "N"{
					done = true
				} else {
					fmt.Printf("\ninput not a valid option - %s\nenter 'Y/y' to save or 'N/n' to exit without saving changes\n\n", choice)
				}
			}

			break
		} else {
			switch choice {
			case "2":
				if err := DisplayMaps(); err != nil {
					fmt.Println("Failed to display maps: ", err.Error())
				}
			case "1":
				if nm, err := NewNoteMap(input); err != nil {
					fmt.Printf("Failed to create new NoteMap: ", err)
					continue
				}
				nm.Open()
			default:
				fmt.Println("\tInput not recognized as valid option")
				time.Sleep(2 * time.Second)
				ClearScreen()
			}
		}
	}
}

func ClearScreen() {
    value, ok := clear[runtime.GOOS]
    if ok {
		value()
    } else {
		LineSkip()
    }
}
func LineSkip() {
	if v, err := strconv.Atoi(configuration.UserSettings.Skips); err == nil {
		for i := 0; i < v; i++ {
			fmt.Println()
		}
	} else if v, err := strconv.Atoi(configuration.DefaultSettings.Skips); err == nil{
		for i := 0; i < v; i++ {
			fmt.Println()
		}
	}
}

func DisplayMaps() error {
	if len(savedMaps) == 0 {
		return errors.New("No existing NoteMap data found")
	} else {
		ClearScreen()

		fmt.Println("\n Select a NoteMap to open\n")
	
		for i, v := range savedMaps {
			fmt.Printf("\t[%s] %s\tLast Updated: %s\tCreated: %s\n\t%s\n\n", i, v.Name, v.LastUpdated, v.Creation, v.Description)
		}

		var choice string
		var err error

		for {
			choice, err = GetInput()

			if err != nil {
				return err
			}

			if val, err := strconv.Atoi(choice); err == nil {
				for i, v := range savedMaps {
					if val == i {
						v.Open()
						break
					}
				}
			} else {
				fmt.Printf("input not recognized - %s\n", choice)
			}
		}

		return nil
	}
}

func Save() {
	fmt.Println("Saving...")
}

type Note struct {
	Subject string
	Content string
	Relations []*Note
}
type NoteMap struct {
	Name string
	Creation time.Time
	LastUpdated time.Time
	Description string
	Root *Note
}

func (nm *NoteMap) Open() {
	active := nm.Root
	navCache = append(navCache, active)

	for {
		ClearScreen()
		/* Display Global Commands
		*	Return To Main Menu
		*	Return to previous Note
				if len(navCache) > 1 {}
		*/

		// Display NoteMap Title
		fmt.Printf("\n\n%s\n%s\n\n", nm.Name, nm.Description)
	
		// Display Nav Options
		for i, v := range active.Relations {
			fmt.Printf("\t[%s]", v.Subject)

			if i % 3 == 0 {
				fmt.Println()
			}
		}

		fmt.Print("\ninput command\n|--> ")
		choice, err := GetInput()

		switch  {

		}
	}
}

func (n *Note) Open() {

	pr, pw := io.Pipe()
	defer pw.Close()

	// Display Note Title
	fmt.Println(n.Subject)
	// Display Note Content
	for i := 0; i < len(n.Content); i++ {
		if i % 16 == 0 {
			fmt.Printf("\n%s", n.Content[i])
		} else {
			fmt.Print(n.Content[i])
		}
	}

	// tell the command to write to our pipe
	pw.Write([]byte(n.Content))

	go func() {
		defer pr.Close()
		// copy the data written to the PipeReader via the cmd to stdout
		if _, err := io.Copy(os.Stdout, pr); err != nil {
			log.Fatal(err)
		}
	}()
}

var (
	savedMaps 	[]*NoteMap 			// notemaps found saved on disk

	navCache 	[]*Note				// cache for navigating back to Root (stores Note.Name)
	saveCache 	map[string]Note		// cache for changes to Notes (stores Note.Name ==> Note)
)



func NewNoteMap(s *bufio.Scanner) (*NoteMap, error) {
	ClearScreen()

	nm := &NoteMap{}

	// Get Name for map
	fmt.Print(" Creating new NoteMap\n Enter Name: ")

	name, err := GetInput()
	if err != nil {
		return nil, err
	}

	// Get Description for map
	fmt.Print("\nEnter Description: ")

	desc, err := GetInput()
	if err != nil {
		return nil, err
	}

	nm.Name = name
	nm.Description = desc

	nm.Creation = time.Now()
	nm.LastUpdated = nm.Creation
	nm.Root = new(Note)

	return nm, nil
}

func LoadMaps() {
	
}

var (
	configSrc = "config.json"
	configuration Config
)

type Config struct {
	Version 		string `json:version""`

	UserSettings	Settings
	DefaultSettings Settings
}

type Settings struct {
	Skips 		string `json:"skips`
	SaveDir 	string `json:"saveDir"`
	SaveFile 	string `json:"saveFile"`
}

func LoadSettings() {
	file, err := ioutil.ReadFile(configSrc)
	if err != nil {
		log.Fatal("Initizilaion Error - Failed to read config file : ", err)
	}

	err = json.Unmarshal(file, &configuration)
	if err != nil {
		fmt.Println("Initialization Error - Failed to extract content from JSON:", err)
	}
}

func GetInput() (string, error) {
	scanOk := input.scan()

	if !scanOk {
		return "", fmt.Errorf("failure reading from standard input: %v", os.Stderr, err)
	} else {
		return input.Text(), nil
	}
}