package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ListVpcs lists all VPCs in the specified session.
func ListVpcs(sess *session.Session) ([]*ec2.Vpc, error) {
	svc := ec2.New(sess)

	result, err := svc.DescribeVpcs(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list VPCs: %v", err)
	}

	return result.Vpcs, nil
}

// ListSubnetsForVpc lists all subnets for the specified VPC ID in the specified session.
func ListSubnetsForVpc(sess *session.Session, vpcID string) ([]*ec2.Subnet, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	result, err := svc.DescribeSubnets(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list subnets for VPC %s: %v", vpcID, err)
	}

	return result.Subnets, nil
}

// ListNatGatewaysForVpc lists all NAT gateways for the specified VPC ID in the specified session.
func ListNatGatewaysForVpc(sess *session.Session, vpcID string) ([]*ec2.NatGateway, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	result, err := svc.DescribeNatGateways(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list NAT gateways for VPC %s: %v", vpcID, err)
	}

	return result.NatGateways, nil
}

// ListVpcEndpointsForVpc lists all VPC endpoints for the specified VPC ID in the specified session.
func ListVpcEndpointsForVpc(sess *session.Session, vpcID string) ([]*ec2.VpcEndpoint, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeVpcEndpointsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	result, err := svc.DescribeVpcEndpoints(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list VPC endpoints for VPC %s: %v", vpcID, err)
	}

	return result.VpcEndpoints, nil
}

// ListEipsForVpc lists all Elastic IP addresses for the specified VPC ID in the specified session.
func ListEipsForVpc(sess *session.Session, vpcID string) ([]*ec2.Address, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("domain"),
				Values: []*string{aws.String("vpc")},
			},
			// constrain to the specified vpcID
			{
				Name:   aws.String("network-interface-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	result, err := svc.DescribeAddresses(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list Elastic IPs for VPC %s: %v", vpcID, err)
	}

	return result.Addresses, nil
}

// ListIgwsForVpc lists all Internet gateways for the specified VPC ID in the specified session.
func ListIgwsForVpc(sess *session.Session, vpcID string) ([]*ec2.InternetGateway, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	result, err := svc.DescribeInternetGateways(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list Internet gateways for VPC %s: %v", vpcID, err)
	}

	return result.InternetGateways, nil
}

func DeleteAllVpcs(sess *session.Session, force bool) error {
	vpcs, err := ListVpcs(sess)
	if err != nil {
		return fmt.Errorf("failed to list VPCs: %v", err)
	}

	for _, vpc := range vpcs {
		err := DeleteVpc(sess, vpc, force)
		if err != nil && !ignoreErrors {
			return err
		}
		fmt.Println("Deleted VPC", *vpc.VpcId)
	}

	return nil
}

