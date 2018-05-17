provider "azurerm" {}

# Create a resource group
resource "azurerm_resource_group" "go2hal" {
  name = "go2hal"
  location = "North Europe"
}

resource "azurerm_cosmosdb_account" "go2hal" {
  name = "go2hal"
  kind = "MongoDB"
  offer_type = "Standard"
  location = "${azurerm_resource_group.go2hal.location}"
  resource_group_name = "${azurerm_resource_group.go2hal.name}"
  consistency_policy {
    consistency_level = "BoundedStaleness"
  }
  geo_location {
    location = "${azurerm_resource_group.go2hal.location}"
    failover_priority = 0
  }

}

output "cosmosdb_account_endpoint" {
  value = "${azurerm_cosmosdb_account.go2hal.connection_strings}"
}

resource "azurerm_container_group" "aci-helloworld" {

  "container" {
    cpu = 1
    image = "weautomateeverything/go2hal:1.807.1"
    memory = 0.5
    name = "go2hal"
    port = "8000"
    environment_variables {
      MONGO = "${azurerm_cosmosdb_account.go2hal.connection_strings[0]}"
      BOT_KEY = "411872276:AAHeaOcCauxP0p7vEoTnl1Jeafil9fulrz0"
    }
  }
  location = "${azurerm_resource_group.go2hal.location}"
  name = "go2hal"
  os_type = "Linux"
  resource_group_name = "${azurerm_resource_group.go2hal.name}"
  dns_name_label = "go2hal"
  ip_address_type = "Public"
}


