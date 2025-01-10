

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_cloud_run_service" "mobile-banking-v3" {
  name     = var.service_name
  location = var.region

  template {
    spec {
      containers {
        image = var.image
        ports {
          container_port = 8080
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

}

resource "google_cloud_run_service" "mobile-banking-v3-development" {
  name     = var.service_name_development
  location = var.region

  template {
    spec {
      containers {
        image = var.image_development
        ports {
          container_port = 8080
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "mobile-banking-v3" {
  location    = google_cloud_run_service.mobile-banking-v3.location
  project     = google_cloud_run_service.mobile-banking-v3.project
  service     = google_cloud_run_service.mobile-banking-v3.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

resource "google_cloud_run_service_iam_policy" "mobile-banking-v3-development" {
  location    = google_cloud_run_service.mobile-banking-v3-development.location
  project     = google_cloud_run_service.mobile-banking-v3-development.project
  service     = google_cloud_run_service.mobile-banking-v3-development.name

  policy_data = data.google_iam_policy.noauth.policy_data
}



resource "google_service_account" "cloud_run_deployer" {
  account_id   = "cloud-run-deployer"
  display_name = "Cloud Run Deployer Service Account"
}

resource "google_project_iam_member" "cloud_run_admin_binding" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_service_account.cloud_run_deployer.email}"
}
resource "google_project_iam_member" "service_account_user_binding" {
  project = var.project_id 
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.cloud_run_deployer.email}"
}


resource "google_service_account_key" "cloud_run_deployer_key" {
  service_account_id = google_service_account.cloud_run_deployer.name
}


