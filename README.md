# BBB (Boundary But Better)

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/freepik-company/bgos)
![GitHub](https://img.shields.io/github/license/freepik-company/bgos)

![YouTube Channel Subscribers](https://img.shields.io/youtube/channel/subscribers/UCeSb3yfsPNNVr13YsYNvCAw?label=achetronic&link=http%3A%2F%2Fyoutube.com%2Fachetronic)
![GitHub followers](https://img.shields.io/github/followers/achetronic?label=achetronic&link=http%3A%2F%2Fgithub.com%2Fachetronic)
![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/achetronic?style=flat&logo=twitter&link=https%3A%2F%2Ftwitter.com%2Fachetronic)

A super UX friendly CLI to make daily connections through H.Boundary easy to do.

It covers common auth, targets listing, target connections by SSH, Kubernetes, etc 

## Motivation

Original H.Boundary CLI is designed to manage every administration aspect of Boundary (even the hardest ones),
but its usage is not friendly, and some flows are even bugged. This makes original CLI not usable on a daily basis.

This CLI wraps original CLI, fixing things such as UX and bugs on top of Boundary CLI, empowering people to use Boundary
in an easy and reliable way.

## Environment Variables

Only few parameters are managed by environment variables.
They are described in the following table:

| Name                           | Description                                      | Default | Example                                  |
|:-------------------------------|:-------------------------------------------------|:-------:|------------------------------------------|
| `BOUNDARY_ADDR`                | Address where your H.Boundary instance is hosted |   `-`   | `https://hashicorp-boundary.company.com` |


## Quickstart

### 1. Install Hashicorp Boundary in your system

Go to the following URL and install it

https://developer.hashicorp.com/boundary/install

> If you are a super expert, just go here and chose a version: 
> https://releases.hashicorp.com/boundary/

### 2. Install BBB

Simply chose a release and download the binary ready for you:

https://github.com/achetronic/bbb/releases


### 3. Use BBB CLI

```console

export BOUNDARY_ADDR="https://your-boundary.you-company.com/"

bbb auth

```

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
