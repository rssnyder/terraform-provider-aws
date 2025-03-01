package vpclattice_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfvpclattice "github.com/hashicorp/terraform-provider-aws/internal/service/vpclattice"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccVPCLatticeTargetGroup_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var targetGroup vpclattice.GetTargetGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_vpclattice_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.VPCLatticeEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.VPCLatticeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "vpc-lattice", regexp.MustCompile("targetgroup/.+$")),
					resource.TestCheckResourceAttr(resourceName, "config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.health_check_interval_seconds", "30"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.health_check_timeout_seconds", "5"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.healthy_threshold_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.matcher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.matcher.0.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.port", "0"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.protocol_version", "HTTP1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.unhealthy_threshold_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "config.0.ip_address_type", ""),
					resource.TestCheckResourceAttr(resourceName, "config.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol_version", "HTTP1"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", "INSTANCE"),
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

func TestAccVPCLatticeTargetGroup_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var targetGroup vpclattice.GetTargetGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_vpclattice_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.VPCLatticeEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.VPCLatticeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckServiceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfvpclattice.ResourceTargetGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCLatticeTargetGroup_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var targetGroup vpclattice.GetTargetGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_vpclattice_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.VPCLatticeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
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
				Config: testAccTargetGroupConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccTargetGroupConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVPCLatticeTargetGroup_lambda(t *testing.T) {
	ctx := acctest.Context(t)
	var targetGroup vpclattice.GetTargetGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_vpclattice_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.VPCLatticeEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.VPCLatticeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupConfig_lambda(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "vpc-lattice", regexp.MustCompile("targetgroup/.+$")),
					resource.TestCheckResourceAttr(resourceName, "config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", "LAMBDA"),
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

func TestAccVPCLatticeTargetGroup_ip(t *testing.T) {
	ctx := acctest.Context(t)
	var targetGroup vpclattice.GetTargetGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_vpclattice_target_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.VPCLatticeEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.VPCLatticeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupConfig_ip(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "vpc-lattice", regexp.MustCompile("targetgroup/.+$")),
					resource.TestCheckResourceAttr(resourceName, "config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.health_check_interval_seconds", "60"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.health_check_timeout_seconds", "10"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.healthy_threshold_count", "6"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.matcher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.matcher.0.value", "200-299"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.path", "/health"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.protocol_version", "HTTP1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.unhealthy_threshold_count", "4"),
					resource.TestCheckResourceAttr(resourceName, "config.0.ip_address_type", "IPV6"),
					resource.TestCheckResourceAttr(resourceName, "config.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol_version", "HTTP2"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTargetGroupConfig_ipUpdated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "vpc-lattice", regexp.MustCompile("targetgroup/.+$")),
					resource.TestCheckResourceAttr(resourceName, "config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.health_check_interval_seconds", "180"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.health_check_timeout_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.healthy_threshold_count", "8"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.matcher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.matcher.0.value", "202"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.path", "/health"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.protocol_version", "HTTP2"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.0.unhealthy_threshold_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "config.0.ip_address_type", "IPV6"),
					resource.TestCheckResourceAttr(resourceName, "config.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol_version", "HTTP2"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
				),
			},
		},
	})
}

func TestAccVPCLatticeTargetGroup_alb(t *testing.T) {
	ctx := acctest.Context(t)
	var targetGroup vpclattice.GetTargetGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_vpclattice_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.VPCLatticeEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.VPCLatticeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupConfig_alb(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "vpc-lattice", regexp.MustCompile("targetgroup/.+$")),
					resource.TestCheckResourceAttr(resourceName, "config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "config.0.health_check.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "config.0.ip_address_type", ""),
					resource.TestCheckResourceAttr(resourceName, "config.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol_version", "HTTP1"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", "ALB"),
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

func testAccCheckTargetGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).VPCLatticeClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_vpclattice_target_group" {
				continue
			}

			_, err := conn.GetTargetGroup(ctx, &vpclattice.GetTargetGroupInput{
				TargetGroupIdentifier: aws.String(rs.Primary.ID),
			})
			if err != nil {
				var nfe *types.ResourceNotFoundException
				if errors.As(err, &nfe) {
					return nil
				}
				return err
			}

			return create.Error(names.VPCLattice, create.ErrActionCheckingDestroyed, tfvpclattice.ResNameService, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckTargetGroupExists(ctx context.Context, name string, targetGroup *vpclattice.GetTargetGroupOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.VPCLattice, create.ErrActionCheckingExistence, tfvpclattice.ResNameService, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.VPCLattice, create.ErrActionCheckingExistence, tfvpclattice.ResNameService, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).VPCLatticeClient(ctx)
		resp, err := conn.GetTargetGroup(ctx, &vpclattice.GetTargetGroupInput{
			TargetGroupIdentifier: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return create.Error(names.VPCLattice, create.ErrActionCheckingExistence, tfvpclattice.ResNameService, rs.Primary.ID, err)
		}

		*targetGroup = *resp

		return nil
	}
}

func testAccTargetGroupConfig_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 0), fmt.Sprintf(`
resource "aws_vpclattice_target_group" "test" {
  name = %[1]q
  type = "INSTANCE"

  config {
    port           = 80
    protocol       = "HTTP"
    vpc_identifier = aws_vpc.test.id
  }
}
`, rName))
}

func testAccTargetGroupConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_vpclattice_target_group" "test" {
  name = %[1]q
  type = "LAMBDA"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccTargetGroupConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_vpclattice_target_group" "test" {
  name = %[1]q
  type = "LAMBDA"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccTargetGroupConfig_lambda(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpclattice_target_group" "test" {
  name = %[1]q
  type = "LAMBDA"
}
`, rName)
}

func testAccTargetGroupConfig_ip(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 0), fmt.Sprintf(`
resource "aws_vpclattice_target_group" "test" {
  name = %[1]q
  type = "IP"

  config {
    port             = 443
    protocol         = "HTTPS"
    vpc_identifier   = aws_vpc.test.id
    ip_address_type  = "IPV6"
    protocol_version = "HTTP2"

    health_check {
      health_check_interval_seconds = 60
      health_check_timeout_seconds  = 10
      healthy_threshold_count       = 6
      unhealthy_threshold_count     = 4

      matcher {
        value = "200-299"
      }

      path             = "/health"
      port             = 8443
      protocol         = "HTTPS"
      protocol_version = "HTTP1"
    }
  }
}
`, rName))
}

func testAccTargetGroupConfig_ipUpdated(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 0), fmt.Sprintf(`
resource "aws_vpclattice_target_group" "test" {
  name = %[1]q
  type = "IP"

  config {
    port             = 443
    protocol         = "HTTPS"
    vpc_identifier   = aws_vpc.test.id
    ip_address_type  = "IPV6"
    protocol_version = "HTTP2"

    health_check {
      health_check_interval_seconds = 180
      health_check_timeout_seconds  = 90
      healthy_threshold_count       = 8
      unhealthy_threshold_count     = 3

      matcher {
        value = "202"
      }

      path             = "/health"
      port             = 8443
      protocol         = "HTTPS"
      protocol_version = "HTTP2"
    }
  }
}
`, rName))
}

func testAccTargetGroupConfig_alb(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 0), fmt.Sprintf(`
resource "aws_vpclattice_target_group" "test" {
  name = %[1]q
  type = "ALB"

  config {
    port           = 80
    protocol       = "HTTP"
    vpc_identifier = aws_vpc.test.id
  }
}
`, rName))
}
