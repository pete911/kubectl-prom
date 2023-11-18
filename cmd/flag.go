package cmd

import (
	"github.com/spf13/cobra"
)

type PromFlags struct {
	PromNamespace string
	PromLabel     string
}

func InitPromFlags(cmd *cobra.Command, flags *PromFlags) {
	cmd.Flags().StringVarP(
		&flags.PromNamespace,
		"prom-namespace",
		"",
		"",
		"prometheus server namespace, if not specified all namespaces are searched",
	)
	cmd.Flags().StringVarP(
		&flags.PromLabel,
		"prom-label",
		"",
		"app.kubernetes.io/name=prometheus,app.kubernetes.io/instance=prometheus",
		"prometheus server label selector",
	)
}
