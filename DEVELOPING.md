# CasaOS Development
Here we will describe the steps required to setup a development environment with CasaOS.  

- [Setting up development environment](#setting-up-development-environment)  
    - [Pre-requisites](#pre-requisites)  
    - [1. Fork the Repo](#1.-fork-the-repo)  
    - [2. Clone the repo down](#2.-clone-the-repo-down)  
    - [3. Install dependencies](#3.-install-dependencies)  


## Setting up a development environment
In this section we will walk you through the general process of setting up your development environment to get started. 

### Pre-requisites
The following must be installed in order to get started. The details of how to install them is outside the scope of this doc, but generally they should be able to be installed with your systems package manager (apt, yum, brew, choco, etc).
- Go > v1.17.0
- yarn
- node.js

### 1. Fork the Repo
[Fork the repo](https://docs.github.com/en/get-started/quickstart/fork-a-repo) onto your own GitHub account for developing.  

### 2. Clone the repo down
1. Navigate into your go workspace (check with `go env GOPATH`).
2. Navigate to the appropriate path for github. It should look something like this: `<path from GOPATH>/github.com/<GitHub Username>/`. If it doesn't exist create it. 
3. Clone down the repo with the following: `git clone --recurse-submodules --remote-submodules https://github.com/<your GitHub Username>/CasaOS.git`  

### 3. Install dependencies
1. `cd UI`
2. `yarn install`
3. `yarn build`
4. `cd ..`
5. `go get`  
