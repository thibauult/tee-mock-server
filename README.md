<a name="readme-top"></a>
<p align="center">
        <a href="https://github.com/thibauult/tee-mock-server/graphs/contributors"><img src="https://img.shields.io/github/contributors/thibauult/tee-mock-server.svg?style=for-the-badge" alt="Contributors"></a>
        <a href="https://github.com/thibauult/tee-mock-server/network/members"><img src="https://img.shields.io/github/forks/thibauult/tee-mock-server.svg?style=for-the-badge" alt="Forks"></a>
        <a href="https://github.com/thibauult/tee-mock-server/stargazers"><img src="https://img.shields.io/github/stars/thibauult/tee-mock-server.svg?style=for-the-badge" alt="Stargazers"></a>
        <a href="https://github.com/thibauult/tee-mock-server/issues"><img src="https://img.shields.io/github/issues/thibauult/tee-mock-server.svg?style=for-the-badge" alt="Issues"></a></p><br/>
<div align="center">

# TEE Mock Server

A mock server written in Go that generates signed JWT tokens for simulating Google Cloud Confidential Space authentication.

<h4>
    <a href="#about-the-project">About the Project</a> • <a href="#features">Features</a> • <a href="#-setup">️Setup</a> • <a href="#license">License</a>
</h4>
</div>

## About the Project

The `tee-mock-server` is a Golang-based project designed to simulate a server that generates JWT tokens, 
specifically tailored for Google Cloud's Confidential Space. 
It listens on a Unix domain socket and responds with a newly signed JWT that includes custom claims related 
to [Confidential Space](https://cloud.google.com/confidential-computing/confidential-space/docs/reference/token-claims). 

The server uses an RSA private key to sign the token and handles graceful shutdown with automatic cleanup of 
the socket file on termination. 

This mock server is useful for testing and simulating token-based authentication workflows in Confidential Space environments.

## Features

1. **JWT Token Generation**: The server generates signed JWT tokens using a predefined RSA private key, including custom claims related to Google Cloud Confidential Space.
2. **Unix Domain Socket**: It listens for incoming requests on a Unix domain socket, providing a simple and efficient way to interact with the server, with automatic cleanup of the socket file upon termination.

## ️Setup

### Prerequisites
Before starting the mock server, you must make sure that the `/run/container_launcher` folder exists and you have the 
right to write in it: 
```shell
sudo mkdir /run/container_launcher
sudo chmod -R 777 /run/container_launcher 
```

### Installation
To install this project using Docker, you can simply run the following command:
```shell
docker compose up
```

### Usage
You can easily generate a new token using the following cURL command: 
```shell
sudo curl -s -N --unix-socket /run/container_launcher/teeserver.sock --data '{ "audience": "foobar", "token_type": "PKI"  }' http://localhost/v1/token
```

## Configuration
The TEE Mock Server allows some level of configuration so that the token it generates can vary depending on your needs. 
Here's a list of the different environment variables that can be set when starting the server: 

| Name                              | Default                                       | Description                                                                               |
|-----------------------------------|-----------------------------------------------|-------------------------------------------------------------------------------------------|
| `TEE_GOOGLE_SERVICE_ACCOUNT`      | tee-mock-server@localhost.gserviceaccount.com | The GCP SA that is set in the "google_service_accounts" <br>claims of the generated token |
| `TEE_TOKEN_EXPIRATION_IN_MINUTES` | 5                                             | The token expiration time in minutes                                                      |

## License
[![GitHub License file](https://img.shields.io/github/license/thibauult/tee-mock-server)](https://github.com/thibauult/tee-mock-server/blob/main/LICENSE)

This project is distributed under the [Apache License 2.0](http://www.apache.org/licenses/),
making it open and free for anyone to use and contribute to.
See the [license](./LICENSE) file for detailed terms.

<p align="right"><a href="#readme-top">(Back to top)</a></p>
