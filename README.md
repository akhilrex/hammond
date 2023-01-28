<p align="center">
  <h1 align="center" style="margin-bottom:0">Hammond</h1>
  <p align="center">Current Version - 2022.07.06</p>

  <p align="center">
    A self-hosted vehicle expense tracking system with support for multiple users.
    <br />
    <a href="https://github.com/AlfHou/hammond"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/AlfHou/hammond/issues">Report Bug</a>
    ·
    <a href="https://github.com/AlfHou/hammond/issues">Request Feature</a>
        ·
    <a href="Screenshots.md">Screenshots</a>
  </p>
</p>

<!-- TABLE OF CONTENTS -->

## Table of Contents

- [About the Project](#about-the-project)
  - [Motivation](#motivation)
  - [Built With](#built-with)
  - [Features](#features)
- [Installation](#installation)
- [Contributing](#contributing)
- [License](#license)
- [Roadmap](#roadmap)
- [Contact](#contact)

<!-- ABOUT THE PROJECT -->

## About The Project

Hammond is a self hosted vehicle management system to track fuel and other 
expenses related to all of your vehicles. 
It supports multiple users sharing multiple vehicles. 
It is the logical successor to Clarkson which has not been updated for quite some time now. 
This repo is again a fork of akhilrex's great [project](https://github.com/akhilrex/hammond).

### Motivation and Developer Notes

As mentioned, this project is a fork of 
akhilrex's [project](https://github.com/akhilrex/hammond) which is no longer active. 
To prevent the same from happeing to this project, we are seeking to add more 
maintainers/collaborators who have access to merge PRs.

We are trying our best to update with new features and feedback is very welcome.

The project is written using Go for the backend and Vuejs for the front end. 

![Product Name Screen Shot][product-screenshot] [More Screenshots](Screenshots.md)

### Built With

- [Go](https://golang.org/)
- [Go-Gin](https://github.com/gin-gonic/gin)
- [GORM](https://github.com/go-gorm/gorm)
- [SQLite](https://www.sqlite.org/index.html)
- [VueJS](https://vuejs.org/)
- [Vuex](https://vuex.vuejs.org/)
- [Buefy](https://buefy.org/)

### Features

- Migrate data from Clarkson
- Add/Manage multiple vehicles
- Add/Manage multiple users
- Track fuel and other expenses
- Share vehicles across multiple users
- Save attachment against vehicles
- Quick Entries (take a photo of a receipt or pump screen to make entry later)
- Vehicle level and overall reporting
- Import from Fuelly and Drivvo

## Installation

The easiest way to run Hammond is to run it as a docker container.

### Using Docker

Simple setup without mounted volumes (for testing and evaluation)

```sh
  docker run -d -p 3000:3000 --name=hammond alfhou/hammond
```

Binding local volumes to the container

```sh
   docker run -d -p 3000:3000 --name=hammond -v "/host/path/to/assets:/assets" -v "/host/path/to/config:/config"  alfhou/hammond
```

### Using Docker-Compose

Modify the docker compose file provided [here](https://github.com/alfhou/hammond/blob/master/docker-compose.yml) 
to update the volume and port binding and run the following command

```yaml
version: '2.1'
services:
  hammond:
    image: alfhou/hammond
    container_name: hammond
    volumes:
      - /path/to/config:/config
      - /path/to/data:/assets
    ports:
      - 3000:3000
    restart: unless-stopped
```

```sh
   docker-compose up -d
```

### Build from Source / Ubuntu Installation

Although personally I feel that using the docker container is the best way of using 
and enjoying something like hammond, a lot of people in the community are still not 
comfortable with using Docker and wanted to host it natively on their Linux servers.
Follow the link below to get a guide on how to build hammond from source.

[Build from source / Ubuntu Guide](docs/ubuntu-install.md)

### Environment Variables

| Name | Description                                                                                                                | Default |
| ---- | -------------------------------------------------------------------------------------------------------------------------- | ------- |
| JWT_SECRET | The secret used to sign the JWT token. There is a default value but it is important that you change it to something else| A super strong secret that needs to be changed | 
| PORT | Change the internal port of the application. If you change this you might have to change your docker configuration as well | (empty) |

### Setup

When you open Hammond for the first time after a fresh install, you will be presented with the option to either import data from an existing Clarkson instance or setup a fresh instance.

#### Migration from Clarkson

You will have to ensure that the Clarkson database is accessible from the Hammond deployment. In case it is not directly possible, you can always take a backup of the Clarkson database and host it somewhere accessible to Hammond using a temporary container. If the access problem is sorted, you will have to enter the connection string the Clarkson database in the following format.

```
        user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

You can check the connectivity from the screen as well.

Note: All the users migrated from Clarkson will have their passwords changed to `hammond`

#### Fresh setup

You will have to provide your name, email and password so that an admin user can be created for you.

Once done you will be taken to the login page.

Go through the settings page once and change relevant settings before you start adding vehicles and expenses.

## Contributing

### Dev Setup

If you want to contribute to the project you need to set it up
for development first.

Fork and clone the project. Once you have it on your own machine,
open up a terminal and navigate to the `server/` directory.

In the `server/` directory run the command `go run main.go`.
After some initial
setup, the server should be listening on at port `3000`.

Next, open a new terminal. Navigate to the `ui/` directory and run `npm install`.
This will install all the dependencies for the frontend.
After the command is done running, run `npm run dev`. After some output, the
frontend should be accessible at `http://localhost:8080`.

If you are sent straight to the login screen, try closing the page and opening
it again. You should be greeted with a setup wizard the first time you run the
project.

Now, simply follow the instructions in order to set up your fresh install.

## License

Distributed under the GPL-3.0 License. See `LICENSE` for more information.

## Roadmap

- [ ] More reports
- [ ] Vehicle specific reminders (servicing etc)
- [ ] Native installer for Windows/Linux/MacOS

<!-- CONTACT -->

## Contact

Project Link: [https://github.com/AlfHou/hammond](https://github.com/AlfHou/hammond)

[product-screenshot]: images/screenshot.jpg
