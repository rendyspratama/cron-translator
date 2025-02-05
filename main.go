package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// joinInts converts a slice of integers into a comma-separated string.
func joinInts(nums []int) string {
	var strNums []string
	for _, num := range nums {
		strNums = append(strNums, fmt.Sprintf("%d", num))
	}
	return strings.Join(strNums, ",")
}

var daysMap = map[string]int{
	"Sunday":    0,
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
}

// translateDaysOfWeek converts day names to their cron equivalents.
func translateDaysOfWeek(days []string) ([]int, error) {
	var result []int
	for _, day := range days {
		if num, exists := daysMap[day]; exists {
			result = append(result, num)
		} else {
			return nil, fmt.Errorf("invalid day of the week: %s", day)
		}
	}
	return result, nil
}

type ScheduleInput struct {
	Type        string    // Daily, Weekly, Monthly, Hourly, Minutes
	TimeRange   [2]string // [Start, End] (e.g., ["14:00", "18:00"]) for Daily, Weekly, Monthly, and Hourly
	MinuteRange [2]int    // [Start, End] (e.g., [0, 59]) for Daily, Weekly, Monthly, and Hourly
	DaysOfWeek  []string  // For Weekly schedules
	Dates       []int     // For Monthly schedules
	Interval    int       // For Minutes or Hourly intervals
}

func generateCron(input ScheduleInput) (string, error) {
	var cron string
	var err error

	switch input.Type {
	case "Daily":
		cron, err = parseDailyWithRange(input.TimeRange, input.MinuteRange)
	case "Weekly":
		cron, err = parseWeeklyWithRange(input.TimeRange, input.MinuteRange, input.DaysOfWeek)
	case "Monthly":
		cron, err = parseMonthlyWithRange(input.TimeRange, input.MinuteRange, input.Dates)
	case "Hourly":
		cron, err = parseHourlyWithRange(input.TimeRange, input.MinuteRange)
	case "Minutes":
		cron, err = parseMinutes(input.Interval)
	default:
		err = errors.New("invalid schedule type")
	}

	return cron, err
}

func parseDailyWithRange(timeRange [2]string, minuteRange [2]int) (string, error) {
	startHour, _, err := parseTime(timeRange[0])
	if err != nil {
		return "", err
	}
	endHour, _, err := parseTime(timeRange[1])
	if err != nil {
		return "", err
	}

	minuteExpr := "*"
	if minuteRange[0] != 0 || minuteRange[1] != 59 {
		minuteExpr = fmt.Sprintf("%d-%d", minuteRange[0], minuteRange[1])
	}

	return fmt.Sprintf("%s %d-%d * * *", minuteExpr, startHour, endHour), nil
}

func parseWeeklyWithRange(timeRange [2]string, minuteRange [2]int, daysOfWeek []string) (string, error) {
	startHour, _, err := parseTime(timeRange[0])
	if err != nil {
		return "", err
	}
	endHour, _, err := parseTime(timeRange[1])
	if err != nil {
		return "", err
	}

	minuteExpr := "*"
	if minuteRange[0] != 0 || minuteRange[1] != 59 {
		minuteExpr = fmt.Sprintf("%d-%d", minuteRange[0], minuteRange[1])
	}

	dayInts, err := translateDaysOfWeek(daysOfWeek)
	if err != nil {
		return "", err
	}

	daysExpr := joinInts(dayInts)
	return fmt.Sprintf("%s %d-%d * * %s", minuteExpr, startHour, endHour, daysExpr), nil
}

func parseMonthlyWithRange(timeRange [2]string, minuteRange [2]int, dates []int) (string, error) {
	startHour, _, err := parseTime(timeRange[0])
	if err != nil {
		return "", err
	}
	endHour, _, err := parseTime(timeRange[1])
	if err != nil {
		return "", err
	}

	minuteExpr := "*"
	if minuteRange[0] != 0 || minuteRange[1] != 59 {
		minuteExpr = fmt.Sprintf("%d-%d", minuteRange[0], minuteRange[1])
	}

	datesExpr := joinInts(dates)
	return fmt.Sprintf("%s %d-%d %s * *", minuteExpr, startHour, endHour, datesExpr), nil
}

func parseHourlyWithRange(timeRange [2]string, minuteRange [2]int) (string, error) {
	startHour, _, err := parseTime(timeRange[0])
	if err != nil {
		return "", err
	}
	endHour, _, err := parseTime(timeRange[1])
	if err != nil {
		return "", err
	}

	// Handle minute range
	minuteExpr := "*"
	if minuteRange[0] != 0 || minuteRange[1] != 59 {
		minuteExpr = fmt.Sprintf("%d-%d", minuteRange[0], minuteRange[1])
	}

	return fmt.Sprintf("%s %d-%d * * *", minuteExpr, startHour, endHour), nil
}