// DeleteVpc deletes the specified VPC, along with all associated resources, in the specified session.
// change to vpc pointer
func DeleteVpc(sess *session.Session, vpc *ec2.Vpc, force bool) error {
	vpcID := *vpc.VpcId
	fmt.Println("Deleting VPC", vpcID)
	// List all associated resources for the VPC.
	subnets, err := ListSubnetsForVpc(sess, vpcID)

	if err != nil {
		fmt.Print("failed to list subnets for VPC %s: %v", vpcID, err)
		if !ignoreErrors {
			return err
		}
	}

	natGateways, err := ListNatGatewaysForVpc(sess, vpcID)
	if err != nil {
		fmt.Print("failed to list NAT gateways for VPC %s: %v", vpcID, err)
		if !ignoreErrors {
			return err
		}
	}

	igws, err := ListIgwsForVpc(sess, vpcID)
	if err != nil {
		fmt.Print("failed to list Internet gateways for VPC %s: %v", vpcID, err)
		if !ignoreErrors {
			return err
		}
	}

	vpcEndpoints, err := ListVpcEndpointsForVpc(sess, vpcID)
	if err != nil {
		fmt.Print("failed to list VPC endpoints for VPC %s: %v", vpcID, err)
		if !ignoreErrors {
			return err
		}
	}

	// delete all non-default route tables
	routeTables, err := ListRouteTablesForVpc(sess, vpc)
	if err != nil {
		fmt.Print("failed to list route tables for VPC %s: %v", vpcID, err)
		if !ignoreErrors {
			return err
		}
	}

	eips, err := ListEipsForVpc(sess, vpcID)
	if err != nil {
		fmt.Print("failed to list Elastic IPs for VPC %s: %v", vpcID, err)
		if !ignoreErrors {
			return err
		}
	}

	nacls, err := ListNaclsForVpc(sess, vpc)
	if err != nil {
		fmt.Print("failed to list network ACLs for VPC %s: %v", vpcID, err)
		if !ignoreErrors {
			return err
		}
	}

	sgs, err := ListSgsForVpc(sess, vpcID)
	if err != nil {
		fmt.Print("failed to list security groups for VPC %s: %v", vpcID, err)
		if !ignoreErrors {
			return err
		}
	}

	if len(vpcEndpoints) > 0 {
		fmt.Printf("Deleting %d VPC endpoints in VPC %s...\n", len(vpcEndpoints), vpcID)
		err := DeleteVpcEndpoints(sess, vpcEndpoints)
		if err != nil {
			fmt.Print("failed to delete VPC endpoints for VPC %s: %v", vpcID, err)
			if !ignoreErrors {
				return err
			}
		}
	}

	if len(natGateways) > 0 {
		fmt.Printf("Deleting %d NAT gateways in VPC %s...\n", len(natGateways), vpcID)
		err := DeleteNatGateways(sess, natGateways)
		if err != nil {
			fmt.Print("failed to delete NAT gateways for VPC %s: %v", vpcID, err)
			if !ignoreErrors {
				return err
			}
		}
	}

	if len(eips) > 0 {
		fmt.Printf("Releasing %d Elastic IPs in VPC %s...\n", len(eips), vpcID)
		err := ReleaseEips(sess, eips)
		if err != nil {
			fmt.Print("failed to release Elastic IPs for VPC %s: %v", vpcID, err)
			if !ignoreErrors {
				return err
			}
		}
		err = DeleteEips(sess, eips)
		if err != nil {
			fmt.Print("failed to delete Elastic IPs for VPC %s: %v", vpcID, err)
			if !ignoreErrors {
				return err
			}
		}
	}

	if len(igws) > 0 {
		fmt.Printf("Detaching and deleting %d Internet gateways in VPC %s...\n", len(igws), vpcID)
		err := DetachAndDeleteIgws(sess, igws)
		if err != nil {
			fmt.Print("failed to detach and delete Internet gateways for VPC %s: %v", vpcID, err)
			if !ignoreErrors {
				return err
			}
		}
	}

	if len(routeTables) > 0 {
		fmt.Printf("Deleting %d route tables in VPC %s...\n", len(routeTables), vpcID)
		err := DeleteRouteTables(sess, routeTables)
		if err != nil {
			fmt.Print("failed to delete route tables for VPC %s: %v", vpcID, err)
			if !ignoreErrors {
				return err
			}
		}
	}

	if len(sgs) > 0 {
		fmt.Printf("Deleting %d security groups in VPC %s...\n", len(sgs), vpcID)
		err := DeleteSgs(sess, sgs)
		if err != nil {
			fmt.Print("failed to delete security groups for VPC %s: %v", vpcID, err)
			if !ignoreErrors {
				return err
			}
		}
	}

	if len(nacls) > 0 {
		fmt.Printf("Deleting %d network ACLs in VPC %s...\n", len(nacls), vpcID)
		err := DeleteNacls(sess, nacls)
		if err != nil {
			fmt.Print("failed to delete network ACLs for VPC %s: %v", vpcID, err)
			if !ignoreErrors {
				return err
			}
		}
	}

	// Delete all associated resources for the VPC.
	if len(subnets) > 0 {
		fmt.Printf("Deleting %d subnets in VPC %s...\n", len(subnets), vpcID)
		err := DeleteSubnets(sess, subnets)
		if err != nil {
			fmt.Print("failed to delete subnets for VPC %s: %v", vpcID, err)
			if !ignoreErrors {
				return err
			}
		}
	}

	// Delete the VPC itself.
	fmt.Printf("Deleting VPC %s...\n", vpcID)
	err = DeleteVpcAndWait(sess, vpc)
	if err != nil {
		fmt.Print("failed to delete VPC %s: %v", vpcID, err)
		if !ignoreErrors {
			return err
		}
	}

	return nil
}

