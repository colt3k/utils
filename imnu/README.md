# Interactive Menu

Provides framework for an interactive menu system on the commandline

EXAMPLE:



menu = imnu.New(menuDisplaySelections, []imnu.Menu{
	{Id: 1, Name: "List xxx", Description: "Lists xxx", Enable: func() bool { return true }, Task: func() { ID = listXXX() }},
	{Id: 2, Name: "List yyy", Description: "List yyy", Enable: func() bool { return len(ID) > 0 }, Task: func() { yyyID = listYYY() }},
	{Id: 999, Name: "Quit", Description: "Quit program", Task: func() { log.Println("Exiting.."); os.Exit(0) }},
})
menu.StartMenu()
