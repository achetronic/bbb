#!/bin/bash



############################
## CONSTANTS AND VARIABLES
############################

# Reset
Color_Off='\033[0m'       # Text Reset

# Regular Colors
Black='\033[0;30m'        # Black
Red='\033[0;31m'          # Red
Green='\033[0;32m'        # Green
Yellow='\033[0;33m'       # Yellow
Blue='\033[0;34m'         # Blue
Purple='\033[0;35m'       # Purple
Cyan='\033[0;36m'         # Cyan
White='\033[0;37m'        # White

# Bold
BBlack='\033[1;30m'       # Black
BRed='\033[1;31m'         # Red
BGreen='\033[1;32m'       # Green
BYellow='\033[1;33m'      # Yellow
BBlue='\033[1;34m'        # Blue
BPurple='\033[1;35m'      # Purple
BCyan='\033[1;36m'        # Cyan
BWhite='\033[1;37m'       # White

# Underline
UBlack='\033[4;30m'       # Black
URed='\033[4;31m'         # Red
UGreen='\033[4;32m'       # Green
UYellow='\033[4;33m'      # Yellow
UBlue='\033[4;34m'        # Blue
UPurple='\033[4;35m'      # Purple
UCyan='\033[4;36m'        # Cyan
UWhite='\033[4;37m'       # White

# Background
On_Black='\033[40m'       # Black
On_Red='\033[41m'         # Red
On_Green='\033[42m'       # Green
On_Yellow='\033[43m'      # Yellow
On_Blue='\033[44m'        # Blue
On_Purple='\033[45m'      # Purple
On_Cyan='\033[46m'        # Cyan
On_White='\033[47m'       # White

# High Intensity
IBlack='\033[0;90m'       # Black
IRed='\033[0;91m'         # Red
IGreen='\033[0;92m'       # Green
IYellow='\033[0;93m'      # Yellow
IBlue='\033[0;94m'        # Blue
IPurple='\033[0;95m'      # Purple
ICyan='\033[0;96m'        # Cyan
IWhite='\033[0;97m'       # White

# Bold High Intensity
BIBlack='\033[1;90m'      # Black
BIRed='\033[1;91m'        # Red
BIGreen='\033[1;92m'      # Green
BIYellow='\033[1;93m'     # Yellow
BIBlue='\033[1;94m'       # Blue
BIPurple='\033[1;95m'     # Purple
BICyan='\033[1;96m'       # Cyan
BIWhite='\033[1;97m'      # White

# High Intensity backgrounds
On_IBlack='\033[0;100m'   # Black
On_IRed='\033[0;101m'     # Red
On_IGreen='\033[0;102m'   # Green
On_IYellow='\033[0;103m'  # Yellow
On_IBlue='\033[0;104m'    # Blue
On_IPurple='\033[0;105m'  # Purple
On_ICyan='\033[0;106m'    # Cyan
On_IWhite='\033[0;107m'   # White

# TODO
repo_owner="achetronic"
repo_name="bbb"

binary_name="${repo_name}"

os=$(uname | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)



##########################
## FUNCTIONS DEFINITION
##########################

# Detect shell configuration file used by the user
detect_shell_config_file() {
  local shell_config_file=""

  case "$SHELL" in
    */zsh)
      shell_config_file="$HOME/.zshrc"
      ;;

    */bash)
      if [[ -f "$HOME/.bashrc" ]]; then
        shell_config_file="$HOME/.bashrc"
      elif [[ -f "$HOME/.bash_profile" ]]; then
        shell_config_file="$HOME/.bash_profile"
      elif [[ -f "$HOME/.profile" ]]; then
        shell_config_file="$HOME/.profile"
      fi
      ;;

    */fish)
      shell_config_file="$HOME/.config/fish/config.fish"
      ;;

    *)
      # Fallback to common shell configuration files if specific shell detection fails
      if [[ -f "$HOME/.zshrc" ]]; then
        shell_config_file="$HOME/.zshrc"
      elif [[ -f "$HOME/.bashrc" ]]; then
        shell_config_file="$HOME/.bashrc"
      elif [[ -f "$HOME/.bash_profile" ]]; then
        shell_config_file="$HOME/.bash_profile"
      elif [[ -f "$HOME/.profile" ]]; then
        shell_config_file="$HOME/.profile"
      elif [[ -f "$HOME/.config/fish/config.fish" ]]; then
        shell_config_file="$HOME/.config/fish/config.fish"
      else
        printf "${Red}No supported shell configuration file found. Setup cancelled.${Color_Off}\n"
        return 1
      fi
      ;;
  esac

  echo "$shell_config_file"
}

