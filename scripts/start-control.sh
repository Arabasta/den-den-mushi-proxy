#!/bin/bash
set -e

app="proxy"

sudo systemctl stop "${app}.service" || true
sudo systemctl daemon-reload
sudo systemctl enable "${app}.service"
sudo systemctl start "${app}.service"
