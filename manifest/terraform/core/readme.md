# Deploy Mobile Banking Apps V3 on Cloud Run GCP Platform
Dokumentasi ini digunakan untuk melakukan installasi mobile-banking-v3 ke GCP Platform Cloud run menggunakan Terraform.
List infrastruktur :
- Cloud Run mobile-banking-v3 (branch main)
- Cloud Run mobile-banking-v3-development (branch development)
- Service account cloud-run-deployer yang digunakan sebagai credentials  update cloud run image untuk gitlab-ci runner

_Output service account cloud-run-deployer adalah file key json (cloud-run-deployer.json), credentials ini disimpan pada variabel Gitlab CI GCLOUD_SERVICE_KEY._

## Requirements
- Terraform
- gcloud cli

## Installation
- Setup variable tfvars.
```bash
cp terraform.tfvars.example terraform.tfvars
```
- Isi nilai pada file tfvars (tambahkan variabel dari variables.tf, jika nilai variabel tersebut tidak menggunakan nilai default).
- Deploy dengan terraform.

```bash

terraform init
terraform plan
terraform apply 
```



