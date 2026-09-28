# S3-compatible storage

Development uses MinIO. Production can use Amazon S3, Cloudflare R2, Backblaze B2 S3 API, or another S3-compatible provider.

Keep buckets private. The backend should issue short-lived signed URLs for uploads/downloads. Do not put S3 credentials into the web or Android client.
