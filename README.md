<a name="readme-top"></a>

<div align="center">

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

</div>

<br />
<div align="center">
  <a href="https://github.com/AntoninoAdornetto/go-issue-summoner/assets/70185688/e16afca7-003d-41f3-94a8-1229b182ac73">
    <img src="https://github.com/AntoninoAdornetto/go-issue-summoner/assets/70185688/e16afca7-003d-41f3-94a8-1229b182ac73" alt="Logo" width="300" height="300">
  </a>

<h3 align="center">Go Issue Summoner</h3>

  <p align="center">
    Turn your comments into trackable issues that are reported to your favorite source code management system. 
    <br />
    <!-- @TODO Uncomment 'explore docs' section once we have added documentation. -->
    <!-- <a href="https://github.com/AntoninoAdornetto/go-issue-summoner"><strong>Explore the docs »</strong></a> -->
    <br />
    <br />
    <a href="https://github.com/AntoninoAdornetto/go-issue-summoner/issues">Report Bug</a>
    ·
    <a href="https://github.com/AntoninoAdornetto/go-issue-summoner/issues">Request Feature</a>
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

## Development Status 🚧

This repo is under active development. I am in the early stages of building out the core features. As such, some parts of the program may be missing and change significantly.

<!-- ABOUT THE PROJECT -->

## About The Project

