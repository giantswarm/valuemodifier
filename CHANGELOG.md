# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.3] - 2025-06-03

- Minor code changes
- Dependency upgrades

## [0.5.3] - 2023-10-25

### Changed

- Upgrade deprecated dependencies.

## [0.5.2] - 2022-11-15

### Fixed

- Modify slice items one by one to prevent loosing the slice content.

## [0.5.1] - 2022-11-10

### Fixed

- Fixed an edge case of modifying fields with dots in their name when the field name got written in a still escaped form

## [0.5.0] - 2022-07-21

### Changed

- Migrate from `golang.org/x/crypto` to `github.com/ProtonMail/go-crypto`, see: https://golang.org/issue/44226.
- Update `github.com/hashicorp/vault/api` to `v1.7.2`
- Use non-deprecated API calls to Vault to encrypt / decrypt secrets

## [0.4.0] - 2021-08-24

### Fixed

- Fix setting values in arrays.

## [0.3.1] - 2021-02-04

### Fixed

- Add support for `null` value in nodes, maps and slices.

## [0.3.0] - 2020-10-21

### Added

- Add modifier for Hashicorp Vault transit secrets.

## [0.2.1] - 2020-09-30

### Changed

- Updated golang dependencies.

## [0.2.0] 2020-03-25

### Changed

- migrate from dep to go modules
- use architect-orb

## [0.1.0] 2020-03-25

### Added

- Added CHANGELOG.md

[Unreleased]: https://github.com/giantswarm/valuemodifier/compare/v0.5.3...HEAD
[0.5.3]: https://github.com/giantswarm/valuemodifier/compare/v0.5.3...v0.5.3
[0.5.3]: https://github.com/giantswarm/valuemodifier/compare/v0.5.2...v0.5.3
[0.5.2]: https://github.com/giantswarm/valuemodifier/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/giantswarm/valuemodifier/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/giantswarm/valuemodifier/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/valuemodifier/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/giantswarm/valuemodifier/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/valuemodifier/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/giantswarm/valuemodifier/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/errors/releases/tag/v0.2.0
[0.1.0]: https://github.com/giantswarm/errors/releases/tag/v0.1.0
