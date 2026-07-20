# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-07-20

### Changed

- Add Has Processed images to index
- Return if sign has processed when `idonly=true`

### Added

- Prometheus metrics
- Run metrics server on seperate port

## [0.2.0] - 2026-07-17

### Changed

- Consolidated database connection management into re-indexer.

## [0.1.0] - 2026-07-15

### Added

- Initial web application to generate a random road sign from https://roadsign.pictures .  Used as a companion to the [rsp.ui](https://github.com/zmaillard/rsp.ui).
