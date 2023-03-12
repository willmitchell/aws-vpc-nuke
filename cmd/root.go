package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

var (
	regionList []string
	forceFlag  bool

	ignoreErrors bool

	debugFlag bool

	profileList []string
)

var rootCmd = &cobra.Command{
	Use:   "aws-vpc-nuke",
	Short: "A command-line tool for deleting all VPC resources in an AWS account",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set the regionList, force, and profileList flags in the Viper configuration.
		viper.Set("regionList", regionList)
		viper.Set("force", forceFlag)
		viper.Set("ignoreErrors", ignoreErrors)
		viper.Set("debug", debugFlag)
		viper.Set("profileList", profileList)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add flags to the root command.
	rootCmd.PersistentFlags().StringSliceVarP(&regionList, "region-list", "r", []string{"us-west-2"}, "Comma-separated list of AWS regions to use")
	rootCmd.PersistentFlags().BoolVarP(&forceFlag, "force", "f", false, "Force the deletion of all VPC resources without confirmation")
	rootCmd.PersistentFlags().BoolVarP(&ignoreErrors, "ignore-errors", "i", false, "Ignore deletion errors and continue deleting resources")
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringSliceVarP(&profileList, "profile-list", "p", []string{""}, "Comma-separated list of AWS profiles to use")
	// why are these not working?
	err := viper.BindPFlag("region-list", rootCmd.PersistentFlags().Lookup("region-list"))
	if err != nil {
		fmt.Println(err)
	}
	err = viper.BindPFlag("force", rootCmd.PersistentFlags().Lookup("force"))
	if err != nil {
		fmt.Println(err)
	}
	err = viper.BindPFlag("ignore-errors", rootCmd.PersistentFlags().Lookup("ignore-errors"))
	if err != nil {
		fmt.Println(err)
	}
	err = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		fmt.Println(err)
	}
	err = viper.BindPFlag("profile-list", rootCmd.PersistentFlags().Lookup("profile-list"))
	if err != nil {
		fmt.Println(err)
	}

	// print out flags and their values

	if debugFlag {
		rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
			fmt.Printf("Flag: %s, Value: %s", flag.Name, flag.Value)
		})
	}
}

func initConfig() {
	// Set up Viper to read configuration from environment variables and/or configuration files.
	viper.AutomaticEnv()

	// Set the default configuration file name.
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file, %s", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		res, err := fmt.Fprintln(os.Stderr, err)
		fmt.Printf("Result: %v, Error: %v\n", res, err)
		os.Exit(1)
	}
}
