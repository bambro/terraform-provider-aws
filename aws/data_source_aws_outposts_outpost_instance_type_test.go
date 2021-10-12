package aws

import (
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/outposts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccAWSOutpostsOutpostInstanceTypeDataSource_InstanceType(t *testing.T) {
	dataSourceName := "data.aws_outposts_outpost_instance_type.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:   acctest.ErrorCheck(t, outposts.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSOutpostsOutpostInstanceTypeDataSourceConfigInstanceType(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "instance_type", regexp.MustCompile(`^.+$`)),
				),
			},
		},
	})
}

func TestAccAWSOutpostsOutpostInstanceTypeDataSource_PreferredInstanceTypes(t *testing.T) {
	dataSourceName := "data.aws_outposts_outpost_instance_type.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:   acctest.ErrorCheck(t, outposts.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSOutpostsOutpostInstanceTypeDataSourceConfigPreferredInstanceTypes(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "instance_type", regexp.MustCompile(`^.+$`)),
				),
			},
		},
	})
}

func testAccAWSOutpostsOutpostInstanceTypeDataSourceConfigInstanceType() string {
	return `
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost_instance_types" "test" {
  arn = tolist(data.aws_outposts_outposts.test.arns)[0]
}

data "aws_outposts_outpost_instance_type" "test" {
  arn           = tolist(data.aws_outposts_outposts.test.arns)[0]
  instance_type = tolist(data.aws_outposts_outpost_instance_types.test.instance_types)[0]
}
`
}

func testAccAWSOutpostsOutpostInstanceTypeDataSourceConfigPreferredInstanceTypes() string {
	return `
data "aws_outposts_outposts" "test" {}

data "aws_outposts_outpost_instance_types" "test" {
  arn = tolist(data.aws_outposts_outposts.test.arns)[0]
}

data "aws_outposts_outpost_instance_type" "test" {
  arn                      = tolist(data.aws_outposts_outposts.test.arns)[0]
  preferred_instance_types = data.aws_outposts_outpost_instance_types.test.instance_types
}
`
}
