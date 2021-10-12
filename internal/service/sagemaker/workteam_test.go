package sagemaker_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/hashicorp/go-multierror"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfsagemaker "github.com/hashicorp/terraform-provider-aws/internal/service/sagemaker"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func init() {
	resource.AddTestSweepers("aws_sagemaker_workteam", &resource.Sweeper{
		Name: "aws_sagemaker_workteam",
		F:    testSweepSagemakerWorkteams,
	})
}

func testSweepSagemakerWorkteams(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).SageMakerConn
	var sweeperErrs *multierror.Error

	err = conn.ListWorkteamsPages(&sagemaker.ListWorkteamsInput{}, func(page *sagemaker.ListWorkteamsOutput, lastPage bool) bool {
		for _, workteam := range page.Workteams {

			r := ResourceWorkteam()
			d := r.Data(nil)
			d.SetId(aws.StringValue(workteam.WorkteamName))
			err := r.Delete(d, client)
			if err != nil {
				log.Printf("[ERROR] %s", err)
				sweeperErrs = multierror.Append(sweeperErrs, err)
				continue
			}
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping SageMaker workteam sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil()
	}

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error retrieving Sagemaker Workteams: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func testAccAWSSagemakerWorkteam_cognitoConfig(t *testing.T) {
	var workteam sagemaker.Workteam
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_sagemaker_workteam.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sagemaker.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSSagemakerWorkteamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerWorkteamCognitoConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", rName),
					resource.TestCheckResourceAttr(resourceName, "member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.cognito_member_definition.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.client_id", "aws_cognito_user_pool_client.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_pool", "aws_cognito_user_pool.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_group", "aws_cognito_user_group.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "subdomain"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workforce_name"},
			},
			{
				Config: testAccAWSSagemakerWorkteamCognitoUpdatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", rName),
					resource.TestCheckResourceAttr(resourceName, "member_definition.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.cognito_member_definition.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.client_id", "aws_cognito_user_pool_client.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_pool", "aws_cognito_user_pool.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_group", "aws_cognito_user_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.1.cognito_member_definition.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.1.cognito_member_definition.0.client_id", "aws_cognito_user_pool_client.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.1.cognito_member_definition.0.user_pool", "aws_cognito_user_pool.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.1.cognito_member_definition.0.user_group", "aws_cognito_user_group.test2", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "subdomain"),
				),
			},
			{
				Config: testAccAWSSagemakerWorkteamCognitoConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", rName),
					resource.TestCheckResourceAttr(resourceName, "member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.cognito_member_definition.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.client_id", "aws_cognito_user_pool_client.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_pool", "aws_cognito_user_pool.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "member_definition.0.cognito_member_definition.0.user_group", "aws_cognito_user_group.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "subdomain"),
				),
			},
		},
	})
}

func testAccAWSSagemakerWorkteam_oidcConfig(t *testing.T) {
	var workteam sagemaker.Workteam
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_sagemaker_workteam.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sagemaker.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSSagemakerWorkteamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerWorkteamOidcConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.oidc_member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.oidc_member_definition.0.groups.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "member_definition.0.oidc_member_definition.0.groups.*", rName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workforce_name"},
			},
			{
				Config: testAccAWSSagemakerWorkteamOidcConfig2(rName, "test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.oidc_member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.oidc_member_definition.0.groups.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "member_definition.0.oidc_member_definition.0.groups.*", rName),
					resource.TestCheckTypeSetElemAttr(resourceName, "member_definition.0.oidc_member_definition.0.groups.*", "test"),
				),
			},
			{
				Config: testAccAWSSagemakerWorkteamOidcConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.oidc_member_definition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "member_definition.0.oidc_member_definition.0.groups.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "member_definition.0.oidc_member_definition.0.groups.*", rName)),
			},
		},
	})
}

