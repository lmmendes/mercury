package main

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"

	"github.com/spf13/pflag"
)

func initFlags() *koanf.Koanf {
	ko := koanf.New(".")

	f := pflag.NewFlagSet("config", pflag.ContinueOnError)

	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	f.Bool("install", false, "setup database (first time)")
	f.Bool("upgrade", false, "upgrade database to the current version")

	if err := f.Parse(os.Args[1:]); err != nil {
		logger.Fatalf("error loading flags: %v", err)
	}

	if err := ko.Load(posflag.Provider(f, ".", ko), nil); err != nil {
		logger.Fatalf("error loading config: %v", err)
	}

	return ko
}
