package cmd

import (
	"errors"
	"fmt"
	"github.com/pete911/kubectl-prom/pkg/prom"
	"github.com/spf13/cobra"
)

var (
	cmdList = &cobra.Command{
		Use:   "query",
		Short: "prometheus query",
		Long:  "",
		RunE:  runQueryCmd,
	}
	promFlags PromFlags
)

func init() {
	RootCmd.AddCommand(cmdList)
	InitPromFlags(cmdList, &promFlags)
}

func runQueryCmd(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("no query supplied")
	}
	prometheus, err := prom.NewPrometheus(Logger(), RestConfig(), promFlags.PromNamespace, promFlags.PromLabel)
	if err != nil {
		return err
	}
	defer prometheus.Stop()

	data, err := prometheus.Query(args[0])
	if err != nil {
		return err
	}

	fmt.Println(string(data.Result))
	return nil
}
