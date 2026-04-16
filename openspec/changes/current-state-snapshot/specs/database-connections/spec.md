## ADDED Requirements

### Requirement: Connection Configuration Storage
The system SHALL store database connection configurations in a JSON configuration file.

#### Scenario: Default config location
- **WHEN** application starts
- **THEN** config file SHALL be created at the same directory as the executable (config.json)

#### Scenario: Persist connection list
- **WHEN** user adds, edits, or deletes a database connection
- **THEN** system SHALL save the complete connection list to config file

#### Scenario: Load config on startup
- **WHEN** application starts
- **THEN** system SHALL load existing connections from config file

#### Scenario: Default export settings
- **WHEN** no export settings exist in config
- **THEN** system SHALL use defaults: export tool=mysql-shell, threads=4, compression=gzip, chunk size=64M
