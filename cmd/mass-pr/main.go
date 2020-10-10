package main

import (
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	authToken string
	prefix    string
	org       string

	rootCmd = &cobra.Command{
		Use:   "mass-pr",
		Short: "mass-pr is a tool for bulk PR creation of GitHub repositories",
		Long:  `mass-pr is a tool for bulk PR creation of GitHub repositories`,
	}
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringVarP(&authToken, "auth-token", "t", os.Getenv("GITHUB_TOKEN"), "GitHub auth token")
	rootCmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", "", "Prefix of repository names")
	rootCmd.PersistentFlags().StringVarP(&org, "org", "o", "", "GitHub org to use")

	flag.Parse()
	err := rootCmd.Execute()
	if err != nil {
		glog.Error(err)
	}
}
