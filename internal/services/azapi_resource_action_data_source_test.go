package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ActionDataSource struct{}

func TestAccActionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionDataSource_providerAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r ActionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource_action" "test" {
  type                   = "Microsoft.Automation/automationAccounts@2022-08-08"
  resource_id            = azapi_resource.test.id
  action                 = "listKeys"
  response_export_values = ["*"]
}
`, GenericResource{}.defaultTag(data))
}

func (r ActionDataSource) providerAction() string {
	return `
data "azapi_resource_action" "test" {
  type        = "Microsoft.ResourceGraph@2020-04-01-preview"
  resource_id = "/providers/Microsoft.ResourceGraph"
  action      = "resources"
  body = jsonencode({
    query = "resources| limit 1"
  })
  response_export_values = ["*"]
}
`
}
