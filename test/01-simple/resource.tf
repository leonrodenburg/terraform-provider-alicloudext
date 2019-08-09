provider "alicloudssl" {}

resource "alicloudssl_certificate" "my-certificate" {
  domain = "test.example.com"
}