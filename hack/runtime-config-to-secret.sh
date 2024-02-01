#!/usr/bin/env bash
#
# runtime-config-to-secret: Convert a runtime-config.toml file to a Kubernetes secret.
# usage: runtime-config-to-secret [PATH_TO_RUNTIME_CONFIG] [SECRET_NAME]

path_to_runtime_config="$1"
secret_name="$2"

# Check if PATH_TO_RUNTIME_CONFIG is provided
if [[ -z $path_to_runtime_config ]]; then
  echo "usage: runtime-config-to-secret [PATH_TO_RUNTIME_CONFIG] [SECRET_NAME]"
  echo "missing required argument: PATH_TO_RUNTIME_CONFIG"
  exit 1
fi

# Check if SECRET_NAME is provided
if [[ -z $secret_name ]]; then
  echo "usage: runtime-config-to-secret [PATH_TO_RUNTIME_CONFIG] [SECRET_NAME]"
  echo "missing required argument: SECRET_NAME"
  exit 1
fi

# Check if the runtime-config file exists
if [[ ! -f $path_to_runtime_config ]]; then
  echo "File not found: $path_to_runtime_config"
  exit 1
fi

# Base64 encode the content of runtime-config.toml
encoded_content=$(base64 -i "$path_to_runtime_config")

# Create the Kubernetes secret YAML
cat <<EOF >"${secret_name}.yaml"
apiVersion: v1
kind: Secret
metadata:
  name: $secret_name
type: Opaque
data:
  runtime-config.toml: $encoded_content
EOF

echo "Kubernetes secret created at ${secret_name}.yaml"