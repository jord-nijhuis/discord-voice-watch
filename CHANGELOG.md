# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Added the `notifications.notify-self` option. When set to `true`, the bot will also send a
  a notification to the user that just joined (if he is registered). By default, this is set to false.
- The bot will now remove a server and all its registrations from the database when it is removed from the server.

### Changed
- Styled the messages that are sent to users
- The database is now properly closed on shutdown
- Added the `users` and `servers` tables to the database

### Fixed
- Fixed a potential concurrency issue
## [0.1.0] - 2024-12-10
### Added
- Added initial version of the bot

[Unreleased]: https://github.com/jord-nijhuis/discord-voice-watch/compare/0.1.0...HEAD
[0.1.0]: https://github.com/jord-nijhuis/discord-voice-watch/releases/tag/0.1.0
