#!/bin/bash

contractName=$1
if  [[ ! -n $contractName ]] ;then
	echo "contractName is empty. use as: ./build.sh contractName"
	exit 1
fi

go build -ldflags="-s -w" -o $contractName

7z a $contractName $contractName
rm -f $contractName