#!/bin/bash
set -e

git_repo_root_dir=$(git rev-parse --show-toplevel)

# pull the agent-api.yaml file from Stainless
AGENT_API_URL="https://app.stainless.com/api/spec/documented/miru-agent/openapi.documented.yml"
AGENT_API_FILE="$git_repo_root_dir/api/agent-api.stainless.yaml"
if ! curl -s -f -o "$AGENT_API_FILE" "$AGENT_API_URL"; then
    echo "❌ Failed to download agent API spec from Stainless"
    exit 1
fi

# Verify the file is not empty
if [ ! -s "$AGENT_API_FILE" ]; then
    echo "❌ Downloaded agent API spec file is empty"
    exit 1
fi

# pull the server-api.yaml file from Stainless
SERVER_API_URL="https://app.stainless.com/api/spec/documented/miru-server/openapi.documented.yml"
SERVER_API_FILE="$git_repo_root_dir/api/server-api.stainless.yaml"
if ! curl -s -f -o "$SERVER_API_FILE" "$SERVER_API_URL"; then
    echo "❌ Failed to download server API spec from Stainless"
    exit 1
fi

# Verify the file is not empty
if [ ! -s "$SERVER_API_FILE" ]; then
    echo "❌ Downloaded server API spec file is empty"
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

# remove the webhooks section from the server-api.yaml file
python3 remove_webhooks.py "$SERVER_API_FILE"

# add Unix socket curl examples to the agent-api.yaml file
python3 add_unix_socket_curl.py "$AGENT_API_FILE"