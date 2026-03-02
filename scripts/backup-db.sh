#!/bin/bash
set -e
BACKUP_DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="backup_green_saas_${BACKUP_DATE}.sql.gz"
S3_BUCKET=${BACKUP_S3_BUCKET:-"k8s-green-backups"}

echo "💾 Backup PostgreSQL: $BACKUP_FILE"
docker compose exec -T postgres pg_dump -U green -d green_saas \
  --no-owner --no-acl --format=plain \
  | gzip > "/tmp/$BACKUP_FILE"

echo "☁️ Upload S3: s3://$S3_BUCKET/postgres/$BACKUP_FILE"
aws s3 cp "/tmp/$BACKUP_FILE" "s3://$S3_BUCKET/postgres/$BACKUP_FILE"

echo "🗑️ Nettoyage backups > 30j"
aws s3 ls "s3://$S3_BUCKET/postgres/" | \
  awk '{print $4}' | \
  while read f; do
    age=$(( ($(date +%s) - $(date -d "${f:7:8}" +%s 2>/dev/null || echo $(date +%s))) / 86400 ))
    [ "$age" -gt 30 ] && aws s3 rm "s3://$S3_BUCKET/postgres/$f" && echo "Supprimé: $f"
  done

echo "✅ Backup terminé: $BACKUP_FILE"