func ListNaclsForVpc(sess *session.Session, vpc *ec2.Vpc) ([]*ec2.NetworkAcl, error) {
	svc := ec2.New(sess)
	input := &ec2.DescribeNetworkAclsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{vpc.VpcId},
			},
		},
	}
	result, err := svc.DescribeNetworkAcls(input)
	if err != nil {
		return nil, err
	}
	return result.NetworkAcls, nil
}

func DeleteNacls(sess *session.Session, nacls []*ec2.NetworkAcl) error {
	svc := ec2.New(sess)
	fmt.Printf("Deleting %d network ACLs...\n", len(nacls))

	for _, nacl := range nacls {
		fmt.Printf("Deleting network ACL %s...\n", *nacl.NetworkAclId)
		input := &ec2.DeleteNetworkAclInput{
			NetworkAclId: nacl.NetworkAclId,
		}
		_, err := svc.DeleteNetworkAcl(input)
		if err != nil {
			fmt.Printf("Error deleting network ACL %s: %v", *nacl.NetworkAclId, err)
			if !ignoreErrors {
				return err
			}
		}
	}
	return nil
}

func DeleteRouteTables(sess *session.Session, tables []*ec2.RouteTable) error {
	svc := ec2.New(sess)
	fmt.Printf("Deleting %d route tables...\n", len(tables))

	for _, table := range tables {
		// do not delete the main route table.  Careful to check for nil.
		if table.Associations != nil {
			for _, association := range table.Associations {
				if *association.Main {
					fmt.Printf("Skipping main route table %s...\n", *table.RouteTableId)
					continue
				}

				fmt.Printf("Disassociating route table %s...\n", *table.RouteTableId)
				input := &ec2.DisassociateRouteTableInput{
					AssociationId: association.RouteTableAssociationId,
				}
				_, err := svc.DisassociateRouteTable(input)
				if err != nil {
					fmt.Printf("Error disassociating route table %s: %v", *table.RouteTableId, err)
					if !ignoreErrors {
						return err
					}
				}
			}
		}

		fmt.Printf("Deleting route table %s...\n", *table.RouteTableId)
		input := &ec2.DeleteRouteTableInput{
			RouteTableId: table.RouteTableId,
		}
		_, err := svc.DeleteRouteTable(input)
		if err != nil {
			fmt.Printf("Error deleting route table %s: %v", *table.RouteTableId, err)
		}
	}

	return nil
}

func ListRouteTablesForVpc(sess *session.Session, vpc *ec2.Vpc) ([]*ec2.RouteTable, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{vpc.VpcId},
			},
		},
	}

	result, err := svc.DescribeRouteTables(input)
	if err != nil {
		return nil, err
	}

	return result.RouteTables, nil
}

func DeleteSgs(sess *session.Session, sgs []*ec2.SecurityGroup) error {
	fmt.Println("Deleting security groups...")
	svc := ec2.New(sess)
	for _, sg := range sgs {
		if *sg.GroupName == "default" {
			continue
		}
		input := &ec2.DeleteSecurityGroupInput{
			GroupId: sg.GroupId,
		}
		_, err := svc.DeleteSecurityGroup(input)
		if err != nil {
			return err
		}
	}
	fmt.Println("Done deleting security groups.")
	return nil
}

