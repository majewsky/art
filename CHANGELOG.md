# v2.0.1 (2020-10-06)

Bugfixes:

- Fix building of packages that produce zst-compressed packages. Previously,
  only xz-compressed packages were supported.

# v2.0 (2018-06-21)

Breaking changes:

- This release requires Pacman 5.1 or later.

New features:

- Allow `${pkgname}.PKGBUILD` instead of `${pkgname}/PKGBUILD` for
  directory-less native packages.
- Allow `GPGKEY` to be set in the environment, like makepkg(1) does.

Bugfixes:

- Fix an issue where native packages would fail to build when Pacman/makepkg is
  at version 5.1 or later.  If you mistakenly ran ART 1.x with makepkg 5.1, you
  will need to delete the `.art-cache` in order to build the skipped packages.

# v1.2 (2017-09-11)

Changes:

- The progress display now looks much nicer, and its layout does not break when
  errors are displayed.
- The error that occurred when a target file was newer than its definition has
  been downgraded to a warning because it can appear in legitimate workflows
  where it is not an actual problem.

# v1.1 (2017-09-04)

New features:

- Add an additional table to the `.art-cache` that caches digests of the
  generated target files. This speeds up the average run of ART by 70%.
- Write the `.art-cache` file only if its contents have changed.
- Add validation for the configuration file `art.toml`: ART will now complain
  about missing fields that are required.

Bugfixes:

- Fix a copy-paste error in the `install` target of the Makefile.

# v1.0 (2017-09-01)

Initial release.
