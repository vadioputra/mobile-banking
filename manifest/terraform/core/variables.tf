variable "project_id" {}
variable "region" {
  default = "asia-southeast2"
}

######### main branch #########  
variable "service_name" {
  default = "mobile-banking-v3"
}
variable "image" {
    default = "fitrakz/mobile-banking-v3:e5da5f98"
}

######### development branch #########  
variable "service_name_development" {
  default = "mobile-banking-v3-development"
}
variable "image_development" {
    default = "fitrakz/mobile-banking-v3:e5da5f98"
}
