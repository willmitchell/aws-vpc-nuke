# aws-vpc-nuke

`aws-vpc-nuke` is a CLI tool for deleting all VPC resources across regions and profiles.

## Warning

This is a highly destructive tool, similar to aws-nuke.
This actually has fewer safeguards thank aws-nuke, so be careful.

You must only use this tool if you are 100% sure you want to delete all VPC resources in all specified regions and profiles.

The one safeguard is that you must specify the `--force` flag to actually delete anything.

## Usage

```

