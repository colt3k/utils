# Print

Adds a small wrapper around `olekukonko/tablewriter`.

## API Surface

`Printer.TablePrint(border bool)` renders the configured header and rows to standard output, letting callers toggle border rendering per table.

## Notes

Use this package when a CLI tool needs quick tabular output without carrying tablewriter configuration at every call site.
