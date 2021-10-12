package securityhub_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfsecurityhub "github.com/hashicorp/terraform-provider-aws/internal/service/securityhub"
)

func testAccAWSSecurityHubMember_basic(t *testing.T) {
	var member securityhub.Member
	resourceName := "aws_securityhub_member.example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, securityhub.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSSecurityHubMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSecurityHubMemberConfig_basic("111111111111", acctest.DefaultEmailAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSecurityHubMemberExists(resourceName, &member),
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

func testAccAWSSecurityHubMember_invite(t *testing.T) {
	var member securityhub.Member
	resourceName := "aws_securityhub_member.example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, securityhub.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSSecurityHubMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSecurityHubMemberConfig_invite("111111111111", acctest.DefaultEmailAddress, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSecurityHubMemberExists(resourceName, &member),
					resource.TestCheckResourceAttr(resourceName, "member_status", "Invited"),
					resource.TestCheckResourceAttr(resourceName, "invite", "true"),
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

func testAccCheckAWSSecurityHubMemberExists(n string, member *securityhub.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SecurityHubConn

		resp, err := conn.GetMembers(&securityhub.GetMembersInput{
			AccountIds: []*string{aws.String(rs.Primary.ID)},
		})

		if err != nil {
			return err
		}

		if len(resp.Members) == 0 {
			return fmt.Errorf("Security Hub member %s not found", rs.Primary.ID)
		}

		member = resp.Members[0]

		return nil
	}
}

func testAccCheckAWSSecurityHubMemberDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SecurityHubConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_securityhub_member" {
			continue
		}

		resp, err := conn.GetMembers(&securityhub.GetMembersInput{
			AccountIds: []*string{aws.String(rs.Primary.ID)},
		})

		if tfawserr.ErrCodeEquals(err, tfsecurityhub.ErrCodeBadRequestException) {
			continue
		}

		if tfawserr.ErrCodeEquals(err, securityhub.ErrCodeResourceNotFoundException) {
			continue
		}

		if err != nil {
			return fmt.Errorf("error getting Security Hub Member (%s): %w", rs.Primary.ID, err)
		}

		if len(resp.Members) != 0 {
			return fmt.Errorf("Security Hub member still exists")
		}

		return nil
	}

	return nil
}

func testAccAWSSecurityHubMemberConfig_basic(accountId, email string) string {
	return fmt.Sprintf(`
resource "aws_securityhub_account" "example" {}

resource "aws_securityhub_member" "example" {
  depends_on = [aws_securityhub_account.example]
  account_id = %[1]q
  email      = %[2]q
}
`, accountId, email)
}

func testAccAWSSecurityHubMemberConfig_invite(accountId, email string, invite bool) string {
	return fmt.Sprintf(`
resource "aws_securityhub_account" "example" {}

resource "aws_securityhub_member" "example" {
  depends_on = [aws_securityhub_account.example]
  account_id = %[1]q
  email      = %[2]q
  invite     = %[3]t
}
`, accountId, email, invite)
}