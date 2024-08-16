terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "privateDnsZone" {
  type                      = "Microsoft.Network/privateDnsZones@2018-09-01"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = "${var.resource_name}.com"
  location                  = "global"
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "MX" {
  type      = "Microsoft.Network/privateDnsZones/MX@2018-09-01"
  parent_id = azapi_resource.privateDnsZone.id
  name      = var.resource_name
  body = {
    properties = {
      metadata = {
      }
      mxRecords = [
        {
          exchange   = "mx1.contoso.com"
          preference = 10
        },
        {
          exchange   = "mx2.contoso.com"
          preference = 10
        },
      ]
      ttl = 300
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

