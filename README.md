# aws-vpc-nuke

`aws-vpc-nuke` is a CLI tool for deleting all VPC resources across regions and profiles.

## Warning

This is a highly destructive tool, similar to aws-nuke.  
This tool actually has fewer safeguards than aws-nuke, so please be careful in using it.

You must only use this tool if you are 100% sure you want to delete all VPC resources in all specified regions and profiles.

The one safeguard is that you must specify the `--force` flag to actually delete anything.

USE AT YOUR OWN RISK.  NO WARRANTIES ARE EXPRESSED OR IMPLIED.

## Supported Resources

The following resources are supported:

- VPCs
- Subnets
- Internet Gateways
- NAT Gateways
- Route Tables
- NACLs
- Security Groups

## Usage

The output of the `--help` flag is below.  Note that profiles and region specifiers are comma-separated lists.

```bash
A command-line tool for deleting all VPC resources in an AWS account

Usage:
  aws-vpc-nuke [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a VPC and all of its associated resources
  help        Help about any command
  list        List all VPC resources in the specified regions and profiles

Flags:
  -d, --debug                  Enable debug logging
  -f, --force                  Force the deletion of all VPC resources without confirmation
  -h, --help                   help for aws-vpc-nuke
  -i, --ignore-errors          Ignore deletion errors and continue deleting resources
  -p, --profile-list strings   Comma-separated list of AWS profiles to use
  -r, --region-list strings    Comma-separated list of AWS regions to use (default [us-west-2])

Use "aws-vpc-nuke [command] --help" for more information about a command.
```

## Why I created this tool

aws-nuke is a great tool, but I found that its super-safe operational model was not suitable for my use case.  I wanted to be able to delete all VPC resources in all regions across a set of profiles (accounts), but I didn't want to have to specify each resource type individually.  I also wanted to be able to delete all resources in a single command.  I also wanted to tool to delete what it could, even if some resources could not be deleted.

My use case actually stems from experimental use of AWS Control Tower.  This is an interesting management tool for enterprise AWS account management, but it can create a lot of resources across multiple regions.  In my case, I set up Control Tower and the associated Account Factory for Terraform, and I deployed it across 3 regions.  I created a couple test accounts across those same regions. My daily costs went up to over $30/day.  Hello VPC Endpoints and multiple NAT Gateways!

## Limitations

- There are likely many resource types that could be added.
- Logging is decent, but messy, and English only.  
- The tool is single-threaded.
- I tried using blocking calls for things like unmapping-then-deleting NAT Gateways, but could not get this working.  So, running the tool repeatedly can be helpful. 

## Thanks

- cobra CLI 
- viper
- aws-nuke
