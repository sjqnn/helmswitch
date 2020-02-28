# helmswitch
cli for switching between different Helm versions.

## How to build this project:
```bash
# Clone this repo locally to your GOPATH
# Build Project
go mod download
go build
# Try it out.
./helmswitch -h
```
## Install Releases
See binary downloads for mac/linux in [Releases](https://github.com/sjqnn/helmswitch/releases).

## Usage 
Will bring up an interactive prompt to select version of Helm.
```
helmswitch
```
Will install the specified version of Helm.
```
helmswitch v3.1.1
```
Will list latest 50 tagged versions of Helm available to install.
```
helmswitch list
```
