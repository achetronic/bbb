#!/bin/bash

# TODO
repo_owner="achetronic"
repo_name="bbb"

binary_name="${repo_name}"

os=$(uname | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

# Convert architecture name to the format used in the releases
case $arch in
    x86_64)
        arch="amd64"
        ;;
    aarch64)
        arch="arm64"
        ;;
    i386)
        arch="386"
        ;;
    *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
esac

# Get the latest release to get the proper download URI depending on the system and the architecture.
# The goal is getting a URL like the following:
# https://github.com/$repo_owner/$repo_name/releases/latest/download/bbb-v0.1.0-linux-amd64.tar.gz
echo "Looking for the proper package for your system"
download_url=$(curl -s https://api.github.com/repos/$repo_owner/$repo_name/releases/latest | \
	grep -oP "https://.+?${repo_name}-v\d+\.\d+\.\d+-${os}-${arch}\.tar\.gz" | \
	head -n 1)

if [ -z "$download_url" ]; then
    echo "No suitable release found for OS: $os, Arch: $arch"
    exit 1
fi

# Download the package into /tmp/bbb-install
# Ask the user for confirmation
printf "Downloading the package from: \n${download_url} \n"

read -p "Do you wanna continue? (y/n): " confirm
if [[ "$confirm" != [yY] ]]; then
    echo "Setup cancelled."
    exit 0
fi

curl --silent -L -o /tmp/${binary_name}.tar.gz "${download_url}"

# Create bbb-install directory and uncompress there
mkdir -p "/tmp/bbb-install"
tar -xzf /tmp/${binary_name}.tar.gz -C "/tmp/bbb-install"
cd "/tmp/bbb-install"

# Assuming the tarball contains a binary with the same name as the repository
echo "Installing the binary on your system"
sudo install -m 0755 $binary_name /usr/local/bin/
