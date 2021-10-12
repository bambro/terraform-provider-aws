package aws

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/hashicorp/go-multierror"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/service/route53resolver/finder"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func init() {
	resource.AddTestSweepers("aws_route53_resolver_firewall_domain_list", &resource.Sweeper{
		Name: "aws_route53_resolver_firewall_domain_list",
		F:    testSweepRoute53ResolverFirewallDomainLists,
		Dependencies: []string{
			"aws_route53_resolver_firewall_rule",
		},
	})
}

func testSweepRoute53ResolverFirewallDomainLists(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).Route53ResolverConn
	var sweeperErrs *multierror.Error

	err = conn.ListFirewallDomainListsPages(&route53resolver.ListFirewallDomainListsInput{}, func(page *route53resolver.ListFirewallDomainListsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, queryLogConfig := range page.FirewallDomainLists {
			id := aws.StringValue(queryLogConfig.Id)

			log.Printf("[INFO] Deleting Route53 Resolver DNS Firewall domain list: %s", id)
			r := ResourceFirewallDomainList()
			d := r.Data(nil)
			d.SetId(id)
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
		log.Printf("[WARN] Skipping Route53 Resolver DNS Firewall domain lists sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}
	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error retrieving Route53 Resolver DNS Firewall domain lists: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func TestAccAWSRoute53ResolverFirewallDomainList_basic(t *testing.T) {
	var v route53resolver.FirewallDomainList
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_route53_resolver_firewall_domain_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSRoute53Resolver(t) },
		ErrorCheck:   acctest.ErrorCheck(t, route53resolver.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckRoute53ResolverFirewallDomainListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoute53ResolverFirewallDomainListConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallDomainListExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccAWSRoute53ResolverFirewallDomainList_domains(t *testing.T) {
	var v route53resolver.FirewallDomainList
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_route53_resolver_firewall_domain_list.test"

	domainName1 := acctest.RandomFQDomainName()
	domainName2 := acctest.RandomFQDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSRoute53Resolver(t) },
		ErrorCheck:   acctest.ErrorCheck(t, route53resolver.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckRoute53ResolverFirewallDomainListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoute53ResolverFirewallDomainListConfigDomains(rName, domainName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallDomainListExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "domains.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "domains.*", domainName1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRoute53ResolverFirewallDomainListConfigDomains(rName, domainName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallDomainListExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "domains.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "domains.*", domainName2),
				),
			},
			{
				Config: testAccRoute53ResolverFirewallDomainListConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallDomainListExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "domains.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSRoute53ResolverFirewallDomainList_disappears(t *testing.T) {
	var v route53resolver.FirewallDomainList
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_route53_resolver_firewall_domain_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSRoute53Resolver(t) },
		ErrorCheck:   acctest.ErrorCheck(t, route53resolver.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckRoute53ResolverFirewallDomainListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoute53ResolverFirewallDomainListConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallDomainListExists(resourceName, &v),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceFirewallDomainList(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSRoute53ResolverFirewallDomainList_tags(t *testing.T) {
	var v route53resolver.FirewallDomainList
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_route53_resolver_firewall_domain_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSRoute53Resolver(t) },
		ErrorCheck:   acctest.ErrorCheck(t, route53resolver.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckRoute53ResolverFirewallDomainListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoute53ResolverFirewallDomainListConfigTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallDomainListExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRoute53ResolverFirewallDomainListConfigTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallDomainListExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccRoute53ResolverFirewallDomainListConfigTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ResolverFirewallDomainListExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckRoute53ResolverFirewallDomainListDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).Route53ResolverConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_route53_resolver_firewall_domain_list" {
			continue
		}

		// Try to find the resource
		_, err := finder.FirewallDomainListByID(conn, rs.Primary.ID)
		// Verify the error is what we want
		if tfawserr.ErrMessageContains(err, route53resolver.ErrCodeResourceNotFoundException, "") {
			continue
		}
		if err != nil {
			return err
		}
		return fmt.Errorf("Route 53 Resolver DNS Firewall domain list still exists: %s", rs.Primary.ID)
	}

	return nil
}

func testAccCheckRoute53ResolverFirewallDomainListExists(n string, v *route53resolver.FirewallDomainList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route 53 Resolver DNS Firewall domain list ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).Route53ResolverConn
		out, err := finder.FirewallDomainListByID(conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		*v = *out

		return nil
	}
}

func testAccRoute53ResolverFirewallDomainListConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_domain_list" "test" {
  name = %[1]q
}
`, rName)
}

func testAccRoute53ResolverFirewallDomainListConfigDomains(rName, domain string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_domain_list" "test" {
  name    = %[1]q
  domains = [%[2]q]
}
`, rName, domain)
}

func testAccRoute53ResolverFirewallDomainListConfigTags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_domain_list" "test" {
  name = %[1]q
  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccRoute53ResolverFirewallDomainListConfigTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_domain_list" "test" {
  name = %[1]q
  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}