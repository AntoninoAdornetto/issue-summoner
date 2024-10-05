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
<h3 align="center">Issue Summoner</h3>

  <p align="center">
    Turn those pesky todo comments into track-able issues that can be reported to your favorite source code hosting platform.
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
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#quick-usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

Issue Summoner is a tool designed to streamline the process of managing issues in your codebase. Its primary function is to locate custom defined annotations within your project and use them to improve your issue tracking workflow. Your code can be scanned for these annotations, generate a summary directly in your terminal, and even report the issues to your preferred source code hosting platform.

After reporting an issue, you can use Issue Summoner to check the status of the issues and automatically clean up your code by removing the corresponding comments.

This tool helps keep track of tasks, and concerns to ensure that nothing is overlooked or forgotten in your development process.

### Built With

[![Go Version](https://img.shields.io/github/go-mod/go-version/AntoninoAdornetto/go-issue-summoner)](https://golang.org/)
[![Cobra](https://img.shields.io/badge/cli-cobra-1abc9c.svg)](https://github.com/spf13/cobra)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

### Installation

Install using go. (**_Ensure you have [Go](https://golang.org/doc/install) on your system first._**)

```sh
go install github.com/AntoninoAdornetto/issue-summoner@latest
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->

## Quick Usage

```sh
# authorize for github. 
# Needed if you want to report issues. Not needed if you want to scan a code base for issue annotations
issue-summoner authorize

# scan code base for issues annotated with default "@TODO" annotation
issue-summoner scan

# scan code base for a custom annotation "@FIXME:"
issue-summoner scan -a @FIXME:

# scan code base for issues annotated with "@FIXME:" and print details of each issue
issue-summoner scan -a @FIXME: -v

# scan for issues annotated with "@FIXME:" and allows you to select issues to report to a source code hosting platform
issue-summoner report -a @FIXME: 

# check the status of all issues, annotated with "@FIXME:", that were reported by issue summoner and remove corresponding 
# comments if the status is resolved/closed
issue-summoner scan -m purge -a @FIXME:
```

## Commands Summary

### Authorize Command

In order to publish issues to a source code hosting platform, we must first authorize the program to allow this. Authorizing will look different for each provider. As of now, I have added support for GitHub. More will be added in the near future.

- `-s`, `--sch` The source code hosting platform to authorize. (default is GitHub).

#### Authorize GitHub

The [device-flow](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow) is utilized to create an access token. The only thing you really need to know here is that when you run the command, you will be given a `user code` in the terminal and your default browser will open to https://github.com/login/device You will then be prompted to enter the user code while the program polls the authorization service for an access token. Once the steps are complete, the program will have all scopes it needs to report issues for you. **Note**: this does grant the program access to both public and private repositories.

```sh
issue-summoner authorize -s github
```

### Scan Command

The `scan` command provides functionality for managing and reviewing issues that reside in your codebase. It serves as an aid to the `report` command through two primary modes. `scan`and `purge` mode. These modes help you manage and track issues directly within your codebase using custom annotations.

##### Scan Mode (Default)

In scan mode, the command analyzes your codebase to locate un-reported issues marked by a specific annotation flag (e.g., `@TODO`, `@FIXME`, etc). This mode is useful when you want to:

- Identify issues currently residing in your codebase that have not been reported yet
- Generate a summary of each un-reported issue that includes a description, location (file name, line number) of the issue.

##### Scan Mode usage

```sh
# will return a count of all issues that are annotated with @TODO
issue-summoner scan

# will return a count of all issues that are annotated with @FIXME
issue-summoner scan --annotation @FIXME

 # will return a count and detailed summary of each issue that is annotated with @TODO
issue-summoner scan --verbose

# short flag examples. Will scan for "@TODO:" annotations and print details about each annotation
issue-summoner scan -a @TODO: -v
```

##### Purge Mode

In purge mode, the command analyzes your codebase to locate reported issues marked by a specific annotation flag that has been appended with an issue number (e.g., `@TODO(#405)`). This mode is useful when you want to:

- Identify issues that have been reported, to a source code hosting platform, but haven't been resolved as of yet.
- Generate a summary of each reported issue that includes a description, location (file name, line number) of the issue.
- Checks the status of all reported issues and will remove the comment entirely if it's been marked as resolved on the hosting platform it was reported to.

**note**: issue summoner cannot keep track of issues that were reported outside of using the `report` command.

##### Purge Mode usage

```sh
# checks the status of each reported issue and returns the count of all open issues that are annotated with @TODO
issue-summoner scan --mode purge

# checks the status of each reported issue and return the count of all open issues that are annotated with @FIXME
issue-summoner scan --mode purge --annotation @FIXME

# shorthand flags
# checks the status of each reported issue and returns a detailed summary of each open issue that is annotated with @TODO
issue-summoner scan -m purge -v
```

##### Flags

- `-a`, `--annotation` **string**: The annotation to search for. Example: @TODO, @FIXME, etc. (Default is "@TODO").

- `-d`, `--debug` Log the stack trace when errors occur

- `-h`, `--help` Help for scan

- `-m`, `--mode`, **string**: `scan`: searches for annotations denoted with the --annotation flag. `purge`: searches for annotations that have been appended with an issue number and removes comments if issues are resolved.(Default is "scan")

- `-p`, `--path` **string**: the path to your local git directory. (Defaults to your working directory).

- `-v`, `--verbose` Log the details about each issue annotation that was located during the scan. Can be used with both `scan`, and `purge` modes.

##### Scan usage

```sh
# search for un-reported issues using @TEST_TODO annotation
# and log details about each issue, with verbose flag
issue-summoner scan -a @TEST_TODO -v

# check the status of all reported issues using the annotation @TEST_TODO,
# removes code comments, automatically, for all issues that have a status of resolved and
# print verbose output about all open issues
issue-summoner scan -a @TEST_TODO -m purge -v
```

##### Code Example

issues that have not been reported yet:

```c
int main() {
  // @TODO do something useful in the main function
  return 0;
}

/*
* @TODO do something useful with the sum function
* for the love of god
*/
int sum(int a, int b) {
  return a + b;
}
```

issues that have been reported, using the `report` command:

```c
int main() {
  // @TODO(#504) do something useful in the main function
  return 0;
}

/*
* @TODO(#505) do something useful with the sum function
* for the love of god
*/
int sum(int a, int b) {
  return a + b;
}
```

### Report Command

![Screenshot__937](https://github.com/user-attachments/assets/b5301301-40dd-4208-8002-118e694669c3)

Report is similar to the scan command but with added functionality. It allows you to report selected comments to a source code hosting platform. After all selections are uploaded, the issue id is written to the same location that the comment token is located. Meaning, your todo annotation will be transformed so that issue summoner can be used to remove the entire comment once the issue has been marked as resolved.

- `-a`, `--annotation` The annotation the program will search for. (default annotation is @TODO)

- `-p`, `--path` The path to your local git repository (defaults to your current working directory if a path is not provided)

- `-s`, `--sch` The souce code hosting platform you would like to upload issues to. Such as, github, gitlab, or bitbucket (default "github")

#### Report usage

```sh
issue-summoner report
```

You are then presented with a list of discovered issues that you can select to report.

`j` - navigate down the list

`k` - navigate up the list

`space` - select an item

`y` - confirm and report the selected issues

After the new issue is published, you will notice that your todo annotation is changed to `@ISSUE(issue_id)`, here is an example of how it may look:

#### Before Report command

```c
int main() {
  // @TODO do something usefull
  return 0;
}
```

#### After Report command

```c
int main() {
  // @TODO(#1999): do something usefull
  return 0;
}
```

<!-- _For more examples, please refer to the [Documentation](https://example.com)_ -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->

## Roadmap

- [ ] `Lexical Analysis`: Develop the core engine that scans source code for comment tokens.

  - [x] `C Lexer`: scan, report, and purge comment tokens for c like languages
  - [x] `Shell Lexer`: scan, report, and purge comment tokens for bash/shell
  - [ ] `Python Lexer`: scan, report, and purge comment tokens for python
  - [ ] `Markdown Lexer`: scan, report, and purge comment tokens for markdown
        <br></br>

- [ ] `Authenticate User to submit issues`: Verify and Authenticate a user to allow the program to submit issues on the users behalf.

  - [x] GitHub Device Flow
  - [ ] GitLab
  - [ ] BitBucket
        <br></br>

- [ ] `Source Code Hosting Drivers`: Implement drivers for issue reporting functionality.

  - [x] GitHub Driver
  - [ ] GitLab Driver
  - [ ] BitBucket Driver

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