func testAccAWSSagemakerWorkteam_tags(t *testing.T) {
	var workteam sagemaker.Workteam
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_sagemaker_workteam.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sagemaker.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSSagemakerWorkteamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerWorkteamTagsConfig1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workforce_name"},
			},
			{
				Config: testAccAWSSagemakerWorkteamTagsConfig2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAWSSagemakerWorkteamTagsConfig1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccAWSSagemakerWorkteam_notificationConfig(t *testing.T) {
	var workteam sagemaker.Workteam
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_sagemaker_workteam.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sagemaker.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSSagemakerWorkteamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerWorkteamNotificationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", rName),
					resource.TestCheckResourceAttr(resourceName, "notification_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_configuration.0.notification_topic_arn", "aws_sns_topic.test", "arn"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workforce_name"},
			},
			{
				Config: testAccAWSSagemakerWorkteamOidcConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", rName),
					resource.TestCheckResourceAttr(resourceName, "notification_configuration.#", "1"),
				),
			},
			{
				Config: testAccAWSSagemakerWorkteamNotificationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					resource.TestCheckResourceAttr(resourceName, "workteam_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexp.MustCompile(`workteam/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", rName),
					resource.TestCheckResourceAttr(resourceName, "notification_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "notification_configuration.0.notification_topic_arn", "aws_sns_topic.test", "arn"),
				),
			},
		},
	})
}

func testAccAWSSagemakerWorkteam_disappears(t *testing.T) {
	var workteam sagemaker.Workteam
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_sagemaker_workteam.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sagemaker.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSSagemakerWorkteamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerWorkteamOidcConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerWorkteamExists(resourceName, &workteam),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceWorkteam(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAWSSagemakerWorkteamDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_sagemaker_workteam" {
			continue
		}

		_, err := tfsagemaker.FindWorkteamByName(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("SageMaker Workteam %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckAWSSagemakerWorkteamExists(n string, workteam *sagemaker.Workteam) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SageMaker Workteam ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerConn

		output, err := tfsagemaker.FindWorkteamByName(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*workteam = *output

		return nil
	}
}

func testAccAWSSagemakerWorkteamCognitoBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_cognito_user_pool" "test" {
  name = %[1]q
}

resource "aws_cognito_user_pool_client" "test" {
  name            = %[1]q
  generate_secret = true
  user_pool_id    = aws_cognito_user_pool.test.id
}

resource "aws_cognito_user_pool_domain" "test" {
  domain       = %[1]q
  user_pool_id = aws_cognito_user_pool.test.id
}

resource "aws_cognito_user_group" "test" {
  name         = %[1]q
  user_pool_id = aws_cognito_user_pool.test.id
}

resource "aws_sagemaker_workforce" "test" {
  workforce_name = %[1]q

  cognito_config {
    client_id = aws_cognito_user_pool_client.test.id
    user_pool = aws_cognito_user_pool_domain.test.user_pool_id
  }
}
`, rName)
}

func testAccAWSSagemakerWorkteamCognitoConfig(rName string) string {
	return acctest.ConfigCompose(testAccAWSSagemakerWorkteamCognitoBaseConfig(rName), fmt.Sprintf(`
resource "aws_sagemaker_workteam" "test" {
  workteam_name  = %[1]q
  workforce_name = aws_sagemaker_workforce.test.id
  description    = %[1]q

  member_definition {
    cognito_member_definition {
      client_id  = aws_cognito_user_pool_client.test.id
      user_pool  = aws_cognito_user_pool_domain.test.user_pool_id
      user_group = aws_cognito_user_group.test.id
    }
  }
}
`, rName))
}

func testAccAWSSagemakerWorkteamCognitoUpdatedConfig(rName string) string {
	return acctest.ConfigCompose(testAccAWSSagemakerWorkteamCognitoBaseConfig(rName), fmt.Sprintf(`
resource "aws_cognito_user_group" "test2" {
  name         = "%[1]s-2"
  user_pool_id = aws_cognito_user_pool.test.id
}

resource "aws_sagemaker_workteam" "test" {
  workteam_name  = %[1]q
  workforce_name = aws_sagemaker_workforce.test.id
  description    = %[1]q

  member_definition {
    cognito_member_definition {
      client_id  = aws_cognito_user_pool_client.test.id
      user_pool  = aws_cognito_user_pool_domain.test.user_pool_id
      user_group = aws_cognito_user_group.test.id
    }
  }

  member_definition {
    cognito_member_definition {
      client_id  = aws_cognito_user_pool_client.test.id
      user_pool  = aws_cognito_user_pool_domain.test.user_pool_id
      user_group = aws_cognito_user_group.test2.id
    }
  }
}
`, rName))
}

func testAccAWSSagemakerWorkteamOidcBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sagemaker_workforce" "test" {
  workforce_name = %[1]q

  oidc_config {
    authorization_endpoint = "https://example.com"
    client_id              = %[1]q
    client_secret          = %[1]q
    issuer                 = "https://example.com"
    jwks_uri               = "https://example.com"
    logout_endpoint        = "https://example.com"
    token_endpoint         = "https://example.com"
    user_info_endpoint     = "https://example.com"
  }
}
`, rName)
}