func ListSgsForVpc(sess *session.Session, id string) ([]*ec2.SecurityGroup, error) {
	// Create an EC2 service client.
	svc := ec2.New(sess)

	// Initialize the input parameters.
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(id)},
			},
		},
	}

	// Retrieve the security groups.
	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}

	return result.SecurityGroups, nil
}

// DeleteEips deletes the specified Elastic IPs.
func DeleteEips(sess *session.Session, eips []*ec2.Address) error {
	svc := ec2.New(sess)
	for _, eip := range eips {
		// Initialize the input parameters.
		input := &ec2.ReleaseAddressInput{
			PublicIp: eip.PublicIp,
		}

		// Release the Elastic IPs.
		_, err := svc.ReleaseAddress(input)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteSubnets deletes the specified subnets.
func DeleteSubnets(sess *session.Session, subnets []*ec2.Subnet) error {
	// Create a new EC2 client using the provided session.
	ec2Svc := ec2.New(sess)

	// Delete each subnet.
	for _, subnet := range subnets {
		// Get the name of the subnet.
		name := getNameTag(subnet.Tags)

		// Delete the subnet.
		fmt.Printf("Deleting subnet %s (%s)...\n", aws.StringValue(subnet.SubnetId), name)
		_, err := ec2Svc.DeleteSubnet(&ec2.DeleteSubnetInput{
			SubnetId: subnet.SubnetId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteVpcEndpoints deletes the specified VPC endpoints.
func DeleteVpcEndpoints(sess *session.Session, vpcEndpoints []*ec2.VpcEndpoint) error {
	// Create a new EC2 client using the provided session.
	ec2Svc := ec2.New(sess)

	// Delete each VPC endpoint.
	for _, vpcEndpoint := range vpcEndpoints {
		fmt.Printf("Deleting VPC endpoint %s...\n", aws.StringValue(vpcEndpoint.VpcEndpointId))

		if forceFlag {
			_, err := ec2Svc.DeleteVpcEndpoints(&ec2.DeleteVpcEndpointsInput{
				VpcEndpointIds: []*string{vpcEndpoint.VpcEndpointId},
			})
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Skipping VPC endpoint deletion. Use the --force flag to force deletion.")
		}
	}

	return nil
}

// DeleteNatGateways deletes the specified NAT gateways and waits for them to be deleted.
func DeleteNatGateways(sess *session.Session, natGateways []*ec2.NatGateway) error {
	fmt.Println("Deleting NAT gateways...")
	// Create a new EC2 client using the provided session.
	ec2Svc := ec2.New(sess)

	// Delete each NAT gateway.
	for _, natGw := range natGateways {
		fmt.Printf("Deleting NAT gateway %s...\n", aws.StringValue(natGw.NatGatewayId))

		if forceFlag {
			_, err := ec2Svc.DeleteNatGateway(&ec2.DeleteNatGatewayInput{
				NatGatewayId: natGw.NatGatewayId,
			})
			if err != nil {
				return err
			}
			// Wait for the NAT gateways to be deleted.
			//fmt.Println("Waiting for NAT gateways to be deleted...")
			//err = ec2Svc.WaitUntilNatGatewayDeleted(&ec2.DescribeNatGatewaysInput{})
			//if err != nil {
			//	return err
			//}
		} else {
			fmt.Println("Skipping NAT gateway deletion. Use the --force flag to force deletion.")
		}
	}

	fmt.Println("NAT gateways deleted.")
	return nil
}

// ReleaseEips releases the specified EIPs.
func ReleaseEips(sess *session.Session, eips []*ec2.Address) error {
	fmt.Println("Releasing EIPs...")
	// Create a new EC2 client using the provided session.
	ec2Svc := ec2.New(sess)

	// Release each EIP.
	for _, eip := range eips {
		fmt.Printf("Releasing EIP %s...\n", aws.StringValue(eip.PublicIp))

		if forceFlag {
			_, err := ec2Svc.ReleaseAddress(&ec2.ReleaseAddressInput{
				PublicIp: eip.PublicIp,
			})
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Skipping EIP release. Use the --force flag to force release.")
		}
	}

	fmt.Println("EIPs released.")
	return nil
}

// DetachAndDeleteIgws detaches and deletes the specified Internet gateways.
func DetachAndDeleteIgws(sess *session.Session, igws []*ec2.InternetGateway) error {
	fmt.Println("Detaching and deleting Internet gateways...")
	// Create a new EC2 client using the provided session.
	ec2Svc := ec2.New(sess)

	// Detach and delete each Internet gateway.
	for _, igw := range igws {
		// Get the name of the Internet gateway.
		name := getNameTag(igw.Tags)

		// Detach the Internet gateway from its VPC.
		vpcId := aws.StringValue(igw.Attachments[0].VpcId)
		fmt.Printf("Detaching Internet gateway %s (%s) from VPC %s...\n", aws.StringValue(igw.InternetGatewayId), name, vpcId)

		if forceFlag {
			_, err := ec2Svc.DetachInternetGateway(&ec2.DetachInternetGatewayInput{
				InternetGatewayId: igw.InternetGatewayId,
				VpcId:             aws.String(vpcId),
			})
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Skipping Internet gateway detachment. Use the --force flag to force detachment.")
		}

		// Wait for the Internet gateway to be detached.
		if forceFlag {
			fmt.Println("Waiting for Internet gateway to be detached...")
			err := ec2Svc.WaitUntilInternetGatewayExists(&ec2.DescribeInternetGatewaysInput{ // TODO WaitUntilInternetGatewayDetached is not available in the SDK.
				InternetGatewayIds: []*string{igw.InternetGatewayId},
			})
			if err != nil {
				return err
			}
		}

		// Delete the Internet gateway.
		fmt.Printf("Deleting Internet gateway %s (%s)...\n", aws.StringValue(igw.InternetGatewayId), name)

		if forceFlag {
			_, err := ec2Svc.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
				InternetGatewayId: igw.InternetGatewayId,
			})
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Skipping Internet gateway deletion. Use the --force flag to force deletion.")
		}

	}

	fmt.Println("Internet gateways detached and deleted.")
	return nil
}

// getNameTag returns the value of the "Name" tag for the specified resource, or an empty string if the tag is not present.
func getNameTag(tags []*ec2.Tag) string {
	for _, tag := range tags {
		if aws.StringValue(tag.Key) == "Name" {
			return aws.StringValue(tag.Value)
		}
	}
	return ""
}

// DeleteVpcAndWait deletes the specified VPC and waits for it to be deleted.
func DeleteVpcAndWait(sess *session.Session, vpc *ec2.Vpc) error {
	fmt.Println("Deleting VPC...")
	// Create a new EC2 client using the provided session.
	ec2Svc := ec2.New(sess)

	// Get the name of the VPC.
	name := getNameTag(vpc.Tags)

	// Delete the VPC.
	fmt.Printf("Deleting VPC %s (%s)...\n", aws.StringValue(vpc.VpcId), name)
	if forceFlag {
		_, err := ec2Svc.DeleteVpc(&ec2.DeleteVpcInput{
			VpcId: vpc.VpcId,
		})
		if err != nil {
			fmt.Println("Error deleting VPC:", err)
			if !ignoreErrors {
				return err
			}
		}
		// Wait for the VPC to be deleted.
		fmt.Println("Waiting for VPC to be deleted...")
		//err = ec2Svc.WaitUntilVpc(&ec2.DescribeVpcsInput{
		//	VpcIds: []*string{vpc.VpcId},
		//})
		//if err != nil {
		//	return err
		//}
	} else {
		fmt.Println("Skipping VPC deletion. Use the --force flag to force deletion.")
	}

	fmt.Println("VPC deleted, I think.")

	return nil
}
