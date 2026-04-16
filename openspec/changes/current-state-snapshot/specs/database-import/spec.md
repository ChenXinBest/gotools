## ADDED Requirements

### Requirement: Import File Selection
The system SHALL allow users to select directory containing exported SQL files.

#### Scenario: Select import directory
- **WHEN** user clicks "Select Directory" and chooses a folder
- **THEN** system SHALL scan the directory for .sql or compressed SQL files

#### Scenario: Display importable files
- **WHEN** user has selected a directory with valid export files
- **THEN** system SHALL display list of available databases/tables to import

### Requirement: Import Execution
The system SHALL execute SQL import using MySQL Shell with parallel processing.

#### Scenario: Parallel import
- **WHEN** user initiates import with multiple tables
- **THEN** system SHALL import tables in parallel using configured thread count

#### Scenario: Wait timeout configuration
- **WHEN** user sets wait timeout (0-86400 seconds)
- **THEN** system SHALL use that timeout for MySQL operations

### Requirement: Import Conflict Handling
The system SHALL detect and handle conflicts between imported and existing objects.

#### Scenario: Detect conflicts before import
- **WHEN** user initiates import
- **THEN** system SHALL check if any target database/table already exists

#### Scenario: Display conflict list
- **WHEN** conflicts are detected
- **THEN** system SHALL display list of conflicting objects (databases/tables)

#### Scenario: User resolves conflicts
- **WHEN** conflicts exist
- **THEN** user SHALL be able to choose: drop existing objects, skip, or cancel import
