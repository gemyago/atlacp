# Contributing to Sonalmod

## Project Setup

Please have the following tools installed: 
* [direnv](https://github.com/direnv/direnv) 
* [nvm](https://github.com/nvm-sh/nvm) - to setup node. We use node to build install packages, you may skip this part.
* For the most part the this is golang, so have [gobrew](https://github.com/kevincobain2000/gobrew#install-or-update) installed.

Some build and deployment scripts use python. Usually you don't need to run thease scripts, only if you are hacking on those scripts themselves. In this case do python setup and `direnv relaod` ater that:
`python -m venv .venv && source .venv/bin/activate && pip install -r requirements.txt`

## AI Frameworks

[OpenSpec](https://github.com/fission-ai/openspec) - good for structured flow. Use `openspec init`. Not committing to the repo for now.

## Typical golang project tasks

Install/Update dependencies (run from go modules): 
```sh
# Install
go mod download
go get -u tool
go install tool

# Update:
go get -u ./... && go mod tidy
```

## Build and Deployment

- To build binaries and docker images locally see [README](./build/README.md)
- Deployment related notes can be found in [README](./deploy/README.md)