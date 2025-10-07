<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

[//]: # ([![Forks][forks-shield]][forks-url])
[![Contributors][contributors-shield]][contributors-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/gary.norman/forum">
    <img src="assets/icons/codex-logo-green-trimmed.svg" alt="Logo" width="350">
  </a>
  <p align="center">
    A forum for finding all the help you need with code.
    <br />
    <a href="https://github.com/gary-norman/forum/wiki"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://codex.loreworld.live">View Demo</a>
    ·
    <a href="https://github.com/gary-norman/forum/issues/new?labels=type%3A+bug">Report Bug</a>
    ·
    <a href="https://github.com/gary-norman/forum/issues/new?labels=type%3A+feature">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->

[![Codex screenshot][product-screenshot]](https://example.com)

## About The Project

**Codex** is a feature-rich web forum application built with Go that provides a modern, community-driven platform for discussions and knowledge sharing. Designed with developers in mind, Codex offers a comprehensive set of features for creating, organizing, and moderating online communities.

### Key Features

#### User Management
- **Authentication & Authorization** - Secure user registration and login with session-based authentication
- **User Profiles** - Customizable profiles with avatars and user information
- **Role-Based Access** - Moderator and admin roles with granular permissions

#### Content & Communication
- **Posts & Comments** - Create threaded discussions with rich content
- **Channels** - Organize discussions into topic-based channels
- **Channel Membership** - Join/leave channels to customize your feed
- **Reactions** - Like and dislike posts and comments to surface quality content
- **Bookmarks** - Save posts for later reference

#### Moderation & Community Management
- **Content Flags** - Report inappropriate content for moderator review
- **Channel Moderation** - Assign moderators to manage specific channels
- **Channel Rules** - Set and enforce community guidelines per channel
- **Mute Channels** - Hide channels you're not interested in
- **Content Filtering** - Filter posts by channel, reactions, and user-created content

#### Media & Design
- **Image Uploads** - Attach images to posts, channels, and user profiles
- **Responsive Design** - Works seamlessly across desktop and mobile devices

#### Engagement & Rewards
- **Loyalty System** - Track user engagement and participation
- **Search Functionality** - Find posts, channels, and users quickly

### Built With

* **Backend:** Go 1.22+
* **Database:** SQLite3
* **Templating:** Go HTML templates
* **Containerization:** Docker
* **Build System:** Make
<!-- Do a search and replace with your text editor for the following: `gary.norman`, `forum`, `twitter_handle`, `linkedin_username`, `email_client`, `email`, `project_title`, `project_description` -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>




[//]: # (* [![Next][Next.js]][Next-url])

[//]: # (* [![React][React.js]][React-url])

[//]: # (* [![Vue][Vue.js]][Vue-url])

[//]: # (* [![Angular][Angular.io]][Angular-url])

[//]: # (* [![Svelte][Svelte.dev]][Svelte-url])

[//]: # (* [![Laravel][Laravel.com]][Laravel-url])

[//]: # (* [![Bootstrap][Bootstrap.com]][Bootstrap-url])

[//]: # (* [![JQuery][JQuery.com]][JQuery-url])

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

* Go 1.22 or higher
* SQLite3

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/gary-norman/forum.git
   cd forum
   ```

### Usage

#### Option 1: Using Make (Recommended)

The project includes an interactive menu system for easy command execution:

```sh
make menu
```

This will display an interactive menu with the following options:
- **build** - Build the web server application
- **run** - Run the web server application
- **build-run** - Build and run the application in sequence
- **Docker** - Docker management submenu (configure, build image, run container)
- **Scripts** - Script management submenu (install, verify, backup scripts)

You can navigate using arrow keys or type the option number.

#### Option 2: Using Make Directly

```sh
# Build the application
make build

# Run the application
make run

# Build and run
make build-run

# Docker commands
make configure        # Configure Docker options
make build-image      # Build Docker image
make run-container    # Run Docker container

# Script management
make install-scripts  # Install/update scripts
make verify-scripts   # Verify script checksums
make backup-scripts   # Backup scripts
```

#### Option 3: Direct Terminal Commands

If you prefer not to use Make:

```sh
# Build the application
go build -o bin/codex github.com/gary-norman/forum/cmd/server

# Run the application
./bin/codex

# Build and run
go build -o bin/codex github.com/gary-norman/forum/cmd/server && ./bin/codex
```

The server will start on `http://localhost:8888` by default.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->

<!-- Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources. -->

_For more examples, please refer to the [Documentation](https://example.com)_

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

### Completed ✓
- [x] Core forum functionality (posts, comments, reactions)
- [x] User authentication and authorization
- [x] Channel-based organization
- [x] Image upload support
- [x] Search functionality
- [x] Docker deployment
- [x] Interactive build menu

### In Progress 🚧
- [ ] Content moderation (flags, moderators)
- [ ] Error handling improvements (400/500 pages)
- [ ] Enhanced UI/UX refinements
- [ ] Performance optimizations

### Planned 📋

#### Advanced Features
- [ ] Bookmark system
- [ ] Real-time notifications
- [ ] Advanced filtering options
- [ ] User reputation system
- [ ] Markdown support for posts

#### Security Enhancements
- [ ] Rate limiting
- [ ] Enhanced session management
- [ ] Two-factor authentication
- [ ] Content sanitization improvements

#### Moderation Tools
- [ ] Admin dashboard
- [ ] Automated content filtering
- [ ] User reporting system enhancements
- [ ] Moderator activity logs

#### Authentication
- [ ] OAuth integration (GitHub, Google)
- [ ] Email verification
- [ ] Password recovery

See the [open issues](https://github.com/gary-norman/forum/issues) for a full list of proposed features (and known issues).

View progress by [milestone](https://github.com/gary-norman/forum/milestones).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<a href="https://github.com/gary.norman/forum/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=gary.norman/forum" alt="contrib.rocks image" />
</a>



<!-- LICENSE -->

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->

Gary Norman - [@twitter_handle](https://twitter.com/twitter_handle) - email@email_client.com

Kamil Ornal - [@twitter_handle](https://twitter.com/twitter_handle) - email@email_client.com

Project Link: [https://github.com/gary.norman/forum](https://github.com/gary.norman/forum)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->

* []()
* []()
* []()

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/gary-norman/forum.svg?style=for-the-badge
[contributors-url]: https://github.com/gary-norman/forum/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/gary-norman/forum.svg?style=for-the-badge
[forks-url]: https://github.com/gary-norman/forum/network/members
[stars-shield]: https://img.shields.io/github/stars/gary-norman/forum.svg?style=for-the-badge
[stars-url]: https://github.com/gary-norman/forum/stargazers
[issues-shield]: https://img.shields.io/github/issues/gary-norman/forum.svg?style=for-the-badge
[issues-url]: https://github.com/gary-norman/forum/issues
[license-shield]: https://img.shields.io/github/license/gary-norman/forum.svg?style=for-the-badge
[license-url]: https://github.com/gary-norman/forum/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/gary-norman/
[product-screenshot]: /assets/images/screenshot.png
[Next.js]: https://img.shields.io/badge/next.js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white
[Next-url]: https://nextjs.org/
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[Vue.js]: https://img.shields.io/badge/Vue.js-35495E?style=for-the-badge&logo=vuedotjs&logoColor=4FC08D
[Vue-url]: https://vuejs.org/
[Angular.io]: https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white
[Angular-url]: https://angular.io/
[Svelte.dev]: https://img.shields.io/badge/Svelte-4A4A55?style=for-the-badge&logo=svelte&logoColor=FF3E00
[Svelte-url]: https://svelte.dev/
[Laravel.com]: https://img.shields.io/badge/Laravel-FF2D20?style=for-the-badge&logo=laravel&logoColor=white
[Laravel-url]: https://laravel.com
[Bootstrap.com]: https://img.shields.io/badge/Bootstrap-563D7C?style=for-the-badge&logo=bootstrap&logoColor=white
[Bootstrap-url]: https://getbootstrap.com
[JQuery.com]: https://img.shields.io/badge/jQuery-0769AD?style=for-the-badge&logo=jquery&logoColor=white
[JQuery-url]: https://jquery.com 
