package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPC resources in the specified regions and profiles",
	Long:  "List all VPC resources in the specified regions and profiles",
	RunE:  listFunc,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceP("region-list", "r", []string{}, "comma-separated list of AWS regions to search for VPC resources")
	viper.BindPFlag("region-list", listCmd.Flags().Lookup("region-list"))

	listCmd.Flags().StringSliceP("profile-list", "p", []string{}, "comma-separated list of AWS profiles to use for authentication")
	viper.BindPFlag("profile-list", listCmd.Flags().Lookup("profile-list"))
}

func listFunc(cmd *cobra.Command, args []string) error {
	regionList := viper.GetStringSlice("region-list")
	profileList := viper.GetStringSlice("profile-list")

	// Print the list of security groups in each region and profile.
	err := IterateOverProfiles(profileList, func(profile string) error {
		return IterateOverRegions(regionList, func(region string) error {
			// Use the current profile and region to create a new session.
			sess, err := GetSession(profile, region)
			if err != nil {
				return fmt.Errorf("failed to create session for profile %s and region %s: %v", profile, region, err)
			}

			// List the VPC resources using the session.
			vpcs, err2 := ListVpcs(sess)
			if err2 != nil {
				return fmt.Errorf("failed to list VPC resources: %v", err2)
			}
			// Print the list of VPCs in the current region and profile.
			fmt.Printf("VPCs in %s (%s):\n", profile, region)
			for _, vpc := range vpcs {
				fmt.Printf("\t%s\n", vpc)
			}

			return nil
		})
	})
	if err != nil {
		return fmt.Errorf("failed to list VPC resources: %v", err)
	}

	return nil
}
