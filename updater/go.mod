module github.com/colt3k/utils/updater

go 1.17

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/colt3k/nglog v0.0.28
	github.com/colt3k/utils/debug v0.0.1
	github.com/colt3k/utils/file v0.0.10
	github.com/colt3k/utils/io v0.0.18
	github.com/colt3k/utils/mathut v0.0.2
	github.com/colt3k/utils/netut v0.0.23
	github.com/colt3k/utils/osut v0.0.5
	github.com/colt3k/utils/ques v0.0.4
	github.com/colt3k/utils/version v0.0.3
)

replace golang.org/x/net => golang.org/x/net v0.19.0 //CVE-2023-48795

require (
	github.com/colt3k/utils/archive v0.0.9 // indirect
	github.com/colt3k/utils/encode v0.0.5 // indirect
	github.com/colt3k/utils/hash v0.0.7 // indirect
	github.com/colt3k/utils/stats v0.0.1 // indirect
	github.com/go-mail/mail v2.3.1+incompatible // indirect
	github.com/gonutz/w32 v1.0.0 // indirect
	github.com/gorilla/http v1.0.0 // indirect
	github.com/iancoleman/orderedmap v0.3.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/str v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/term v0.15.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)
