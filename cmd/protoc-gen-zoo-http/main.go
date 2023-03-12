package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.1.0"

var (
	showVersion         = flag.Bool("version", false, "print the version and exit")
	omitempty           = flag.Bool("omitempty", true, "omit if google.api is empty")
	allowDeleteBody     = flag.Bool("allow_delete_body", false, "allow delete body")
	allowEmptyPatchBody = flag.Bool("allow_empty_patch_body", false, "allow empty patch body")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-zoo-http %v\n", version)
		return
	}

	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(run)
}
