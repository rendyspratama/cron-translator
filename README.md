# Schedule Cron Generator

## Overview
This project provides a Go implementation for generating cron expressions based on different scheduling types. It supports various scheduling options such as Daily, Weekly, Monthly, Hourly, and Minutes-based intervals. The generated cron expressions can be used in cron jobs, task schedulers, or automation workflows.

## Features
- Convert structured scheduling input into cron expressions
- Support for:
  - Daily schedules with time and minute range
  - Weekly schedules with specific days of the week
  - Monthly schedules with specific dates
  - Hourly schedules with specific minute range
  - Minute-based intervals (e.g., every 10 minutes)
- Human-readable translations of the cron expressions
- Error handling for invalid inputs

## Usage

### Example Input
The application takes a `ScheduleInput` struct with the following fields:
```go
ScheduleInput {
    Type        string    // Daily, Weekly, Monthly, Hourly, Minutes
    TimeRange   [2]string // [Start, End] (e.g., ["14:00", "18:00"])
    MinuteRange [2]int    // [Start, End] (e.g., [0, 59])
    DaysOfWeek  []string  // For Weekly schedules (e.g., ["Monday", "Wednesday"])
    Dates       []int     // For Monthly schedules (e.g., [1, 15, 30])
    Interval    int       // For Minutes-based intervals
}
```

### Example Execution
The program processes multiple scheduling inputs:
```go
examples := []ScheduleInput{
    {
        Type:        "Daily",
        TimeRange:   [2]string{"14:00", "18:00"},
        MinuteRange: [2]int{30, 45},
    },
    {
        Type:        "Weekly",
        TimeRange:   [2]string{"10:00", "12:00"},
        MinuteRange: [2]int{10, 30},
        DaysOfWeek:  []string{"Monday", "Wednesday"},
    },
    {
        Type:     "Minutes",
        Interval: 10,
    },
}
```

The application generates cron expressions and prints them:
```json
{
  "cron_type": "Daily",
  "cron": "30-45 14-18 * * *",
  "translation": "Every minute between 14:00 and 18:00 from minute 30 to 45 every day.",
  "error": ""
}
```

## How It Works
- `generateCron(input ScheduleInput)`: Converts input schedule details into a cron expression.
- `parseDailyWithRange()`, `parseWeeklyWithRange()`, `parseMonthlyWithRange()`, `parseHourlyWithRange()`, `parseMinutes()`: Handle specific schedule types.
- `translateCron()`: Generates a human-readable explanation of the schedule.
- `handleGenerateCron(input ScheduleInput)`: Wrapper function that returns structured API responses.

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/schedule-cron-generator.git
   cd schedule-cron-generator
   ```
2. Run the application:
   ```sh
   go run main.go
   ```

## Contributing
Feel free to open an issue or submit a pull request if you have improvements or bug fixes.

## License
This project is licensed under the MIT License.

