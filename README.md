![issue-summoner-scan](https://github.com/AntoninoAdornetto/go-issue-summoner/assets/70185688/e9073e64-d160-4857-9dae-bf470d2e50f9)

# Go Issue Summoner

## Development Status :construction:

This repo is under active development. I am in the early stages of building out the core features. As such, some parts of the program may change significantly.

## Overview :world_map:

Scan source code files for actionable comments marked by user-defined tag annotations and report information about the comment to a preferred source code management system.

## Phase 1

- `Comment/Tag-Annotation Scanning Engine`: Develop the core engine that scans source code files for user defined tag annotations, such as `@TODO` for to-do items. It can recognize the file's extension and appropriately handle language specific syntax for both single and multi-line comments.

- `SCM Adapter`: Implement a basic adapter for GitHub to demonstrate issue reporting functionality.

## Phase 2

- `Expand SCM Support`: Develop and integrate additional adapters for other popular SCMs. Such as GitLab, BitBucket, etc.

## Phase 3

- `Advanced Tagging`: Expand the scanning engine to account for **reported** tag annotations and clean up source code files by removing the annotation and surrounding comments if the issue has been resolved.
