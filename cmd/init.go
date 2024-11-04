package main

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"

	flag "github.com/spf13/pflag"
)

func initFlags() *koanf.Koanf {
	ko := koanf.New(".")

	f := flag.NewFlagSet("config", flag.ContinueOnError)

	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	f.Bool("install", false, "setup database (first time)")
	f.Bool("upgrade", false, "upgrade database to the current version")

	if err := f.Parse(os.Args[1:]); err != nil {
		lo.Fatalf("error loading flags: %v", err)
	}

	if err := ko.Load(posflag.Provider(f, ".", ko), nil); err != nil {
		lo.Fatalf("error loading config: %v", err)
	}

	return ko
}

func initConfigFiles(files []string, ko *koanf.Koanf) {
	for _, f := range files {
		lo.Printf("Reading config: %s", f)
		if err := ko.Load(file.Provider(f), yaml.Parser()); err != nil {
			if os.IsNotExist(err) {
				lo.Fatalf("Config file not found. If there isn't one yet, run --new-config to generate one.")
			}
			lo.Fatalf("error loading config from file: %v.", err)
		}
	}
}
