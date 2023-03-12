package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// GetSession creates a new AWS session using the provided profile and region.
func GetSession(profile, region string) (*session.Session, error) {
	if debugFlag {
		fmt.Println("GetSession called")
	}
	// Create a new AWS session using the provided profile and region.
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return sess, nil
}

// IterateOverProfiles calls the provided function for each profile in the profileList.
func IterateOverProfiles(profileList []string, fn func(string) error) error {
	fmt.Println("IterateOverProfiles called, profileList: ", profileList)
	for _, profile := range profileList {
		err := fn(profile)
		if err != nil {
			return err
		}
	}
	return nil
}

// IterateOverRegions calls the provided function for each region in the regionList.
func IterateOverRegions(regionList []string, fn func(string) error) error {
	fmt.Println("IterateOverRegions called, regionList: ", regionList)
	for _, region := range regionList {
		err := fn(region)
		if err != nil {
			return err
		}
	}
	return nil
}
