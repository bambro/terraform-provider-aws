package aws

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/service/eks/finder"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
)

func init() {
	resource.AddTestSweepers("aws_eks_cluster", &resource.Sweeper{
		Name: "aws_eks_cluster",
		F:    testSweepEksClusters,
		Dependencies: []string{
			"aws_eks_addon",
			"aws_eks_fargate_profile",
			"aws_eks_node_group",
		},
	})
}

func testSweepEksClusters(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*conns.AWSClient).EKSConn
	input := &eks.ListClustersInput{}
	sweepResources := make([]*sweep.SweepResource, 0)

	err = conn.ListClustersPages(input, func(page *eks.ListClustersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, cluster := range page.Clusters {
			r := ResourceCluster()
			d := r.Data(nil)
			d.SetId(aws.StringValue(cluster))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EKS Clusters sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing EKS Clusters (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping EKS Clusters (%s): %w", region, err)
	}

	return nil
}

func TestAccAWSEksCluster_basic(t *testing.T) {
	var cluster eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_Required(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "eks", regexp.MustCompile(fmt.Sprintf("cluster/%s$", rName))),
					resource.TestCheckResourceAttr(resourceName, "certificate_authority.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_authority.0.data"),
					resource.TestMatchResourceAttr(resourceName, "endpoint", regexp.MustCompile(`^https://`)),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identity.0.oidc.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "identity.0.oidc.0.issuer", regexp.MustCompile(`^https://`)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_network_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "kubernetes_network_config.0.service_ipv4_cidr"),
					resource.TestMatchResourceAttr(resourceName, "platform_version", regexp.MustCompile(`^eks\.\d+$`)),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "status", eks.ClusterStatusActive),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestMatchResourceAttr(resourceName, "version", regexp.MustCompile(`^\d+\.\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.endpoint_private_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.endpoint_public_access", "true"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.security_group_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.subnet_ids.#", "2"),
					resource.TestMatchResourceAttr(resourceName, "vpc_config.0.vpc_id", regexp.MustCompile(`^vpc-.+`)),
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

func TestAccAWSEksCluster_disappears(t *testing.T) {
	var cluster eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_Required(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceCluster(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSEksCluster_EncryptionConfig_Create(t *testing.T) {
	var cluster eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"
	kmsKeyResourceName := "aws_kms_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_EncryptionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.0.provider.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "encryption_config.0.provider.0.key_arn", kmsKeyResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.0.resources.#", "1"),
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

func TestAccAWSEksCluster_EncryptionConfig_Update(t *testing.T) {
	var cluster1, cluster2 eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"
	kmsKeyResourceName := "aws_kms_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_Required(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster1),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.#", "0"),
				),
			},
			{
				Config: testAccAWSEksClusterConfig_EncryptionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster2),
					testAccCheckAWSEksClusterNotRecreated(&cluster1, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.0.provider.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "encryption_config.0.provider.0.key_arn", kmsKeyResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.0.resources.#", "1"),
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

// https://github.com/hashicorp/terraform-provider-aws/issues/19968.
func TestAccAWSEksCluster_EncryptionConfig_VersionUpdate(t *testing.T) {
	var cluster1, cluster2 eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"
	kmsKeyResourceName := "aws_kms_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_EncryptionConfig_Version(rName, "1.19"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster1),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.0.provider.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "encryption_config.0.provider.0.key_arn", kmsKeyResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.0.resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.19"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSEksClusterConfig_EncryptionConfig_Version(rName, "1.20"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster2),
					testAccCheckAWSEksClusterNotRecreated(&cluster1, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.0.provider.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "encryption_config.0.provider.0.key_arn", kmsKeyResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "encryption_config.0.resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.20"),
				),
			},
		},
	})
}

func TestAccAWSEksCluster_Version(t *testing.T) {
	var cluster1, cluster2 eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_Version(rName, "1.16"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster1),
					resource.TestCheckResourceAttr(resourceName, "version", "1.16"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSEksClusterConfig_Version(rName, "1.17"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster2),
					testAccCheckAWSEksClusterNotRecreated(&cluster1, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "version", "1.17"),
				),
			},
		},
	})
}

func TestAccAWSEksCluster_Logging(t *testing.T) {
	var cluster1, cluster2 eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_Logging(rName, []string{"api"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster1),
					resource.TestCheckResourceAttr(resourceName, "enabled_cluster_log_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "enabled_cluster_log_types.*", "api"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSEksClusterConfig_Logging(rName, []string{"api", "audit"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster2),
					testAccCheckAWSEksClusterNotRecreated(&cluster1, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "enabled_cluster_log_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "enabled_cluster_log_types.*", "api"),
					resource.TestCheckTypeSetElemAttr(resourceName, "enabled_cluster_log_types.*", "audit"),
				),
			},
			// Disable all log types.
			{
				Config: testAccAWSEksClusterConfig_Required(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster2),
					testAccCheckAWSEksClusterNotRecreated(&cluster1, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "enabled_cluster_log_types.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSEksCluster_Tags(t *testing.T) {
	var cluster1, cluster2, cluster3 eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfigTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster1),
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
				Config: testAccAWSEksClusterConfigTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAWSEksClusterConfigTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccAWSEksCluster_VpcConfig_SecurityGroupIds(t *testing.T) {
	var cluster eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_VpcConfig_SecurityGroupIds(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.security_group_ids.#", "1"),
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

func TestAccAWSEksCluster_VpcConfig_EndpointPrivateAccess(t *testing.T) {
	var cluster1, cluster2, cluster3 eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_VpcConfig_EndpointPrivateAccess(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster1),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.endpoint_private_access", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSEksClusterConfig_VpcConfig_EndpointPrivateAccess(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster2),
					testAccCheckAWSEksClusterNotRecreated(&cluster1, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.endpoint_private_access", "false"),
				),
			},
			{
				Config: testAccAWSEksClusterConfig_VpcConfig_EndpointPrivateAccess(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster3),
					testAccCheckAWSEksClusterNotRecreated(&cluster2, &cluster3),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.endpoint_private_access", "true"),
				),
			},
		},
	})
}

func TestAccAWSEksCluster_VpcConfig_EndpointPublicAccess(t *testing.T) {
	var cluster1, cluster2, cluster3 eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_VpcConfig_EndpointPublicAccess(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster1),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.endpoint_public_access", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSEksClusterConfig_VpcConfig_EndpointPublicAccess(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster2),
					testAccCheckAWSEksClusterNotRecreated(&cluster1, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.endpoint_public_access", "true"),
				),
			},
			{
				Config: testAccAWSEksClusterConfig_VpcConfig_EndpointPublicAccess(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster3),
					testAccCheckAWSEksClusterNotRecreated(&cluster2, &cluster3),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.endpoint_public_access", "false"),
				),
			},
		},
	})
}

func TestAccAWSEksCluster_VpcConfig_PublicAccessCidrs(t *testing.T) {
	var cluster eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEksClusterConfig_VpcConfig_PublicAccessCidrs(rName, `["1.2.3.4/32", "5.6.7.8/32"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.public_access_cidrs.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSEksClusterConfig_VpcConfig_PublicAccessCidrs(rName, `["4.3.2.1/32", "8.7.6.5/32"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.public_access_cidrs.#", "2"),
				),
			},
		},
	})
}

func TestAccAWSEksCluster_NetworkConfig_ServiceIpv4Cidr(t *testing.T) {
	var cluster1, cluster2 eks.Cluster
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_eks_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEks(t) },
		ErrorCheck:   acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSEksClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAWSEksClusterConfig_NetworkConfig_ServiceIpv4Cidr(rName, `"10.0.0.0/11"`),
				ExpectError: regexp.MustCompile(`expected .* to contain a network Value with between`),
			},
			{
				Config:      testAccAWSEksClusterConfig_NetworkConfig_ServiceIpv4Cidr(rName, `"10.0.0.0/25"`),
				ExpectError: regexp.MustCompile(`expected .* to contain a network Value with between`),
			},
			{
				Config:      testAccAWSEksClusterConfig_NetworkConfig_ServiceIpv4Cidr(rName, `"9.0.0.0/16"`),
				ExpectError: regexp.MustCompile(`must be within`),
			},
			{
				Config:      testAccAWSEksClusterConfig_NetworkConfig_ServiceIpv4Cidr(rName, `"172.14.0.0/24"`),
				ExpectError: regexp.MustCompile(`must be within`),
			},
			{
				Config:      testAccAWSEksClusterConfig_NetworkConfig_ServiceIpv4Cidr(rName, `"192.167.0.0/24"`),
				ExpectError: regexp.MustCompile(`must be within`),
			},
			{
				Config: testAccAWSEksClusterConfig_NetworkConfig_ServiceIpv4Cidr(rName, `"192.168.0.0/24"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster1),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_network_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_network_config.0.service_ipv4_cidr", "192.168.0.0/24"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccAWSEksClusterConfig_NetworkConfig_ServiceIpv4Cidr(rName, `"192.168.0.0/24"`),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccAWSEksClusterConfig_NetworkConfig_ServiceIpv4Cidr(rName, `"192.168.1.0/24"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEksClusterExists(resourceName, &cluster2),
					testAccCheckAWSEksClusterRecreated(&cluster1, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_network_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_network_config.0.service_ipv4_cidr", "192.168.1.0/24"),
				),
			},
		},
	})
}

