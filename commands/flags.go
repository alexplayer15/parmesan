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

func bindWithServerFlag(cmd *cobra.Command, flags *Flags) {
	cmd.Flags().IntVar(&flags.WithServer, "with-server", 0, "Which server url to use from OAS. 0 = First URL.")
}

func bindMethodFlag(cmd *cobra.Command, flags *Flags) {
	cmd.Flags().StringVar(&flags.Method, "method", "*", "Filter by HTTP method.")
}

func bindPathsFlag(cmd *cobra.Command, flags *Flags) {
	cmd.Flags().StringSliceVar(&flags.Paths, "path", []string{}, "Filter by path.")
}

func bindOutputDirFlag(cmd *cobra.Command, flags *Flags) {
	cmd.Flags().StringVar(&flags.OutputDir, "output", ".", "Directory of output.")
}

func bindHooksFileFlag(cmd *cobra.Command, flags *Flags) {
	cmd.Flags().StringVar(&flags.HooksFile, "hooks", "", "Location of hooks file.")
}

func bindRulesFileFlag(cmd *cobra.Command, flags *Flags) {
	cmd.Flags().StringVar(&flags.RulesFile, "rules", "", "Location of rules file.")
}

func bindGenerateRequestFlags(cmd *cobra.Command) *Flags {
	flags := &Flags{}
	bindWithServerFlag(cmd, flags)
	bindOutputDirFlag(cmd, flags)
	return flags
}

func bindSendRequestFlags(cmd *cobra.Command) *Flags {
	flags := &Flags{}
	bindWithServerFlag(cmd, flags)
	bindOutputDirFlag(cmd, flags)
	bindPathsFlag(cmd, flags)
	bindMethodFlag(cmd, flags)
	bindHooksFileFlag(cmd, flags)
	return flags
}

func bindChainRequestFlags(cmd *cobra.Command) *Flags {
	flags := &Flags{}
	bindWithServerFlag(cmd, flags)
	bindOutputDirFlag(cmd, flags)
	bindHooksFileFlag(cmd, flags)
	bindRulesFileFlag(cmd, flags)
	return flags
}
