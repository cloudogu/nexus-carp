# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
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
