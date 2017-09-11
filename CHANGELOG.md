# v1.2 (TBD)

New features:

- The progress display now looks much nicer, and its layout does not break when
  errors are displayed.

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
