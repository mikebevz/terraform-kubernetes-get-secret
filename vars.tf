variable "name" {
  type        = string
  description = "Secret name"
}

variable "key" {
  type        = string
  description = "Key to fetch"
}

variable "namespace" {
  type        = string
  default     = "default"
  description = "Namespace"
}

variable "context" {
  type        = string
  description = "Kubernetes context"
  default     = "default"
}

variable "authinfo" {
  type = object({
    host : string,
    cluster_ca_certificate : string
    token : string
  })
}
