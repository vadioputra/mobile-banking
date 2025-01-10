# Mobile Banking Apps V3
Aplikasi dummy mobile app.

## Requirements
- Golang

## Installation
For apps
```bash
go mod download
go run main.go
```

## ENV GITHUB
Tambahkan pada secret github. Env ini digunakan untuk melakukan sinkronisasi.
```bash
GITLAB_URL                  => URL Gitlab
GITLAB_USERNAME             => Username Gitlab
GITLAB_PAT                  => Gitlab Personal Access Token
GITLAB_PROJECT_ID           => Project ID Gitlab
```


## ENV Gitlab CI
Tambahkan pada variable gitlab ci
```bash
DOCKERHUB_USERNAME          => Image registry account username
DOCKERHUB_PASSWORD          => Image registry account password
DOCKER_IMAGE                => Image name
SERVICE_NAME                => Name cloud run for main branch
SERVICE_NAME_DEVELOPMENT    => Name cloud run for development branch
GCLOUD_SERVICE_KEY          => Credentials gcloud
GCP_PROJECT_ID              => Id GCP Project
```
_GCLOUD_SERVICE_KEY credentials ini didapat dari manifest/terraform/core/cloud-run-deployer.json setelah deploy infratructure cloud run ([link](https://gitlab.com/fitraelbi/mobile-banking-v3/-/tree/main/manifest/terraform/core?ref_type=heads))._