# Add or modify environment variable BOUNDARY_ADDR
set_boundary_addr() {
  local new_url=$1
  local shell_config_file

  shell_config_file=$(detect_shell_config_file)
  if [[ $? -ne 0 ]]; then
    return 1
  fi

  # Look for the line containing BOUNDARY_ADDR
  if grep -q "export BOUNDARY_ADDR=" "$shell_config_file"; then
    # Modify the line on existence
    sed -i'' -e "s|^export BOUNDARY_ADDR=.*|export BOUNDARY_ADDR=${new_url}|g" "$shell_config_file"
  else
    # Add the line when missing
    echo "export BOUNDARY_ADDR=${new_url}" >> "$shell_config_file"
  fi

  return 0
}



##########################
## 1. MAIN SCRIPT FLOW
##########################

FLOW_EXIT_CODE=0

# Convert architecture name to the format used in the releases
case $arch in
    x86_64)
        arch="amd64"
        ;;
    aarch64)
        arch="arm64"
        ;;
    arm64)
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



########################
## 1.1 DOWNLOAD PACKAGE
########################

printf "${BPurple}1. DOWNLOAD PACKAGE ${Color_Off}\n\n"

# Get the latest release to get the proper download URI depending on the system and the architecture.
# The goal is getting a URL like the following:
# https://github.com/$repo_owner/$repo_name/releases/latest/download/bbb-v0.1.0-linux-amd64.tar.gz
printf "${White}* Looking for the proper package for your system${Color_Off}"
download_url=$(curl -s https://api.github.com/repos/$repo_owner/$repo_name/releases/latest | \
	grep -oE "https://.+?${repo_name}-v[0-9]+\.[0-9]+\.[0-9]+-${os}-${arch}\.tar\.gz" | \
	head -n 1)

printf "\n"

if [ -z "$download_url" ]; then
    echo "${Red}No suitable release found for ${Cyan} OS: $os, Arch: $arch ${Color_Off}"
    exit 1
fi

# Download the package into /tmp/bbb-install
# Ask the user for confirmation
printf "${White}* Downloading the package from: \n${download_url} ${Color_Off}\n\n"

printf "${Green}Do you want to continue? (y/n): ${Color_Off}"
read confirm
if [[ "$confirm" != [yY] ]]; then
    printf "${Red}Setup cancelled.${Color_Off}"
    exit 0
fi

curl --silent -L -o /tmp/${binary_name}.tar.gz "${download_url}"

# Create bbb-install directory and uncompress there
mkdir -p "/tmp/bbb-install"
tar -xzf /tmp/${binary_name}.tar.gz -C "/tmp/bbb-install"
cd "/tmp/bbb-install"



######################
## 1.2 INSTALL BINARY
######################

printf "\n\n\n"
printf "${BPurple}2. INSTALL BINARY ${Color_Off}\n\n"

# Assuming the tarball contains a binary with the same name as the repository
printf "${White}* Installing the binary on your system${Color_Off}"
sudo install -m 0755 $binary_name /usr/local/bin/



#############################
## 1.3 SET ENVIRONMENT VARS
#############################

printf "\n\n\n"
printf "${BPurple}3. SET ENVIRONMENT VARS ${Color_Off}\n\n"

# TODO
printf "${White}* Type URL of your public H.Boundary instance.
The url will be used to set ${Cyan}BOUNDARY_ADDR${White} environment var that is required to properly function${Color_Off}\n"


printf "${Green}BOUNDARY_ADDR (https://example.com): ${Color_Off}"
read hboundary_url

url_regex="^https?://[a-zA-Z0-9.-]+(:[0-9]+)?(/.*)?$"
if [[ ! "$hboundary_url" =~ $url_regex ]]; then
    printf "${Red}Invalid URL format. Setup cancelled.${Color_Off}\n"
    exit 1
fi

set_boundary_addr "$hboundary_url" || FLOW_EXIT_CODE=$?
if [ $FLOW_EXIT_CODE -ne 0 ]; then
  printf "${Red}Impossible to set BOUNDARY_ADDR environment var in your user profile. Please, set it manually.${Color_Off}\n"
  exit 1
fi

printf "\n\n\n"
printf "${BPurple}4. ALMOST COMPLETE ${Color_Off}\n\n"
printf "${White}Yay, you have 'bbb' installed on your system.
${Cyan}Restart your terminal to complete setup.

${White}Then you can authenticate against H.Boundary as follows:
${BWhite}console ~$ ${Green}bbb auth${Color_Off}"
