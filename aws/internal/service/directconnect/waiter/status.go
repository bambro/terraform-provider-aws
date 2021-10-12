package waiter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/service/directconnect/finder"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
)

func statusConnectionState(conn *directconnect.DirectConnect, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := tfdirectconnect.FindConnectionByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.ConnectionState), nil
	}
}

func statusGatewayState(conn *directconnect.DirectConnect, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := tfdirectconnect.FindGatewayByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.DirectConnectGatewayState), nil
	}
}

func statusGatewayAssociationState(conn *directconnect.DirectConnect, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := tfdirectconnect.FindGatewayAssociationByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.AssociationState), nil
	}
}

func statusHostedConnectionState(conn *directconnect.DirectConnect, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := tfdirectconnect.FindHostedConnectionByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.ConnectionState), nil
	}
}

func statusLagState(conn *directconnect.DirectConnect, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := tfdirectconnect.FindLagByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.LagState), nil
	}
}
