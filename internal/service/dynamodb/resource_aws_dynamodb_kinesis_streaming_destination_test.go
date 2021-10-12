package aws

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/service/dynamodb/finder"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	tfdynamodb "github.com/hashicorp/terraform-provider-aws/internal/service/dynamodb"
	tfdynamodb "github.com/hashicorp/terraform-provider-aws/internal/service/dynamodb"
	tfdynamodb "github.com/hashicorp/terraform-provider-aws/internal/service/dynamodb"
	tfdynamodb "github.com/hashicorp/terraform-provider-aws/internal/service/dynamodb"
	tfdynamodb "github.com/hashicorp/terraform-provider-aws/internal/service/dynamodb"
)

func TestAccAwsDynamoDbKinesisStreamingDestination_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_dynamodb_kinesis_streaming_destination.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, dynamodb.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAWSDynamoDbKinesisStreamingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsDynamodbKinesisStreamingDestinationConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynamoDbKinesisStreamingDestinationExists(resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "stream_arn", "kinesis", regexp.MustCompile(fmt.Sprintf("stream/%s", rName))),
					resource.TestCheckResourceAttr(resourceName, "table_name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsDynamoDbKinesisStreamingDestination_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_dynamodb_kinesis_streaming_destination.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, dynamodb.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAWSDynamoDbKinesisStreamingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsDynamodbKinesisStreamingDestinationConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynamoDbKinesisStreamingDestinationExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceKinesisStreamingDestination(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAwsDynamoDbKinesisStreamingDestination_disappears_DynamoDbTable(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resourceName := "aws_dynamodb_kinesis_streaming_destination.test"
	tableResourceName := "aws_dynamodb_table.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, dynamodb.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAWSDynamoDbKinesisStreamingDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsDynamodbKinesisStreamingDestinationConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDynamoDbKinesisStreamingDestinationExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceTable(), tableResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAwsDynamodbKinesisStreamingDestinationConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = %[1]q
  read_capacity  = 10
  write_capacity = 10
  hash_key       = "hk"

  attribute {
    name = "hk"
    type = "S"
  }
}

resource "aws_kinesis_stream" "test" {
  name        = %[1]q
  shard_count = 2
}

resource "aws_dynamodb_kinesis_streaming_destination" "test" {
  table_name = aws_dynamodb_table.test.name
  stream_arn = aws_kinesis_stream.test.arn
}
`, rName)
}

func testAccCheckDynamoDbKinesisStreamingDestinationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		tableName, streamArn, err := dynamoDbKinesisStreamingDestinationParseId(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBConn

		output, err := tfdynamodb.FindDynamoDBKinesisDataStreamDestination(context.Background(), conn, streamArn, tableName)

		if err != nil {
			return err
		}

		if output == nil {
			return fmt.Errorf("DynamoDB Kinesis Streaming Destination (%s) not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckAWSDynamoDbKinesisStreamingDestinationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_dynamodb_kinesis_streaming_destination" {
			continue
		}

		tableName, streamArn, err := dynamoDbKinesisStreamingDestinationParseId(rs.Primary.ID)

		if err != nil {
			return err
		}

		output, err := tfdynamodb.FindDynamoDBKinesisDataStreamDestination(context.Background(), conn, streamArn, tableName)

		if tfawserr.ErrCodeEquals(err, dynamodb.ErrCodeResourceNotFoundException) {
			continue
		}

		if err != nil {
			return err
		}

		if output != nil {
			return fmt.Errorf("DynamoDB Kinesis Streaming Destination (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}