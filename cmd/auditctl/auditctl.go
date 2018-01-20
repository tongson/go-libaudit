package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/elastic/go-libaudit"
)

var (
	fs		= flag.NewFlagSet("auditctl", flag.ExitOnError)
	del		= fs.Bool("D", false, "delete all rules")
	list	= fs.Bool("l", false, "list rules")
)

func main() {
	fs.Parse(os.Args[1:])
	if err := isRoot(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if *del {
		if err := deleteRules(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		    os.Exit(1)
		}
	}
	if *list {
		if err:= listRules(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
}

func isRoot() error {
	if os.Geteuid() != 0 {
		return errors.New("you must be root")
	}
	return nil
}

func deleteRules() error {
	var d io.WriteCloser
	c, err := libaudit.NewAuditClient(d)
	if err != nil {
		return errors.New("cannot setup audit client")
	}
	defer c.Close()
	n, err := c.DeleteRules()
	if err != nil {
		return errors.New("cannot delete rules")
	}
	fmt.Fprintf(os.Stdout, "%v rules deleted\n", n)
	return nil
}

func listRules() error {
	var d io.WriteCloser
	c, err := libaudit.NewAuditClient(d)
	if err != nil {
		return errors.New("cannot setup audit client")
	}
	defer c.Close()
	rules, err := c.GetRules()
	if err != nil {
		return errors.New("cannot get rules")
	}
	for i, rule := range rules {
		fmt.Fprintf(os.Stdout, "rule %v:\n%v\n", i, hex.EncodeToString(rule))
    }
	return nil
}
