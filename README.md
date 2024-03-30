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
    Transform your todo comments into trackable issues that are reported to a source code management system of your choosing.
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

Go Issue Summoner is a tool that will stremline the process of creating issues within your code base. It works by scanning source code files for special annotations (that you define and pass into the program via a flag) and automatically creates issues on a source code management platform of your choosing. This process will ensure that no important task or concern is overlooked.

## Core Features

- `Customizable Annotations`: Define your own set of annotations, that you would use in a single or multi line comment to mark tasks, concerns, or areas of code that require attention.

<!-- @TODO Uncomment language support Note in README when more lanagues are added -->

- `Language Agnostic`: Annotations are scanned and discovered by locating single and multi line comments and then parsing the information surrounding the annotation. This process is language agnostic and uses the current file extension (when walking the directory) to determine the the proper syntax for a single or multi line comment. **Note: Additional language support will be added soon**

- `SCM Adapters`: Support multiple source code management platforms. GitHub, GitLab, BitBucket etc...

- `Minimized Context Switching`: Developers can write a quick note in their source code file about the issue and then run the report command. Those details will be pushed to the source code management platform you selected and will allow the developer to continue on with their original task with minimal context switching.

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
go install github.com/AntoninoAdornetto/go-issue-summoner@latest
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

Annotations are scanned and located through single and multi line comments. Lets take a look at an example using a single line comment in go:

```go
// @TODO handle error in main function
// we should not ignore the error that may be returned from `importantFunc`
func main(){
  data, _ := importantFunc()
}
```

We add the annotation and the action that should be taken. Next, we run the `scan` command to see information about all of the issues that are using the `@TODO` annotation in our project. We only have one item at this time, but if you had multiple they would all be displayed.

```sh
# default tag annotation is @TODO you dont have to specify it, if you use this tag.
issue-summoner scan --tag @TODO --verbose

Filename:  main.go
Title:  handle error in main function
Description:  we should not ignore the error that may be returned from `importantFunc`
Start Line number:  8
End Line number:  9
Annotation Line number:  8
Multi line comment:  false
Single line comment:  true

Found 1 (@TODO) tag annotations in your project.
```

You can read more about the scan command [here - (does not exist yet)]() but in short, it provides us information about each annotation that is discovered in your source code. If you wanted to report this issue to GitHub, you could then run the following command after [authorizing - (does not exist yet)]() the program to submit issues on your behalf.

```sh
issue-summoner report --tag @TODO
```

Running the above command will provide a multi select option where you can choose which items you would like to report. The default source code management adapter for the `report` command is GitHub. You can read more about the command [here - (does not exist yet)]()

![issue-summoner-report](https://github.com/AntoninoAdornetto/go-issue-summoner/assets/70185688/04b8ad6b-0791-4dd7-840f-201796d75c97)

<!-- _For more examples, please refer to the [Documentation](https://example.com)_ -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->

## Roadmap

- [x] `Comment/Tag-Annotation Scanning Engine`: Develop the core engine that scans source code files for user defined tag annotations, such as @TODO for to-do items. It can recognize the file's extension and appropriately handle language specific syntax for both single and multi-line comments.

- [x] `Authenticate User to submit issues`: Verify and Authenticate a user to allow the program to submit issues on the users behalf.

  - [x] GitHub Device Flow
  - [ ] GitLabl
  - [ ] BitBucket
        <br></br>

- [ ] `SCM Adapter`: Implement a basic adapter for issue reporting functionality.

  - [ ] GitHub Adapter
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
