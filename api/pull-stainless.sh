#!/bin/bash
set -e

git_repo_root_dir=$(git rev-parse --show-toplevel)

# pull the device-api.yaml file from Stainless
DEVICE_API_URL="https://app.stainless.com/api/spec/documented/miru-device/openapi.documented.yml"
DEVICE_API_FILE="$git_repo_root_dir/api/device-api.stainless.yaml"
if ! curl -s -f -o "$DEVICE_API_FILE" "$DEVICE_API_URL"; then
    echo "❌ Failed to download device API spec from Stainless"
    exit 1
fi

# Verify the file is not empty
if [ ! -s "$DEVICE_API_FILE" ]; then
    echo "❌ Downloaded device API spec file is empty"
    exit 1
fi

# pull the server-api.yaml file from Stainless
PLATFORM_API_URL="https://app.stainless.com/api/spec/documented/miru-platform/openapi.documented.yml"
PLATFORM_API_FILE="$git_repo_root_dir/api/platform-api.stainless.yaml"
if ! curl -s -f -o "$PLATFORM_API_FILE" "$PLATFORM_API_URL"; then
    echo "❌ Failed to download platform API spec from Stainless"
    exit 1
fi

# Verify the file is not empty
if [ ! -s "$PLATFORM_API_FILE" ]; then
    echo "❌ Downloaded platform API spec file is empty"
    exit 1
fi

# use a virtual environment to run the python script
if [ ! -d ".venv" ]; then
    echo "Creating virtual environment..."
    python3 -m venv .venv
fi
. .venv/bin/activate

# ensure python3 and pyyaml are installed
if ! command -v python3 >/dev/null 2>&1; then
    echo "Python 3 is required but not installed"
    exit 1
fi
if ! python3 -c "import pyyaml" 2>/dev/null; then
    echo "Installing pyyaml..."
    if command -v pip3 >/dev/null 2>&1; then
        python3 -m pip install pyyaml types-pyyaml
    else
        echo "pip not found"
        exit 1
    fi
fi

# add Unix socket curl examples to the device-api.yaml file
python3 add_unix_socket_curl.py "$DEVICE_API_FILE"

# insert the scopes into the platform-api.yaml file
python3 inject_scopes.py "$PLATFORM_API_FILE"

# generate event type MDX pages from the device API spec
python3 generate_event_pages.py "$git_repo_root_dir/docs/references/device-api/v0.2.1/api.yaml" events.yaml "$git_repo_root_dir"
