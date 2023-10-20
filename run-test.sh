#!/bin/bash
docker compose -f docker-compose-test.yaml rm --stop --force
docker compose -f docker-compose-test.yaml up integration --build --abort-on-container-exit --exit-code-from integration
docker compose -f docker-compose-test.yaml rm --stop --force
