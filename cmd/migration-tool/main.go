package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/IceWhaleTech/CasaOS/types"
)

func main() {
	versionFlag := flag.Bool("v", false, "version")
	flag.Parse()

	if *versionFlag {
		fmt.Println(types.CURRENTVERSION)
		os.Exit(0)
	}

	fmt.Println("This migration tool is not implemented yet.")
}
