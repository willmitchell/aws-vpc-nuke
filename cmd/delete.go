package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a VPC and all of its associated resources",
	Long:  "Delete a VPC and all of its associated resources",
	RunE:  deleteFunc,
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("vpc-id", "v", "", "the ID of the VPC to delete")
	viper.BindPFlag("vpc-id", deleteCmd.Flags().Lookup("vpc-id"))
}

func deleteFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("delete called")

	regionList := viper.GetStringSlice("region-list")
	profileList := viper.GetStringSlice("profile-list")

	// Delete the VPC and all associated resources.
	err := IterateOverProfiles(profileList, func(profile string) error {
		return IterateOverRegions(regionList, func(region string) error {
			fmt.Printf("Deleting VPCs in %s (%s)\n", profile, region)

			// Use the current profile and region to create a new session.
			sess, err := GetSession(profile, region)
			if err != nil {
				return fmt.Errorf("failed to create session for profile %s and region %s: %v", profile, region, err)
			}

			// Delete the VPC and all associated resources using the session.
			err = DeleteAllVpcs(sess, forceFlag)
			if err != nil {
				return fmt.Errorf("failed to delete VPC resources: %v", err)
			}

			return nil
		})
	})
	if err != nil {
		return fmt.Errorf("failed to IterateOverProfiles: %v", err)
	}
	return nil
}
