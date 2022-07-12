[![Contributors][contributors-shield]][contributors-url] [![Forks][forks-shield]][forks-url] [![Stargazers][stars-shield]][stars-url] [![Issues][issues-shield]][issues-url] [![MIT License][license-shield]][license-url] [![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <!-- <a href="https://github.com/akhilrex/hammond">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a> -->

  <h1 align="center" style="margin-bottom:0">Hammond</h1>
  <p align="center">Current Version - 2022.07.06</p>

  <p align="center">
    A self-hosted vehicle expense tracking system with support for multiple users.
    <br />
    <a href="https://github.com/akhilrex/hammond"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <!-- <a href="https://github.com/akhilrex/hammond">View Demo</a>
    · -->
    <a href="https://github.com/akhilrex/hammond/issues">Report Bug</a>
    ·
    <a href="https://github.com/akhilrex/hammond/issues">Request Feature</a>
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

Hammond is a self hosted vehicle management system to track fuel and other expenses related to all of your vehicles. It supports multiple users sharing multiple vehicles. It is the logical successor to Clarkson which has not been updated for quite some time now.

_Developers Note: This project is under active development which means I release new updates very frequently. It is recommended that you use something like [watchtower](https://github.com/containrrr/watchtower) which will automatically update your containers whenever I release a new version or periodically rebuild the container with the latest image manually._

__Also check out my other self-hosted, open-source solution - [Podgrab](https://github.com/akhilrex/podgrab) - Podcast download and archive manager and player.__
### Motivation and Developer Notes

I was looking for a fuel tracking system and stumbled upon Clarkson. Although it did most of what I needed it has not been updated for quite a lot of time. Since I had some bandwidth available as my previous open source project [Podgrab](http://github.com/akhilrex/podgrab) had become quite stable now, my first thought was to contribute to the Clarkson project only. I soon realized that the architecture that Clarkson had used was not really be that extensible now and would warrant a complete rewrite only. So I decided to build Hammond - The successor to Clarkson.

The current version of Hammond is written using GO for backend and Vuejs for the front end. Originally I had thought of using the same tech stack for both frontend and the backend so that it became easier for users and other developers to use, deploy and contribute. Which is why the first version of Hammond has a NestJS backend complete with all the bells and whistles (GraphQL, Prisma and what nots). But I eventually decided to rebuild the backend in GO just to keep the container size small. No matter how much you can optimize the sheer size of the node_modules will always add bulk to your containers. I host all my tools on my Raspberry Pi. It only makes sense to keep the container size as small as possible.

Also I had initially thought of a 2 container approach (1 for backend and 1 for the frontend) so that they can be independently maintained and updated. I eventually decided against this idea for the sake of simplicity. Although it is safe to assume that most self-hosters are fairly tech capable it still is much better to have a single container that you can fire and forget.

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
- Import from Fuelly (more apps coming soon)

## Installation

The easiest way to run Hammond is to run it as a docker container.

### Using Docker

Simple setup without mounted volumes (for testing and evaluation)

```sh
  docker run -d -p 3000:3000 --name=hammond akhilrex/hammond
```

Binding local volumes to the container

```sh
   docker run -d -p 3000:3000 --name=hammond -v "/host/path/to/assets:/assets" -v "/host/path/to/config:/config"  akhilrex/hammond
```

### Using Docker-Compose

Modify the docker compose file provided [here](https://github.com/akhilrex/hammond/blob/master/docker-compose.yml) to update the volume and port binding and run the following command

```yaml
version: '2.1'
services:
  hammond:
    image: akhilrex/hammond
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

Although personally I feel that using the docker container is the best way of using and enjoying something like hammond, a lot of people in the community are still not comfortable with using Docker and wanted to host it natively on their Linux servers. Follow the link below to get a guide on how to build hammond from source.

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

Akhil Gupta - [@akhilrex](https://twitter.com/akhilrex)

Project Link: [https://github.com/akhilrex/hammond](https://github.com/akhilrex/hammond)

<a href="https://www.buymeacoffee.com/akhilrex" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="width: 217px !important;height: 60px !important;" ></a>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/akhilrex/hammond.svg?style=flat-square
[contributors-url]: https://github.com/akhilrex/hammond/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/akhilrex/hammond.svg?style=flat-square
[forks-url]: https://github.com/akhilrex/hammond/network/members
[stars-shield]: https://img.shields.io/github/stars/akhilrex/hammond.svg?style=flat-square
[stars-url]: https://github.com/akhilrex/hammond/stargazers
[issues-shield]: https://img.shields.io/github/issues/akhilrex/hammond.svg?style=flat-square
[issues-url]: https://github.com/akhilrex/hammond/issues
[license-shield]: https://img.shields.io/github/license/akhilrex/hammond.svg?style=flat-square
[license-url]: https://github.com/akhilrex/hammond/blob/master/LICENSE
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=flat-square&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/akhilrex
[product-screenshot]: images/screenshot.jpg
