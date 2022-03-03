terraform {
  required_providers {
    hashicups = {
      versions = ["0.3.0"]
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

provider "hashicups" {
  username = "education"
  password = "test123"
}