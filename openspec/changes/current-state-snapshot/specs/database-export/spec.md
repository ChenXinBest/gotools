## ADDED Requirements

### Requirement: MySQL Shell Export Options
The system SHALL support configurable options for MySQL Shell export.

#### Scenario: Multi-threaded export
- **WHEN** user sets thread count (0-64)
- **THEN** system SHALL use that many parallel threads for export

#### Scenario: Compression selection
- **WHEN** user selects compression type (gzip, gz, zstd, none)
- **THEN** system SHALL apply the selected compression to output

#### Scenario: Chunk size configuration
- **WHEN** user sets chunk size (e.g., 64M, 1G)
- **THEN** system SHALL use that chunk size for dividing large tables

#### Scenario: Skip Definer
- **WHEN** user enables "skip definer" option
- **THEN** exported SQL SHALL NOT include DEFINER clauses

#### Scenario: Skip Binlog
- **WHEN** user enables "skip binlog" option
- **THEN** exported SQL SHALL include SET SQL_LOG_BIN=0

### Requirement: mysqldump Export Options
The system SHALL support configurable options for mysqldump export.

#### Scenario: Single transaction mode
- **WHEN** user enables single-transaction
- **THEN** mysqldump SHALL use --single-transaction flag for consistent snapshots

#### Scenario: Export routines
- **WHEN** user enables routines option
- **THEN** exported SQL SHALL include stored procedures and functions

#### Scenario: Export events
- **WHEN** user enables events option
- **THEN** exported SQL SHALL include event scheduler events
