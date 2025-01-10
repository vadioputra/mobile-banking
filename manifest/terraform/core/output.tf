output "service_account_key" {
  value       = google_service_account_key.cloud_run_deployer_key.private_key
  sensitive   = true
  description = "Private key for the cloud-run-deployer service account."
}

output "service_account_key_file" {
  value       = "cloud-run-deployer.json"
  description = "Path to the service account key file."
}

resource "local_file" "cloud_run_deployer_key_file" {
  filename   = "cloud-run-deployer.json"
  content    = base64decode(google_service_account_key.cloud_run_deployer_key.private_key)
  depends_on = [google_service_account_key.cloud_run_deployer_key]
}