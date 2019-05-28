package imnu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/colt3k/nglog/ng"
)

type Menu struct {
	Id          int
	Name        string
	Description string
	Enable      interface{}
	Task        interface{}
}

type InteractiveMenu struct {
	menus             []Menu
	displaySelections interface{}
	lastErr           bool
	lastMsg           bool
	msg               string
}

/*
Flow
	1. Show menu with options to enter
	2. User pushes # and enter
	3. System performs and outputs data on screen with option to choose what to do next
*/
func New(displaySelections interface{}, menulayout []Menu) *InteractiveMenu {
	t := new(InteractiveMenu)
	t.menus = menulayout
	t.displaySelections = displaySelections

	return t
}

func (i *InteractiveMenu) StartMenu() {
	for {
		//log.Printf("%s", log.ClearTerminalSequence) // Clear Terminal
		log.Printf(strings.Repeat("\n", 20))
		err := i.menu()
		if err == -1 {
			i.lastErr = true
		} else {
			i.lastErr = false
		}
	}
}

func (i *InteractiveMenu) SetMsg(msg string) {
	i.msg = msg
	i.lastMsg = true
}
func (i *InteractiveMenu) menu() int {

	if len(i.menus) <= 0 {
		log.Printf("%s", log.ClearTerminalSequence) // Clear Terminal
	}

	log.SetFormatter(&log.TextLayout{ForceColor: true, DisableTimestamp: true})

	i.displaySelections.(func())()

	fmt.Println("")
	tmpMenuFuncs := make(map[string]interface{})

	var count = 1
	for _, d := range i.menus {
		if d.Id != 999 && d.Enable.(func() bool)() {
			tmpMenuFuncs[strconv.Itoa(count)] = d.Task
			fmt.Printf("%d. %-25s - %s\n", count, d.Name, d.Description)
			count++
		} else if d.Id == 999 {
			tmpMenuFuncs["q"] = d.Task
			fmt.Printf("q. %-25s - %s\n", d.Name, d.Description)
		}
	}
	if i.lastErr {
		fmt.Printf("%s Invalid choice, try again. %s", log.ColorFmt(log.FgRed), log.CLRRESET)
	}
	if i.lastMsg {
		fmt.Printf("%s%s%s\n", log.ColorFmt(log.FgRed), i.msg, log.CLRRESET)
		i.lastMsg = false
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter menu item: ")
	selection, _ := reader.ReadString('\n')
	if strings.Index(selection, "\r") > -1 {
		selection = strings.TrimSuffix(selection, "\r\n")
	} else {
		selection = strings.TrimSuffix(selection, "\n")
	}

	if strings.TrimSpace(selection) != "" {
		if tmpMenuFuncs[selection] != nil {
			tmpMenuFuncs[selection].(func())()
		} else {
			//Invalid choice
			return -1
		}
	} else {
		//Invalid choice
		return -1
	}

	log.SetFormatter(&log.TextLayout{ForceColor: true})

	return 0
}

func CaptureSelection(datamap map[string]string, msg ...string) []string {
	defltmsg := "Enter row # (#:# and #,# supported, other ignored) and push enter or return (to go back): "
	if msg != nil {
		defltmsg = msg[0]
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(defltmsg)
	selection, _ := reader.ReadString('\n')
	if strings.Index(selection, "\r") > -1 {
		selection = strings.TrimSuffix(selection, "\r\n")
	} else {
		selection = strings.TrimSuffix(selection, "\n")
	}

	if strings.TrimSpace(selection) != "" {
		//Rows start at 1 and data at 0, subtract one to get actual row
		if strings.Index(selection, ":") > -1 {
			vals := strings.Split(selection, ":")
			start, _ := strconv.Atoi(vals[0])
			end, _ := strconv.Atoi(vals[1])
			idxAr := make([]int, 0)
			selections := make([]string, 0)
			for i := start; i <= end; i++ {
				idxAr = append(idxAr, i)
				log.Println("Ids selected: ", i)
				actRow := i - 1
				selection = strconv.Itoa(actRow)
				if len(datamap[selection]) > 0 {
					selections = append(selections, datamap[selection])
				}
			}
			return selections

		} else if strings.Index(selection, ",") > -1 {
			// list of numbers separated by commas
			vals := strings.Split(selection, ",")
			selections := make([]string, 0)
			for _,d := range vals {
				// get value and subtract one for actual value of each
				v,err := strconv.Atoi(d)
				if err != nil {
					panic(err)
				}
				selection = strconv.Itoa(v-1)
				if len(datamap[selection]) > 0 {
					selections = append(selections, datamap[selection])
				}
			}
			return selections
		} else {
			row, _ := strconv.Atoi(selection)
			actRow := row - 1
			selection = strconv.Itoa(actRow)
			if len(datamap[selection]) > 0 {
				return []string{datamap[selection]}
			}
		}
	}
	return []string{}
}

func Pause() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Push enter to continue")
	_, _ = reader.ReadString('\n')
}
