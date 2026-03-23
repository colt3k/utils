# webut

Small helpers for web and JSON interoperability.

## API Surface

`ConvertibleBoolean` accepts JSON booleans that may arrive either as native booleans or quoted strings such as `"true"` and `"false"`.

## Notes

Use this type inside request or response structs when the upstream API is inconsistent about boolean encoding.
