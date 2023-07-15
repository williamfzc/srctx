#!/bin/bash

set -e

# Set version number and base url
version="0.10.0"
base_url="https://github.com/williamfzc/srctx/releases/download"

echo "Starting SRCTX..."

# Determine the filename of the srctx executable based on the OS type and architecture
if [[ "$OSTYPE" == "darwin"* ]]; then
  # Mac OS
  if [[ "$(uname -m)" == "arm64" ]]; then
    # M1 Mac
    filename="srctx-darwin-arm64"
  else
    # Intel Mac
    filename="srctx-darwin-amd64"
  fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
  # Linux
  if [[ "$(uname -m)" == "aarch64" ]]; then
    # ARM 64
    filename="srctx-linux-arm64"
  else
    # x86_64
    filename="srctx-linux-amd64"
  fi
elif [[ "$OSTYPE" == "msys" ]]; then
  # Windows
  filename="srctx-windows-amd64.exe"
else
  echo "Unsupported OS type: $OSTYPE"
  exit 1
fi

# Check if the srctx executable already exists
if [[ ! -f "srctx_bin" ]]; then
  echo "Downloading srctx executable..."
  # Download the srctx executable
  wget "${base_url}/v${version}/${filename}" -O srctx_bin
  chmod +x srctx_bin
fi

# Run srctx on the specified file or directory, with language parameter
if [[ "$SRCTX_LANG" == "GOLANG" ]]; then
  echo "Running srctx for Golang..."
  ./srctx_bin diff --lang GOLANG --withIndex --src "$SRCTX_SRC" --outputHtml ./output.html

elif [[ "$SRCTX_LANG" == "JAVA" ]]; then
  echo "Running srctx for Java..."
  # Download and unpack scip-java.zip
  if [[ ! -f "scip-java.zip" ]]; then
    echo "Downloading scip-java.zip..."
    # we do not always ship this zip
    wget "${base_url}/v0.8.0/scip-java.zip"
  fi
  echo "Extracting scip-java.zip..."
  unzip -o scip-java.zip

  ./scip-java index "$SRCTX_BUILD_CMD"
  ./srctx_bin diff --lang JAVA --src "$SRCTX_SRC" --scip ./index.scip --outputHtml ./output.html

elif [[ "$SRCTX_LANG" == "KOTLIN" ]]; then
  echo "Running srctx for Kotlin..."
  # Download and unpack scip-java.zip
    if [[ ! -f "scip-java.zip" ]]; then
      echo "Downloading scip-java.zip..."
      # we do not always ship this zip
      wget "${base_url}/v0.8.0/scip-java.zip"
    fi
    echo "Extracting scip-java.zip..."
    unzip -o scip-java.zip

  ./scip-java index "$SRCTX_BUILD_CMD"
  ./srctx_bin diff --lang KOTLIN --src "$SRCTX_SRC" --scip ./index.scip --outputHtml ./output.html

else
  echo "Unsupported language: $SRCTX_LANG"
  exit 1
fi

echo "SRCTX finished successfully."
