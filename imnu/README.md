# Interactive Menu

Provides framework for an interactive menu system on the command line.

## API Surface

- `New(displaySelections, menus)` creates an interactive menu controller.
- `CaptureSelection` collects a filtered choice from a map of values.
- `Pause` blocks until the operator confirms they are ready to continue.

## Menu Flow

Each `Menu` entry defines an `Id`, `Name`, `Description`, `Enable` predicate, and `Task`. The menu only exposes items that are currently enabled, making it suitable for CLI workflows with prerequisites.

EXAMPLE:



menu = imnu.New(menuDisplaySelections, []imnu.Menu{
	{Id: 1, Name: "List xxx", Description: "Lists xxx", Enable: func() bool { return true }, Task: func() { ID = listXXX() }},
	{Id: 2, Name: "List yyy", Description: "List yyy", Enable: func() bool { return len(ID) > 0 }, Task: func() { yyyID = listYYY() }},
	{Id: 999, Name: "Quit", Description: "Quit program", Task: func() { log.Println("Exiting.."); os.Exit(0) }},
})
menu.StartMenu()

## Development

- `cd imnu && go test ./...`
