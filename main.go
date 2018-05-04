package main

import (
	"fmt"
	"os"

	gocongresscmd "github.com/welaw/go-congress/cmd/go-congress"
)

func main() {
	if err := gocongresscmd.MakeRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
