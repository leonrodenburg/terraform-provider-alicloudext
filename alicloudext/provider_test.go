package alicloudext

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

var client = &sdk.Client{}

func createMockClient(regionId string, accessKeyId string, accessKeySecret string) (*sdk.Client, error) {
	return client, nil
}

func Test_Provider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}