func testAccAWSSagemakerWorkteamOidcConfig(rName string) string {
	return acctest.ConfigCompose(testAccAWSSagemakerWorkteamOidcBaseConfig(rName), fmt.Sprintf(`
resource "aws_sagemaker_workteam" "test" {
  workteam_name  = %[1]q
  workforce_name = aws_sagemaker_workforce.test.id
  description    = %[1]q

  member_definition {
    oidc_member_definition {
      groups = [%[1]q]
    }
  }
}
`, rName))
}

func testAccAWSSagemakerWorkteamOidcConfig2(rName, group string) string {
	return acctest.ConfigCompose(testAccAWSSagemakerWorkteamOidcBaseConfig(rName), fmt.Sprintf(`
resource "aws_sagemaker_workteam" "test" {
  workteam_name  = %[1]q
  workforce_name = aws_sagemaker_workforce.test.id
  description    = %[1]q

  member_definition {
    oidc_member_definition {
      groups = [%[1]q, %[2]q]
    }
  }
}
`, rName, group))
}

func testAccAWSSagemakerWorkteamNotificationConfig(rName string) string {
	return acctest.ConfigCompose(testAccAWSSagemakerWorkteamOidcBaseConfig(rName), fmt.Sprintf(`
resource "aws_sns_topic" "test" {
  name = %[1]q
}

resource "aws_sns_topic_policy" "test" {
  arn = aws_sns_topic.test.arn

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Id" : "default",
    "Statement" : [
      {
        "Sid" : "%[1]s",
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "sagemaker.amazonaws.com"
        },
        "Action" : [
          "sns:Publish"
        ],
        "Resource" : "${aws_sns_topic.test.arn}"
      }
    ]
  })
}

resource "aws_sagemaker_workteam" "test" {
  workteam_name  = %[1]q
  workforce_name = aws_sagemaker_workforce.test.id
  description    = %[1]q

  member_definition {
    oidc_member_definition {
      groups = [%[1]q]
    }
  }

  notification_configuration {
    notification_topic_arn = aws_sns_topic.test.arn
  }
}
`, rName))
}

func testAccAWSSagemakerWorkteamTagsConfig1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccAWSSagemakerWorkteamOidcBaseConfig(rName), fmt.Sprintf(`
resource "aws_sagemaker_workteam" "test" {
  workteam_name  = %[1]q
  workforce_name = aws_sagemaker_workforce.test.id
  description    = %[1]q

  member_definition {
    oidc_member_definition {
      groups = [%[1]q]
    }
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccAWSSagemakerWorkteamTagsConfig2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccAWSSagemakerWorkteamOidcBaseConfig(rName), fmt.Sprintf(`
resource "aws_sagemaker_workteam" "test" {
  workteam_name  = %[1]q
  workforce_name = aws_sagemaker_workforce.test.id
  description    = %[1]q

  member_definition {
    oidc_member_definition {
      groups = [%[1]q]
    }
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}