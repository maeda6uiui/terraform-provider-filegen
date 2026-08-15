terraform {
  required_providers {
    filegen = {
      source  = "maeda6uiui/filegen"
      version = "0.0.1-alpha1"
    }
  }

  backend "local" {
    path = "terraform.tfstate"
  }

  required_version = "~>1.9"
}

provider "filegen" {

}
