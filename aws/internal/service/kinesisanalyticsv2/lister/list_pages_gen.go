// Code generated by "aws/internal/generators/listpages/main.go -function=ListApplications -paginator=NextToken github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"; DO NOT EDIT.

package lister

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ListApplicationsPages(conn *kinesisanalyticsv2.KinesisAnalyticsV2, input *kinesisanalyticsv2.ListApplicationsInput, fn func(*kinesisanalyticsv2.ListApplicationsOutput, bool) bool) error {
	return ListApplicationsPagesWithContext(context.Background(), conn, input, fn)
}

func ListApplicationsPagesWithContext(ctx context.Context, conn *kinesisanalyticsv2.KinesisAnalyticsV2, input *kinesisanalyticsv2.ListApplicationsInput, fn func(*kinesisanalyticsv2.ListApplicationsOutput, bool) bool) error {
	for {
		output, err := conn.ListApplicationsWithContext(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.StringValue(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