[![Product Name Screen Shot][product-screenshot]](https://example.com)

Go Issue Summoner is a tool that will streamline the process of creating issues within your code base. It works by scanning source code files for special annotations (that you define and pass into the program via a flag) and automatically creates issues on a source code management platform of your choosing. This process will ensure that no important task or concern is overlooked.

## Core Features

- `Customizable Annotations`: Define your own set of annotations, that you would use in a single or multi line comment to mark tasks, concerns, or areas of code that require attention.

<!-- @TODO Uncomment language support Note in README when more lanagues are added -->

- `Language Agnostic`: Annotations are scanned and discovered by locating single and multi line comments and then parsing the information surrounding the annotation. This process is language agnostic and uses the current file extension (when walking the directory) to determine the the proper syntax for a single or multi line comment. **Note: Additional language support will be added soon**

- `SCM Adapters`: Support multiple source code management platforms. GitHub, GitLab, BitBucket etc...

- `Minimized Context Switching`: Developers can write a quick note in their source code file about the issue and then run the report command. Those details will be pushed to the source code management platform you selected and will allow the developer to continue on with their original task with minimal context switching.

- `Discover Issues for contributing to open source projects`: Contributing to open source can be a daunting task. Where does one start? What issue should I tackle first? Well, issue-summoner can be used to locate forgotten issues that may have never been reported and were forgotten about. Simply running `issue-summoner scan` on your favorite open source project may return hundreds of `TODO:` annotations that went under the radar. What a great place to start!

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

[![Go Version](https://img.shields.io/github/go-mod/go-version/AntoninoAdornetto/go-issue-summoner)](https://golang.org/)
[![Cobra](https://img.shields.io/badge/cli-cobra-1abc9c.svg)](https://github.com/spf13/cobra)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

To get started, follow these steps:

### Installation

Install using go. (**_Ensure you have [Go](https://golang.org/doc/install) on your system first._**)

```sh
go install github.com/AntoninoAdornetto/issue-summoner@latest
```

<!-- Install using archive file (**_Helpful if you don't want to install go on your system_**) -->
<!---->
<!-- ### Unix -->
<!---->
<!-- Visit releases page and download the latest version and correct architecture for your system -->
<!---->
<!-- ```sh -->
<!-- # replace X with the correct architecture -->
<!-- tar -xzf go-issue-summoner_X.tar.gz -->
<!---->
<!-- # If you want to make the program executable from anywhere, move to your PATH -->
<!-- sudo mv go-issue-summoner /usr/local/bin -->
<!-- ``` -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->

## Usage

### Authorize Command

In order to publish issues to a source code management system, we must first authorize the program to allow this. Authorizing will look different for each provider. As of now, I have added support for GitHub. I will be adding more in the near future.

- `-s`, `--scm` The source code management platform to authorize. (default is GitHub).

#### Authorize for GitHub

The [device-flow](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow) is utilized to create an access token. The only thing you really need to know here is that when you run the command, you will be given a `user code` in the terminal and your default browser will open to https://github.com/login/device You will then be prompted to enter the user code while the program polls the authorization service for an access token. Once the steps are complete, the program will have all scopes it needs to report issues for you. **Note**: this does grant the program access to both public and private repositories.

```sh
issue-summoner authorize -s github
```

### Scan Command

Scans your local git project for comments that are denoted with an annotation. Details about the comment are constructed through lexical analysis. Each programming language uses it's own lexer to gather the comment tokens and parse information about the comment. Scan is a preliminary command that may be used prior to the `report` command. This will give you an idea of the issue annotations that reside in your project.

- `-a`, `--annotation` The annotation the program will search for. (default annotation is @TODO)

- `-p`, `--path` The path to your local git repository (defaults to your current working directory if a path is not provided)

- `-m`, `--mode` The two modes are `pending` and `processed`. Meaning, you can scan for annotations that have not been uploaded to a source code management platform, I.E pending, or you can scan for annotations that have been published, I.E processed. Processed annotations will look differently than pending annotations because when issues are reported, the program will update the comment, write to the file at the location of the comment, and append the issue id that is tied to the comment. This is so the comment can be removed after it's been resolved.

- `-v`, `--verbose` Logs detailed information about each issue annotation that was located during the scan.

#### Scan Usage

```sh
issue-summoner scan
```

The command will walk your git project directory and check each source file. It adheres to the rules of your projects .gitignore file and skips entire directories and files when it finds a match. Yes, you do not need to worry about your node_modules folder being scanned! The comment syntax to use for each file is based on the files extension. Most languages are supported and more are to come! Let's take a look at an example that uses a single line comment for a C file:

```c
#include <stdio.h>

// @TODO implement the main function
int main() {
    printf("Hello world\n");
    return 0;
}
```

Basic usage of the command would result in the following:

![issue-summoner-scan](https://github.com/AntoninoAdornetto/go-issue-summoner/assets/70185688/f9eaef15-ac50-49d1-b8b2-1c0dd72f8393)

We can get a little more information about the annotation by passing the verbose flag `-v` the result would be:

![issue-summoner-scan-verbose](https://github.com/AntoninoAdornetto/issue-summoner/assets/70185688/45373977-8828-4f57-9371-6486c634bb52)

You may have noticed that there is not a description. This is because single line comments are concise. However, we can be more granular by utilizing a multi line comment:

```c
#include <stdio.h>

int main() {
  /*
   * @TODO implement the main function
   * The main function does nothing useful.
   * Remove the print statement and build something that is useful!
   * */
  printf("Hello world\n");
  return 0;
}
```

The new result using a multi line comment:

![issue-summoner-scan-verbose-multi-line](https://github.com/AntoninoAdornetto/issue-summoner/assets/70185688/09313924-2a02-4000-898e-09b2aeca07a1)

### Report Command

Report is similar to the scan command but with added functionality. It allows you to report the discovered comments, that contain an issue annotation, to a source code management platform. Upon selecting an issue to upload to a source code platform, a template is used to format the data.

- `-a`, `--annotation` The annotation the program will search for. (default annotation is @TODO)

- `-p`, `--path` The path to your local git repository (defaults to your current working directory if a path is not provided)

- `-s`, `--scm` The souce code management platform you would like to upload issues to. Such as, github, gitlab, or bitbucket (default "github")

#### Report usage

```sh
issue-summoner report
```

You are then presented with a list of discovered issues that you can select to report.

`j` - navigate down the list

`k` - navigate up the list

`space` - select an item

`y` - confirm and report the selected issues

![Screenshot_05-Jun_01-18-10_15255](https://github.com/AntoninoAdornetto/issue-summoner/assets/70185688/68769010-031f-4b73-84c0-1d2b59072490)

`Feature note for report command:`

- I plan to take the id returned from each issue that is reported and write the id to the file of where it was discovered. This will allow issue-summoner to programatically remove comments once issues have been marked as resolved. This is one of the next features I plan to implement. The result would be something like this:

```c
int main() {
  // @TODO(1999): do something usefull
  return 0;
  // the id: 1999 can be used to check the status and remove the comment once marked as "resolved"
}
```

<!-- _For more examples, please refer to the [Documentation](https://example.com)_ -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->

## Roadmap

- [x] `Comment/Tag-Annotation Scanning Engine`: Develop the core engine that scans source code files for user defined tag annotations, such as @TODO for to-do items. It can recognize the file's extension and appropriately handle language specific syntax for both single and multi-line comments.

- [ ] `Authenticate User to submit issues`: Verify and Authenticate a user to allow the program to submit issues on the users behalf.

  - [x] GitHub Device Flow
  - [ ] GitLab
  - [ ] BitBucket
        <br></br>

- [ ] `SCM Adapter`: Implement a basic adapter for issue reporting functionality.

  - [x] GitHub Adapter
  - [ ] GitLab Adapater
  - [ ] BitBucket Adapater

See the [open issues](https://github.com/AntoninoAdornetto/go-issue-summoner/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/AntoninoAdornetto/go-issue-summoner.svg?style=for-the-badge
[contributors-url]: https://github.com/AntoninoAdornetto/go-issue-summoner/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/AntoninoAdornetto/go-issue-summoner.svg?style=for-the-badge
[forks-url]: https://github.com/AntoninoAdornetto/go-issue-summoner/network/members
[stars-shield]: https://img.shields.io/github/stars/AntoninoAdornetto/go-issue-summoner.svg?style=for-the-badge
[stars-url]: https://github.com/AntoninoAdornetto/go-issue-summoner/stargazers
[issues-shield]: https://img.shields.io/github/issues/AntoninoAdornetto/go-issue-summoner.svg?style=for-the-badge
[issues-url]: https://github.com/AntoninoAdornetto/go-issue-summoner/issues
[license-shield]: https://img.shields.io/github/license/AntoninoAdornetto/go-issue-summoner.svg?style=for-the-badge
[license-url]: https://github.com/AntoninoAdornetto/go-issue-summoner/blob/master/LICENSE.txt
[product-screenshot]: https://github.com/AntoninoAdornetto/go-issue-summoner/assets/70185688/ccf65400-f43d-4b5b-91ac-46694ccf7d08
