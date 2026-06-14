<!--
  Copyright (c) 2026 Bradley Larrick. All rights reserved.

  Licensed under the Apache License v2.0
  https://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
 -->

# Copyright Updater

This program updates the copyright header in all supported files in the specified directory(s).
Supported file types inherent to the program include:

```
	.apt        — Apache APT markup files
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
	.rs         — Rust source files
	.sh         — Shell script files
	.toml       — TOML files
	.txt        — Text files
	.vm         — Apache Velocity template files
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
copyright [options] <source> ...
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

Where `<source>` is the path to a file to process or a directory containing files to process.

### Options
The following options are available:

```bash
  -d <directory> — destination directory for updated files (defaults to updating in-place)
  -e <patterns>  — a list of patterns of directories and/or files to exclude
  -h             — print this message and exit
  -p             — only lists the files that will be updated; no changes are made to the files
  -t <template>  — a copyright template file (defaults to ".copyright.txt")
  -v             — sets verbose logging
  -version       — print the version and exit
```

## Examples

This updates all matching files in the current directory and its subdirectories in-place:
```bash
copyright .
```

This updates all matching files in the `src/main` and `src/test` directories and their subdirectories in-place:
```bash
copyright src/main src/test
```

This updates the copyright header in the `Sample.java` and `SampleTest.java` files:
```bash
copyright src/main/Sample.java src/test/SampleTest.java
```

This lists the files in the current directory and its subdirectories that will be updated without making any changes
to the files:
```bash
copyright -p .
```

This updates all matching files in the current directory and its subdirectories and writes the updated files to `C:/temp`:
```bash
copyright -d C:/temp .
```

This updates all matching files in the current directory and its subdirectories in-place using the template file
`C:/MyCopyrightTemplate.txt`:
```bash
copyright -t C:/MyCopyrightTemplate.txt .
```

This updates all matching files in the current directory and its subdirectories in-place, excluding files in the
`test` and `dist` directories:
```bash
copyright -e '**/test/**/*,**/dist/**/*' .
```

## Global Configuration File

The global configuration file is a TOML file that can be used to set default value for the copyright template
and add custom file extensions to the list of supported files. It is located in the user's home directory under the `.copyright`
subdirectory and is named `.config.toml`; e.g., `C:/Users/username/.copyright/config.toml`.

The format of the file is as follows:

```toml
Copyright = [
'Copyright (c) ${year} Company Name. All rights reserved.',
'',
'Licensed under the Apache License v2.0',
'https://www.apache.org/licenses/LICENSE-2.0']

[[Extensions]]
Extension = '.hcl'
Processor = 'hashtag'
Protected = []
```

### Copyright Template Processing

The copyright template is a multi-line string that is used to add or replace the copyright notice in the files. The template
can include a placeholder for the year (`'${year}'`) which is replaced with an actual value when the files are updated.
If the file does not already contain a copyright notice, the template is added to the top of the file and the placeholder
is replaced with the current year. If the file already contains a copyright notice, the existing notice is replaced with the
new template. If the current copyright notice contains a year range (e.g. `2020-2022`), the placeholder is replaced with
the first year in the range and the current year; e.g. `2020-2022` becomes `2020-2026`. If the current copyright notice
contains a single year that is not the current year (e.g. `2022`), the placeholder is replaced with a range starting with
the old year and ending with the current year; e.g. `2022` becomes `2022-2026`.

The program processes the copyright template in the following order:

1. If a template file is specified on the command line, that template is used. If the file does not exist, the program fails.
2. If no template file is specified, the program searches for a default template file in the current directory (`./.copyright.txt`).
3. If no default template file is found, the program uses the template string defined in the global configuration file.

If the previous three steps fail to find a copyright template, the program fails.

### Supported File Processors

The program supports five types of file processors:

 -  `apt`     — Apache APT (apt) style comments ('`~~`')
 -  `bat`     — MSDos Batch comments ('`REM`')
 -  `hashtag` — Hashtag style comments ('`#`')
 -  `java`    — Java style comments ('`/*`' and '`*/`')
 -  `vm`      — Apache Velocity template comments ('`##`')
 -  `xml`     — XML style comments ('`<!--`' and '`-->`')

The file processors also support 'protected' lines at the beginning of the file that must precede the copyright comment block.
For example, `xml` files have protected lines that start with '`<?xml version`' or '`<!DOCTYPE`'. When processing these files,
the copyright template is not applied until after the protected line(s). Predefined protected lines include:

- `bat`     — '`@echo off`'
- `hashtag` — lines starting with '`#!`'
- `xml`     — lines starting with '`<?xml version`' or '`<!DOCTYPE`'

`APT` files are a special case because there may be multiple protected lines, including ones containing regular text, at the
beginning of the file. In this case, the copyright template is not applied until the processor finds the first blank or
comment line.

Files with a `.vm` extension are also a special case. The file name is inspected for "additional" extensions which may determine which file processor to use. For example, file names ending with `.apt.vm` extension are processed using the `apt` file processor. Files with a `.vm` extension alone or with the extension `txt.vm` are processed using the `vm` file processor.

## Updating with Pre-commit Hooks

To use the program with pre-commit hooks, create a global hooks path in your home directory, copy the included
pre-commit hook script to the hooks path, and set the global hooks path in your Git configuration:
```bash
mkdir -p ~/.githooks
cp <path-to-copyright>/pre-commit.sample ~/.githooks/precommit
git config --global core.hooksPath ~/.githooks
```
The pre-commit script assumes the copyright program is installed in `~/go/bin`.

If, during a commit, some of the staged files require a copyright update, the pre-commit hook will automatically run and update the files as needed. If updates are made, the updated files will need to be re-staged and the commit re-run.
