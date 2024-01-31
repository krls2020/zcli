package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeropsio/zcli/src/cmdBuilder"
	"github.com/zeropsio/zcli/src/i18n"
	"github.com/zeropsio/zerops-go/dto/input/path"
	"github.com/zeropsio/zerops-go/types"
	"github.com/zeropsio/zerops-go/types/enum"
)

func bucketZeropsDeleteCmd() *cmdBuilder.Cmd {
	return cmdBuilder.NewCmd().
		Use("delete").
		Short(i18n.T(i18n.CmdBucketDelete)).
		ScopeLevel(cmdBuilder.Service).
		Arg("bucketName").
		BoolFlag("confirm", false, i18n.T(i18n.ConfirmFlag)).
		LoggedUserRunFunc(func(ctx context.Context, cmdData *cmdBuilder.LoggedUserCmdData) error {
			uxBlocks := cmdData.UxBlocks

			if cmdData.Service.ServiceTypeCategory != enum.ServiceStackTypeCategoryEnumObjectStorage {
				return errors.New(i18n.T(i18n.BucketGenericOnlyForObjectStorage))
			}

			serviceId := cmdData.Service.ID
			// TODO - janhajek duplicate
			bucketName := fmt.Sprintf("%s.%s", strings.ToLower(serviceId.Native()), cmdData.Args["bucketName"][0])

			if !cmdData.Params.GetBool("confirm") {
				err := YesNoPromptDestructive(ctx, cmdData, i18n.T(i18n.BucketDeleteConfirm, bucketName))
				if err != nil {
					return err
				}
			}

			uxBlocks.PrintLine(i18n.T(i18n.BucketDeleteDeletingZeropsApi, bucketName))
			uxBlocks.PrintLine(i18n.T(i18n.BucketGenericBucketNamePrefixed))

			resp, err := cmdData.RestApiClient.DeleteS3(
				ctx,
				path.S3Bucket{
					ServiceStackId: serviceId,
					Name:           types.NewString(bucketName),
				},
			)
			if err != nil {
				return err
			}
			if _, err := resp.Output(); err != nil {
				return err
			}

			uxBlocks.PrintSuccessLine(i18n.T(i18n.BucketDeleted))

			return nil
		})
}
