package transfer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transfer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/service/transfer/finder"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	tftransfer "github.com/hashicorp/terraform-provider-aws/internal/service/transfer"
	tftransfer "github.com/hashicorp/terraform-provider-aws/internal/service/transfer"
	tftransfer "github.com/hashicorp/terraform-provider-aws/internal/service/transfer"
)

const (
	userStateExists = "exists"
)

func statusServerState(conn *transfer.Transfer, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := tftransfer.FindServerByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.State), nil
	}
}

func statusUserState(conn *transfer.Transfer, serverID, userName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := tftransfer.FindUserByServerIDAndUserName(conn, serverID, userName)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, userStateExists, nil
	}
}