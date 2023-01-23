#!/usr/bin/env sh

docker build -t wasa-photos-frontend:latest -f Dockerfile.frontend .
docker build -t wasa-photos-backend:latest -f Dockerfile.backend .