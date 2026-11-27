package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type playbookList []string

func (p *playbookList) String() string { return strings.Join(*p, ",") }

func (p *playbookList) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func main() {
	var playbooks playbookList
	flag.Var(&playbooks, "playbook", "Path to an Ansible playbook YAML file (repeatable)")
	rolesPath := flag.String("roles", "", "Path to the Ansible roles directory")
	showVersion := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	if len(playbooks) == 0 || *rolesPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	unused, err := unusedRoles(playbooks, *rolesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, role := range unused {
		fmt.Println(role)
	}
}