func parseMinutes(interval int) (string, error) {
	if interval < 1 || interval > 59 {
		return "", errors.New("invalid minute interval: must be between 1 and 59")
	}

	// Cron notation for "Every X minutes"
	return fmt.Sprintf("*/%d * * * *", interval), nil
}

func parseTime(time string) (int, int, error) {
	parts := strings.Split(time, ":")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid time format, must be HH:mm")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return 0, 0, errors.New("invalid hour value")
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return 0, 0, errors.New("invalid minute value")
	}

	return hour, minute, nil
}

// APIResponse represents the structure of the API response
type APIResponse struct {
	CronType    string `json:"cron_type"`
	Cron        string `json:"cron"`
	Translation string `json:"translation"`
	Error       string `json:"error"`
}

func main() {
	// Examples of user inputs
	examples := []ScheduleInput{
		{
			Type:        "Daily",
			TimeRange:   [2]string{"14:00", "18:00"},
			MinuteRange: [2]int{30, 45},
		},
		{
			Type:        "Hourly",
			TimeRange:   [2]string{"14:00", "18:00"},
			MinuteRange: [2]int{0, 15},
		},
		{
			Type:     "Minutes",
			Interval: 10,
		},
		{
			Type:        "Weekly",
			TimeRange:   [2]string{"10:00", "12:00"},
			MinuteRange: [2]int{10, 30},
			DaysOfWeek:  []string{"Monday", "Wednesday"},
		},
		{
			Type:        "Monthly",
			TimeRange:   [2]string{"18:00", "20:00"},
			MinuteRange: [2]int{0, 10},
			Dates:       []int{1, 15, 30},
		},
	}

	for _, input := range examples {
		response := handleGenerateCron(input)
		jsonResponse, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(jsonResponse))
	}
}

// handleGenerateCron wraps the generateCron logic and adds human-readable translation
func handleGenerateCron(input ScheduleInput) APIResponse {
	cron, err := generateCron(input)
	if err != nil {
		return APIResponse{
			CronType:    input.Type,
			Cron:        "",
			Translation: "",
			Error:       err.Error(),
		}
	}
	translation := translateCron(input)
	return APIResponse{
		CronType:    input.Type,
		Cron:        cron,
		Translation: translation,
		Error:       "",
	}
}

// translateCron generates a human-readable explanation based on ScheduleInput
func translateCron(input ScheduleInput) string {
	var translation strings.Builder

	switch input.Type {
	case "Daily":
		translation.WriteString(fmt.Sprintf("Every minute between %s and %s ", input.TimeRange[0], input.TimeRange[1]))
		if input.MinuteRange[0] != 0 || input.MinuteRange[1] != 59 {
			translation.WriteString(fmt.Sprintf("from minute %d to %d ", input.MinuteRange[0], input.MinuteRange[1]))
		}
		translation.WriteString("every day.")
	case "Hourly":
		translation.WriteString(fmt.Sprintf("Every minute between %s and %s ", input.TimeRange[0], input.TimeRange[1]))
		if input.MinuteRange[0] != 0 || input.MinuteRange[1] != 59 {
			translation.WriteString(fmt.Sprintf("from minute %d to %d ", input.MinuteRange[0], input.MinuteRange[1]))
		}
		translation.WriteString("every hour.")
	case "Minutes":
		translation.WriteString(fmt.Sprintf("Every %d minutes.", input.Interval))
	case "Weekly":
		translation.WriteString(fmt.Sprintf("Every minute between %s and %s ", input.TimeRange[0], input.TimeRange[1]))
		if input.MinuteRange[0] != 0 || input.MinuteRange[1] != 59 {
			translation.WriteString(fmt.Sprintf("from minute %d to %d ", input.MinuteRange[0], input.MinuteRange[1]))
		}
		days := strings.Join(input.DaysOfWeek, ", ")
		translation.WriteString(fmt.Sprintf("on %s.", days))
	case "Monthly":
		translation.WriteString(fmt.Sprintf("Every minute between %s and %s ", input.TimeRange[0], input.TimeRange[1]))
		if input.MinuteRange[0] != 0 || input.MinuteRange[1] != 59 {
			translation.WriteString(fmt.Sprintf("from minute %d to %d ", input.MinuteRange[0], input.MinuteRange[1]))
		}
		dates := joinInts(input.Dates)
		translation.WriteString(fmt.Sprintf("on the %s of the month.", dates))
	}

	return translation.String()
}
