# ART - the Arch Repository building Tool

This is ART. Please do not throw it away. :)

ART builds a package repository for Arch Linux. It ties together the standard tools (`makepkg`, `repo-add` and
`repo-remove`) into a workflow that is quicker and more robust than when you just glue these tools together with
Makefiles or shell scripts.

## Installation

```bash
$ make && make install
```

## Usage

Invoke as `art`, without any arguments. ART expects a configuration file `./art.toml` in the current working directory,
like this one:

```toml
[[source]]
path = "/path/to/source/directory"

[[source]]
path = "/path/to/another/directory"

[target]
path = "/path/to/output"
name = "my-packages"
```

For each source in this configuration file, the following globs will be expanded to find packages to build:

* `$SOURCE_PATH/*/PKGBUILD` will be built with [makepkg(8)](https://www.archlinux.org/pacman/makepkg.8.html).
* `$SOURCE_PATH/*.pkg.toml` will be built with [holo-build(8)](https://github.com/holocm/holo-build).

Note that directories will not be traversed any deeper than that. For example, any PKGBUILD file must be in a direct
subdirectory of the source path specified in the configuration file.

The packages thus produced will be stored in the `target.path`. The `target.name` defines the file name of the
repository metadata archive. In the example above, it will be at `/path/to/output/my-packages.db.tar.xz`, so Pacman
would find it with the following configuration snippet:

```
[my-packages]
Server = file:///path/to/output
```

ART keeps a cache file (`.art-cache`) in its current working directory to speed up incremental rebuilds.
