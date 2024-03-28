![issue-summoner-ascii](https://github.com/AntoninoAdornetto/go-issue-summoner/assets/70185688/ccf65400-f43d-4b5b-91ac-46694ccf7d08)

# Go Issue Summoner: Automated Issue Creation

Issue Summoner is a tool that will streamline the process of creating issues within your code base. It works by scanning source code files for
special annotations (that you define and pass into the program via a flag) and automatically creates issues on a source code managment system of your choosing.
This process will ensure that no important task or concern gets overlooked. It also reduces the amount of context shifting we developers have to endure by enabling
us to write a simple comment in our code, with a description of the todo or action item, and run the report command to create an issue for us to tackle at a later time.

## Development Status :construction:

This repo is under active development. I am in the early stages of building out the core features. As such, some parts of the program may change significantly.

## Features

- `Language Agnostic`: Supports source code files written in many programming languages. **Note: more languages will be added soon**
- `Customizable Annotations`: Define your own set of annotations, that you would use in a single or multi line comment (e.g., `// @TODO`, `/** @IMPORTANT /*`), to mark tasks, concerns, or areas of code that require attention.
- `Minimized Context Switching`: Developers can write a quick note in the code, along with the annotation, about what the issue is and then run the report command. The issue will be uploaded to a SCM and the developer can continue on with their original task without having to shift their focus. This is great for people with ADHD :)

## Overview :world_map:

Break down of what has been implemented and what is to come in the near future.

### Phase 1

- [x] `Comment/Tag-Annotation Scanning Engine`: Develop the core engine that scans source code files for user defined tag annotations, such as `@TODO` for to-do items. It can recognize the file's extension and appropriately handle language specific syntax for both single and multi-line comments.

- [x] `Authenticate User to submit issues`: Verify and Authenticate a user to allow the program to submit issues on the users behalf. We will start with GitHub and progress to other source code management platforms.

- [ ] `SCM Adapter`: Implement a basic adapter for GitHub to demonstrate issue reporting functionality.

### Phase 2

- [ ] `Expand SCM Support`: Develop and integrate additional adapters for other popular SCMs. Such as GitLab, BitBucket, etc.

### Phase 3

- [ ] `Advanced Tagging`: Expand the scanning engine to account for **reported** tag annotations and clean up source code files by removing the annotation and surrounding comments if the issue has been resolved.

# Authorization

### GitHub

We use the oauth device flow to authorize a user that wishes to use GitHub as their adapter. This allows us to Automate issue creation on both public and private repos. Here is an overview of the Authorization process:

1. App requests device and user verification codes and gets a URL where the user will be prompted to enter the user verification code.

2. App opens a new browser (using your default browser) and prompts the user to enter the verification code at [https://github.com/login/device](https://github.com/login/device)

3. App polls for the user authentication status. Once the user has authorized the device, the app will make an HTTP request to get the new access token.

4. The access token is written to a configuration file that the program will check for before creating an issue
