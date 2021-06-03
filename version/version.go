package version

// VERSION indicates which version of the binary is running.
var VERSION = "0.0.1"// OR go run -ldflags '-X version.VERSION=1.0'

// GITCOMMIT indicates which git hash the binary was built off of
var GITCOMMIT = "0abc12d" //go run -ldflags "-X version.GITCOMMIT=`git rev-parse --short HEAD`"

// BUILDDATE date in unixtime
var BUILDDATE = "1538584284"//= time.Now().UTC().Unix() //.Add(-2 * time.Hour), ON LINUX `date +%s`

// GOVERSION go lang version application was built upon
var GOVERSION = "go1.xx.x" //go run -ldflags "-X version.GOVERSION=`go version | awk '{ print $3 }'`"