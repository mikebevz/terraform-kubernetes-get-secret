data "external" "secret" {
  program = ["${path.module}/bin/get-secret-value"]

  query = {
    namespace              = "${var.namespace}"
    name                   = "${var.name}"
    key                    = "${var.key}"
    host                   = var.authinfo.host
    cluster_ca_certificate = var.authinfo.cluster_ca_certificate
    token                  = var.authinfo.token
  }
}
