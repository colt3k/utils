# Regx

Various defined regular expressions

## Included Helpers

The package exposes common validation and extraction helpers such as:

- `Find`, `Match`
- `CC` for credit-card style patterns
- `DATE`
- `ISBN`
- `POSTALCODE`
- `PRICE`

## Notes

Pattern definitions live in `patterns.go`. The package is most useful when a caller needs a named regular expression without carrying the pattern text inline.
