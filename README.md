# Ref: https://developer.hashicorp.com/boundary/install


# BT (Boundary Tools)

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/freepik-company/bgos)
![GitHub](https://img.shields.io/github/license/freepik-company/bgos)

![YouTube Channel Subscribers](https://img.shields.io/youtube/channel/subscribers/UCeSb3yfsPNNVr13YsYNvCAw?label=achetronic&link=http%3A%2F%2Fyoutube.com%2Fachetronic)
![GitHub followers](https://img.shields.io/github/followers/achetronic?label=achetronic&link=http%3A%2F%2Fgithub.com%2Fachetronic)
![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/achetronic?style=flat&logo=twitter&link=https%3A%2F%2Ftwitter.com%2Fachetronic)

A super specific process to 

## Motivation

Boundary is...

## Flags

Every configuration parameter can be defined by flags that can be passed to the CLI.
They are described in the following table:

| Name                           | Description                                                           |      Default      | Example                                           |
|:-------------------------------|:----------------------------------------------------------------------|:-----------------:|---------------------------------------------------|
| `--log-level`                  | Define the verbosity of the logs                                      |      `info`       | `--log-level info`                                |
| `--disable-trace`              | Disable traces from logs                                              |      `false`      | `--disable-trace true`                            |
| `--google-sa-credentials-path` | Google ServiceAccount credentials JSON file path                      |   `google.json`   | `--google-sa-credentials-path="~/something.json"` |   
| `--sync-time`                  | Waiting time between group synchronizations (in duration type)        |       `10m`       | `--sync-time 5m`                                  |
| `--google-group`               | (Repeatable or comma-separated list) G.Workspace groups               |        `-`        | `--google-group group1@company.com`               |
| `--boundary-oidc-id`           | Boundary oidc auth method ID to compare its users against G.Workspace | `amoidc_changeme` | `--boundary-oidc-id "amoidc_example"`             |
| `--boundary-scope-id`          | Boundary scope ID where the users and groups are synchronized         |     `global`      | `--boundary-scope-id "global"`                    |

## Environment Variables

Security-critical parameters are managed by environment variables.
They are described in the following table:

| Name                           | Description                                                       | Default | Example                                  |
|:-------------------------------|:------------------------------------------------------------------|:-------:|------------------------------------------|
| `BOUNDARY_ADDR`                | Address where your Boundary instance is hosted                    |   `-`   | `https://hashicorp-boundary.company.com` |
| `BOUNDARY_AUTHMETHODPASS_ID`   | ID of boundary auth method where the privileged user is stored    |   `-`   | `ampw_example`                           |
| `BOUNDARY_AUTHMETHODPASS_USER` | Username of boundary privileged user that perform synchronization |   `-`   | `user_example_changeit`                  |
| `BOUNDARY_AUTHMETHODPASS_PASS` | Password of boundary privileged user that perform synchronization |   `-`   | `super_secure_password`                  |

## Examples

Here you have a complete example to use this command.

> Output is thrown always in JSON as it is more suitable for automations

```console

export BOUNDARY_ADDR="https://your-boundary.you-company.com/"
export BOUNDARY_AUTHMETHODPASS_ID="ampw_example"
export BOUNDARY_AUTHMETHODPASS_USER="automation-google-workspace-groups-syncer" 
export BOUNDARY_AUTHMETHODPASS_PASS='super_secure_password'

bgos run \
     --log-level=info \
     --google-sa-credentials-path=le_credentials.json \
     --google-group sre@your-company.com \
     --google-group developers@your-company.com
```

## How to use

This project provides binary files and Docker images to make it easy to use wherever wanted

### Binaries

Binary files for the most popular platforms will be added to the [releases](https://github.com/freepik-company/bgos/releases)

### Docker

Docker images can be found in GitHub's [packages](https://github.com/freepik-company/bgos/pkgs/container/bgos)
related to this repository

> Do you need it in a different container registry? We think this is not needed, but if we're wrong, please, let's discuss
> it in the best place for that: an issue

## How to contribute

We are open to external collaborations for this project: improvements, bugfixes, whatever.

For doing it, open an issue to discuss the need of the changes, then:

- Fork the repository
- Make your changes to the code
- Open a PR and wait for review

The code will be reviewed and tested (always)

> We are developers and hate bad code. For that reason we ask you the highest quality
> on each line of code to improve this project on each iteration.

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Special mention

This project was done using IDEs from JetBrains. They helped us to develop faster, so we recommend them a lot! 🤓

<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains Logo (Main) logo." width="150">
