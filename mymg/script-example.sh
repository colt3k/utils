#!/bin/bash

# parameters passed from mymg
#   from is the local location
from=$1
#   to is the file name without a path
to=$2

echo put $1 /path/$2 | sftp hostname