package commands

import "github.com/spf13/cobra"

type Flags struct {
	WithServer int
	Method     string
	Paths      []string
	OutputDir  string
	HooksFile  string
	RulesFile  string
}

func bindFlags(cmd *cobra.Command) *Flags {
	flags := &Flags{}

	cmd.Flags().IntVar(&flags.WithServer, "with-server", 0, "Which server url to use from OAS. 0 = First URL.")
	cmd.Flags().StringVar(&flags.Method, "method", "*", "Choose which requests you want to send from your OAS by method. Default is all methods.")
	cmd.Flags().StringSliceVar(&flags.Paths, "path", []string{}, "Choose which requests you want to send from your OAS by path. Default is all paths.")
	cmd.Flags().StringVar(&flags.OutputDir, "output", ".", "Directory of output for HTTP responses.")
	cmd.Flags().StringVar(&flags.HooksFile, "hooks", "", "Location of hooks file to modify request values.")
	cmd.Flags().StringVar(&flags.RulesFile, "rules", "", "Location of rules file for chain requests.")

	return flags
}
