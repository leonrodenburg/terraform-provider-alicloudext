#!/bin/bash
PROVIDER_NAME="terraform-provider-alicloudssl"

echo "Building Terraform provider..."
go build -o $PROVIDER_NAME ./

PLUGIN_DIR="$HOME/.terraform.d/plugins"
if [ ! -d "$PLUGIN_DIR" ]; then
  echo "Creating local provider directory..."
  mkdir -p $PLUGIN_DIR
fi

cp $PROVIDER_NAME $PLUGIN_DIR/$PROVIDER_NAME
echo "Custom provider installed!"