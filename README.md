<!--
Copyright (c) Bradley Larrick
-->

# Copyright Updater

This program updates the copyright header in all supported files in the specified directory(s).
Supported file types include:

```
    		.bash       — Bash script files
    		.bat        — Batch script files
    		.cs         — C# source files
    		.css        — Cascading style sheets
    		.csv        — Comma-separated values
    		.go         — Go source files
    		.gradle     — Gradle build files
    		.groovy     — Groovy source files
    		.html       — HTML files
    		.java       — Java source files
    		.js         — JavaScript files
    		.mk         — Makefiles
    		.properties — Properties files
    		.py         — Python source files
    		.sh         — Shell script files
    		.toml       — TOML files
    		.txt        — Text files
    		.xaml       — XAML files
    		.xmi        — XMI files
    		.xml        — XML files
    		.xsd        — XSD files
    		.yaml       — YAML files
```

In addition, files without extensions are checked for the Bash script (`#!/bin/bash`) or Shell script (`#!/bin/sh`)
shebang and are processed appropriately; files with a name starting with `Jenkinsfile` are processed the same as Java files;
and files with a name starting with `Makefile` are processed as `.mk` files.

## Build

To build the Copyright Updater, first download the latest version of go from the [official website](https://golang.org/dl/).

Then, clone the repository and run `go build` in the root directory.

```bash
git clone https://github.com/bradleylarrick/copyright-updater.git
cd copyright-updater
go build
```

## Execution

The Copyright Updater can be executed from the command line as follows:

```bash
copyright [options] <source directory> ...
  -d directory
        destination directory (defaults to source)
  -e patterns
        a list of directory patterns to exclude
  -h    print this message and exit
  -p    only list files that will be updated
  -t template
        a copyright template file (default ".copyright.txt")
  -v    set verbose logging
```

Where `<source directory>` is the path to the directory containing the files to process.

### Options
The following options are available:

```bash
  -d <directory> — destination directory for updated files (defaults to updating in-place)
  -e <patterns>  — a list of patterns of directories and/or files to exclude
  -p             — only lists the files that will be updated; no changes are made to the files
  -t <template>  — a copyright template file (defaults to ".copyright.txt")
  -v             — sets verbose logging
```

## Examples

```bash
copyright .
```
This updates all matching files in the current directory and its subdirectories in-place.

```bash
copyright src/main src/test
```
This updates all matching files in the `src/main` and `src/test` directories and their subdirectories in-place.

```bash
copyright -p .
```
This lists the files that will be updated without making any changes to the files.

```bash
copyright -d C:/temp .
```
This updates all matching files in the current directory and its subdirectories and writes the updated files to `C:/temp`.

```bash
copyright -t C:/MyCopyrightTemplate.txt .
```
This updates all matching files in the current directory and its subdirectories in-place using the template file
`C:/MyCopyrightTemplate.txt`.

```bash
copyright -e '**/test/**/*,**/dist/**/*' .
```
This updates all matching files in the current directory and its subdirectories in-place, excluding files in the
`test` and `dist` directories.
