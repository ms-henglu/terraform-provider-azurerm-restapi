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
  default = "eastus2euap"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "communicationsGateway" {
  type      = "Microsoft.VoiceServices/communicationsGateways@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location

  body = {
    properties = {
      autoGeneratedDomainNameLabelScope = "NoReuse"
      codecs = [
        "PCMA",
      ]
      connectivity = "PublicAddress"
      e911Type     = "Standard"
      platforms = [
        "OperatorConnect",
      ]
      serviceLocations = [
        {
          name = "useast"
          primaryRegionProperties = {
            allowedMediaSourceAddressPrefixes = [
              "10.1.2.0/24",
            ]
            allowedSignalingSourceAddressPrefixes = [
              "10.1.1.0/24",
            ]
            operatorAddresses = [
              "198.51.100.1",
            ]
          }
        },
        {
          name = "useast2"
          primaryRegionProperties = {
            allowedMediaSourceAddressPrefixes = [
              "10.2.2.0/24",
            ]
            allowedSignalingSourceAddressPrefixes = [
              "10.2.1.0/24",
            ]
            operatorAddresses = [
              "198.51.100.2",
            ]
          }
        },
      ]
      teamsVoicemailPilotNumber = "1234567890"
    }
  }
}

resource "azapi_resource" "TestLine" {
  type = "Microsoft.VoiceServices/communicationsGateways/testLines@2023-09-01"
  parent_id = azapi_resource.communicationsGateway.id
  name      = var.resource_name
  location  = var.location

  body = {
    properties = {
      phoneNumber = "123456789"
      purpose     = "Automated"
    }
  }
}
