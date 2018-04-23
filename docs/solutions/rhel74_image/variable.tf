# DO NOT ALTER THIS FILE

provider "oci" {
  tenancy_ocid         = "${var.tenancy_ocid}"
  user_ocid            = "${var.user_ocid}"
  fingerprint          = "${var.fingerprint}"
  private_key_path     = "${var.private_key_path}"
  private_key_password = "${var.private_key_password}"
  region               = "${var.region}"
}

variable "tenancy_ocid" {}
variable "user_ocid" {}
variable "fingerprint" {}
variable "private_key_path" {}
variable "private_key_password" {}
variable "ssh_public_key" {}
variable "region" {}

variable "ipxe_instance" {
	type = "map"
	default = {
		name = "ipxe-rhel74"
		hostname = "ipxe-rhel74"
		os = "Oracle Linux"
		os-version = "7.4"
		shape = "VM.Standard1.1"
	}
}

variable "region_all_zeros_ocid" {
	type = "map"
	default = {
		us-phoenix-1   = "ocid1.image.oc1..aaaaaaaadevqufnklkexuu6z62f7riocqigz6zng5mxhuhghy3e6zurwct2a"
		us-ashburn-1   = "ocid1.image.oc1..aaaaaaaadevqufnklkexuu6z62f7riocqigz6zng5mxhuhghy3e6zurwct2a"
		eu-frankfurt-1 = "ocid1.image.oc1..aaaaaaaadevqufnklkexuu6z62f7riocqigz6zng5mxhuhghy3e6zurwct2a"
		uk-london-1    = "ocid1.image.oc1..aaaaaaaadevqufnklkexuu6z62f7riocqigz6zng5mxhuhghy3e6zurwct2a"
	}
}

variable "ipxe_image_ocid" {
	type = "map"
	default = {
		us-phoenix-1 = "ocid1.image.oc1.phx.aaaaaaaaxklzl52nmabfp3466ilzfpo7lv737k44kih4jpo7nsmxjehwrdsq"
		us-ashburn-1 = "ocid1.image.oc1.iad.aaaaaaaahglw45opiuf6zrbhyuywh7lv5nxsxqbm4yznjwkac6zkapg6zi4a"
		eu-frankfurt-1 = "ocid1.image.oc1.eu-frankfurt-1.aaaaaaaagixzcssj76xeehppbnobvhais57zrvxe32bncaalty4wrhpossfa"
	}
}
