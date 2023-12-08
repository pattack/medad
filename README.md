# Medad - Content publisher
Converts Markdown documents to HTML files and uploads them to the server.

# Contents
* [Features](#features)
* [Installation](#installation)
	* [Go](#go)
	* [AUR (Arch User Repository)](#aur-arch-user-repository)
* [Usage](#usage)
* [License](#license)

# Features
* Creates HTML files using overridable templates
* Multilingual
* Filter target articles to be compiled. This reflects on table of contents automatically
* Uploads over FTP

# Installation
## Go
If you're a Gopher and want to install Medad locally in your GOPATH:

```shell
go install -v github.com/pattack/medad/cmd/medad@latest
```

## AUR (Arch User Repository)
If you're using [Arch Linux][archlinux] install [medad][aur-medad] or [medad-git][aur-medad-git] package from [AUR][wiki-aur]

```shell
yay -S medad
```

# Usage
Medad is a cli application consists of subcommands. Run with `--help` switch go get started.
```shell
medad --help
```

# License
This software is [licensed](LICENSE) under the [GPL v3 License][gpl]. © 2023 [Janstun][janstun]

[archlinux]: https://www.archlinux.org/
[aur-medad]: https://aur.archlinux.org/packages/medad
[aur-medad-git]: https://aur.archlinux.org/packages/medad-git
[wiki-aur]: https://wiki.archlinux.org/index.php/AUR
[gpl]: https://www.gnu.org/licenses/gpl-3.0.en.html
[janstun]: http://janstun.com
