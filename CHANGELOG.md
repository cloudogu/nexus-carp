# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v1.5.0] - 2024-09-18
### Changed
- Relicense to AGPL-3.0-only

## [v1.4.1] - 2024-09-16
### Fixed
- [#17] Random password generator did not use `java.security.SecureRandom`.

## [v1.4.0] - 2024-09-04
### Changed
- [#16] update carp version to v1.2.0
  - This adds the ability to bypass CAS-authentication for certain non-browser-requests. Which prevents request-throttling in CAS for requests that only have dogu-internal authentication

## [v1.3.1] - 2022-03-09
### Added
- log level mapping for enhanced logging #14  

## [v1.3.0]
### Fixed
- Forwarding no longer worked successfully with CAS 6 in combination with OIDC. This is now fixed
  - When a user logs in via OIDC, a separate, unique user ID is transmitted by the OIDC provider. 
    This user ID is now used as username (and at the same time as a unique ID); #12 

## [v1.2.0]
### Changed
- Use port 8081 for nexus health check as the Nexus healthcheck endpoint throws a server error on any failed health check; #10

## [v1.1.0] - 2020-09-04
### Changed
- Upgrade to carp v1.1.0; #8

## [v1.0.0] - 2020-07-01
### Changed
- Changed logger to go-logging instead of glog
- Changed log output format
- Make log-level configurable
- Changed from dep to go modules
### Added
- Added modular Makefiles
