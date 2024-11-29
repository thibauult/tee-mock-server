<a name="readme-top"></a>
<p align="center">
        <a href="https://github.com/thibauult/tee-mock-server/graphs/contributors"><img src="https://img.shields.io/github/contributors/thibauult/tee-mock-server.svg?style=for-the-badge" alt="Contributors"></a>
        <a href="https://github.com/thibauult/tee-mock-server/network/members"><img src="https://img.shields.io/github/forks/thibauult/tee-mock-server.svg?style=for-the-badge" alt="Forks"></a>
        <a href="https://github.com/thibauult/tee-mock-server/stargazers"><img src="https://img.shields.io/github/stars/thibauult/tee-mock-server.svg?style=for-the-badge" alt="Stargazers"></a>
        <a href="https://github.com/thibauult/tee-mock-server/issues"><img src="https://img.shields.io/github/issues/thibauult/tee-mock-server.svg?style=for-the-badge" alt="Issues"></a></p><br/>
<div align="center">

# TEE Mock Server

A mock server in Go that generates signed JWT tokens for simulating Google Cloud Confidential Space authentication

</div>

<div align="center"><h4><a href="#-table-of-contents">️Table of Contents</a> • <a href="#about-the-project">About the Project</a> • <a href="#features">Features</a> • <a href="#-setup">️Setup</a> • <a href="#about-the-author">About the Author</a> • <a href="#license">License</a></h4></div>

## ️Table of Contents
 <details>
<summary>Open Contents</summary>

- [TEE Mock Server](#tee-mock-server)
    - [About the Project](#about-the-project)
    - [Features](#features)
    - [️Setup](#setup)
        - [Installation](#installation)
    - [License](#license)
</details>

## About the Project

The tee-mock-server is a Golang-based project designed to simulate a server that generates JWT tokens, 
specifically tailored for Google Cloud's Confidential Space. 
It listens on a Unix domain socket and responds with a newly signed JWT that includes custom claims, 
such as eat_profile, secboot, and others related to confidential computing. 
The server uses an RSA private key to sign the token and handles graceful shutdown with automatic cleanup of 
the socket file on termination. This mock server is useful for testing and simulating token-based authentication 
workflows in confidential computing environments.

## Features

1. **JWT Token Generation**: The server generates signed JWT tokens using RSA private keys, including custom claims related to Google Cloud Confidential Space, for testing and simulating authentication in confidential computing environments.
2. **Unix Domain Socket**: It listens for incoming requests on a Unix domain socket, providing a simple and efficient way to interact with the server, with automatic cleanup of the socket file upon termination.

## ️Setup

### Installation
To install this project, follow these steps:
```shell
docker build --tag tee-server-mock .
docker run --name "tee-mock-server" --rm -it -v /run/container_launcher:/run/container_launcher tee-server-mock
```

### Usage
You can easily generate a new token using the following cURL command: 
```shell
sudo curl -s -N --unix-socket /run/container_launcher/teeserver.sock http://localhost/v1/token
```

## License
[![GitHub License file](https://img.shields.io/github/license/thibauult/tee-mock-server)](https://github.com/thibauult/tee-mock-server/blob/main/LICENSE)

This project is distributed under the [Apache License 2.0](http://www.apache.org/licenses/),
making it open and free for anyone to use and contribute to.
See the [license](./LICENSE) file for detailed terms.

<p align="right"><a href="#readme-top">(Back to top)</a></p>