func testAccCheckAWSEksClusterExists(resourceName string, cluster *eks.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EKS Cluster ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EKSConn

		output, err := tfeks.FindClusterByName(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*cluster = *output

		return nil
	}
}

func testAccCheckAWSEksClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_eks_cluster" {
			continue
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EKSConn

		_, err := tfeks.FindClusterByName(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EKS Cluster %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckAWSEksClusterRecreated(i, j *eks.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.TimeValue(i.CreatedAt).Equal(aws.TimeValue(j.CreatedAt)) {
			return errors.New("EKS Cluster was not recreated")
		}

		return nil
	}
}

func testAccCheckAWSEksClusterNotRecreated(i, j *eks.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !aws.TimeValue(i.CreatedAt).Equal(aws.TimeValue(j.CreatedAt)) {
			return errors.New("EKS Cluster was recreated")
		}

		return nil
	}
}

func testAccPreCheckAWSEks(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EKSConn

	input := &eks.ListClustersInput{}

	_, err := conn.ListClusters(input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccAWSEksClusterConfig_Base(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.${data.aws_partition.current.dns_suffix}"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "test-AmazonEKSClusterPolicy" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.test.name
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name                          = %[1]q
    "kubernetes.io/cluster/%[1]s" = "shared"
  }
}

resource "aws_subnet" "test" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name                          = %[1]q
    "kubernetes.io/cluster/%[1]s" = "shared"
  }
}
`, rName)
}

func testAccAWSEksClusterConfig_Required(rName string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName))
}

func testAccAWSEksClusterConfig_Version(rName, version string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn
  version  = %[2]q

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName, version))
}

func testAccAWSEksClusterConfig_Logging(rName string, logTypes []string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name                      = %[1]q
  role_arn                  = aws_iam_role.test.arn
  enabled_cluster_log_types = ["%[2]v"]

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName, strings.Join(logTypes, "\", \"")))
}

func testAccAWSEksClusterConfigTags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  tags = {
    %[2]q = %[3]q
  }

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName, tagKey1, tagValue1))
}

func testAccAWSEksClusterConfigTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccAWSEksClusterConfig_EncryptionConfig(rName string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  encryption_config {
    resources = ["secrets"]

    provider {
      key_arn = aws_kms_key.test.arn
    }
  }

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName))
}

func testAccAWSEksClusterConfig_EncryptionConfig_Version(rName, version string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn
  version  = %[2]q

  encryption_config {
    resources = ["secrets"]

    provider {
      key_arn = aws_kms_key.test.arn
    }
  }

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName, version))
}

func testAccAWSEksClusterConfig_VpcConfig_SecurityGroupIds(rName string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_security_group" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  vpc_config {
    security_group_ids = [aws_security_group.test.id]
    subnet_ids         = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName))
}

func testAccAWSEksClusterConfig_VpcConfig_EndpointPrivateAccess(rName string, endpointPrivateAccess bool) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  vpc_config {
    endpoint_private_access = %[2]t
    endpoint_public_access  = true
    subnet_ids              = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName, endpointPrivateAccess))
}

func testAccAWSEksClusterConfig_VpcConfig_EndpointPublicAccess(rName string, endpointPublicAccess bool) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  vpc_config {
    endpoint_private_access = true
    endpoint_public_access  = %[2]t
    subnet_ids              = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName, endpointPublicAccess))
}

func testAccAWSEksClusterConfig_VpcConfig_PublicAccessCidrs(rName string, publicAccessCidr string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  vpc_config {
    endpoint_private_access = true
    endpoint_public_access  = true
    public_access_cidrs     = %[2]s
    subnet_ids              = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName, publicAccessCidr))
}

func testAccAWSEksClusterConfig_NetworkConfig_ServiceIpv4Cidr(rName string, serviceIpv4Cidr string) string {
	return acctest.ConfigCompose(testAccAWSEksClusterConfig_Base(rName), fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  kubernetes_network_config {
    service_ipv4_cidr = %[2]s
  }

  depends_on = [aws_iam_role_policy_attachment.test-AmazonEKSClusterPolicy]
}
`, rName, serviceIpv4Cidr))
}
