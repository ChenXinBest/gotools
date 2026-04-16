## ADDED Requirements

### Requirement: Process List Display
The system SHALL display a list of all running system processes with their PID, name, command line, CPU usage, memory usage, listen port, and status.

#### Scenario: Display process list
- **WHEN** user navigates to the Process Manager page
- **THEN** system SHALL display all running processes in a table format with columns: PID, Name, CPU%, Memory(MB), Port, Status

#### Scenario: Show CPU usage as percentage
- **WHEN** process list is displayed
- **THEN** CPU usage SHALL be normalized by dividing by CPU core count and displayed as percentage

#### Scenario: Show memory in MB
- **WHEN** process list is displayed
- **THEN** memory usage SHALL be displayed in megabytes (MB)

#### Scenario: Display process status
- **WHEN** process has network connections
- **THEN** system SHALL display the TCP connection status (LISTEN, ESTABLISHED, TIME_WAIT, CLOSE_WAIT)

### Requirement: Process Search
The system SHALL allow users to search processes by name, PID, command line, or port number.

#### Scenario: Search by process name
- **WHEN** user enters a keyword matching a process name
- **THEN** system SHALL highlight and focus on the matching process

#### Scenario: Search by PID
- **WHEN** user enters a numeric PID
- **THEN** system SHALL find the process with that PID

#### Scenario: Search by port
- **WHEN** user enters a port number
- **THEN** system SHALL find the process listening on that port

### Requirement: Process Sorting
The system SHALL allow users to sort the process list by clicking column headers.

#### Scenario: Sort by column
- **WHEN** user clicks on a column header (PID, Name, CPU%, Memory, Port, Status)
- **THEN** system SHALL sort the list by that column in ascending order

#### Scenario: Toggle sort order
- **WHEN** user clicks on the same column header again
- **THEN** system SHALL toggle to descending order

### Requirement: Process Selection
The system SHALL allow users to select multiple processes for batch operations.

#### Scenario: Single selection with checkbox
- **WHEN** user clicks the checkbox next to a process
- **THEN** system SHALL add/remove that process from the selection

#### Scenario: Select all processes
- **WHEN** user clicks the checkbox in the table header
- **THEN** system SHALL select/deselect all visible processes

### Requirement: Process Termination
The system SHALL allow users to terminate selected processes.

#### Scenario: Kill single process
- **WHEN** user clicks the "Kill" button for a specific process
- **THEN** system SHALL terminate that process by PID

#### Scenario: Kill selected processes
- **WHEN** user clicks "Kill Selected" button with multiple processes selected
- **THEN** system SHALL terminate all selected processes sequentially

#### Scenario: Kill confirmation
- **WHEN** user initiates process termination
- **THEN** system SHALL immediately terminate the process without additional confirmation

### Requirement: Auto Refresh
The system SHALL support automatic periodic refresh of the process list.

#### Scenario: Enable auto refresh
- **WHEN** user enables the auto-refresh toggle
- **THEN** system SHALL refresh the process list at the specified interval

#### Scenario: Configure refresh interval
- **WHEN** user sets a refresh interval (in seconds)
- **THEN** system SHALL use that interval for automatic refresh

#### Scenario: Disable auto refresh
- **WHEN** user disables the auto-refresh toggle
- **THEN** system SHALL stop automatic refresh
