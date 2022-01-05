#!/bin/sh
set -e

if [[ -z "${GITHUB_TOKEN}" ]]; then
  echo "GITHUB_TOKEN is required"
  exit 1
fi

git config --global credential.helper store
echo "https://github-backup:${GITHUB_TOKEN}@github.com" > ~/.git-credentials

python /app/backup.py --https -c /app/config.json -b /tmp/backups/github.com --concurrent 50
