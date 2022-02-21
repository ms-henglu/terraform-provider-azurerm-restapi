package services_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/clients"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type GenericPatchResource struct{}

func TestAccGenericPatchResource_automationAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_patch_resource", "test")
	r := GenericPatchResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.automationAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericPatchResource_withNameParentId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_patch_resource", "test")
	r := GenericPatchResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.automationAccountWithNameParentId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (r GenericPatchResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.NewResourceID(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	resp, _, err := client.NewResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		var responseErr *azcore.ResponseError
		if errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusNotFound {
			exist := false
			return &exist, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	exist := len(utils.GetId(resp)) != 0
	return &exist, nil
}

func (r GenericPatchResource) automationAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm-restapi_patch_resource" "test" {
  resource_id = azurerm_automation_account.test.id
  type        = "Microsoft.Automation/automationAccounts@2021-06-22"
  body        = <<BODY
{
  "properties": {
    "publicNetworkAccess": true
  }
}
  BODY
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (r GenericPatchResource) automationAccountWithNameParentId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm-restapi_patch_resource" "test" {
  name      = azurerm_automation_account.test.name
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  body      = <<BODY
{
  "properties": {
    "publicNetworkAccess": true
  }
}
  BODY
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (GenericPatchResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    azurerm = {
      version = "= 2.75.0"
      source  = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.LocationPrimary, data.RandomStringOfLength(10))
}
