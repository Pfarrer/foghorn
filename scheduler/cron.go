package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CronField struct {
	min, max int
	values   map[int]bool
}

type CronExpression struct {
	minute     CronField
	hour       CronField
	dayOfMonth CronField
	month      CronField
	dayOfWeek  CronField
}

func ParseCronExpression(expr string) (*CronExpression, error) {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(parts))
	}

	minute, err := parseField(parts[0], 0, 59)
	if err != nil {
		return nil, fmt.Errorf("invalid minute field: %w", err)
	}

	hour, err := parseField(parts[1], 0, 23)
	if err != nil {
		return nil, fmt.Errorf("invalid hour field: %w", err)
	}

	dayOfMonth, err := parseField(parts[2], 1, 31)
	if err != nil {
		return nil, fmt.Errorf("invalid day of month field: %w", err)
	}

	month, err := parseField(parts[3], 1, 12)
	if err != nil {
		return nil, fmt.Errorf("invalid month field: %w", err)
	}

	dayOfWeek, err := parseField(parts[4], 0, 6)
	if err != nil {
		return nil, fmt.Errorf("invalid day of week field: %w", err)
	}

	return &CronExpression{
		minute:     minute,
		hour:       hour,
		dayOfMonth: dayOfMonth,
		month:      month,
		dayOfWeek:  dayOfWeek,
	}, nil
}

func parseField(field string, min, max int) (CronField, error) {
	cf := CronField{
		min:    min,
		max:    max,
		values: make(map[int]bool),
	}

	if field == "*" {
		for i := min; i <= max; i++ {
			cf.values[i] = true
		}
		return cf, nil
	}

	parts := strings.Split(field, ",")
	for _, part := range parts {
		if err := parsePart(part, min, max, cf.values); err != nil {
			return CronField{}, err
		}
	}

	return cf, nil
}

func parsePart(part string, min, max int, values map[int]bool) error {
	if strings.Contains(part, "/") {
		return parseStep(part, min, max, values)
	}

	if strings.Contains(part, "-") {
		return parseRange(part, min, max, values)
	}

	val, err := strconv.Atoi(part)
	if err != nil {
		return fmt.Errorf("invalid value: %s", part)
	}

	if val < min || val > max {
		return fmt.Errorf("value %d out of range [%d, %d]", val, min, max)
	}

	values[val] = true
	return nil
}

func parseRange(part string, min, max int, values map[int]bool) error {
	rangeParts := strings.Split(part, "-")
	if len(rangeParts) != 2 {
		return fmt.Errorf("invalid range: %s", part)
	}

	start, err := strconv.Atoi(rangeParts[0])
	if err != nil {
		return fmt.Errorf("invalid range start: %s", rangeParts[0])
	}

	end, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return fmt.Errorf("invalid range end: %s", rangeParts[1])
	}

	if start < min || start > max {
		return fmt.Errorf("range start %d out of range [%d, %d]", start, min, max)
	}

	if end < min || end > max {
		return fmt.Errorf("range end %d out of range [%d, %d]", end, min, max)
	}

	if start > end {
		return fmt.Errorf("range start %d greater than end %d", start, end)
	}

	for i := start; i <= end; i++ {
		values[i] = true
	}

	return nil
}

func parseStep(part string, min, max int, values map[int]bool) error {
	stepParts := strings.Split(part, "/")
	if len(stepParts) != 2 {
		return fmt.Errorf("invalid step: %s", part)
	}

	step, err := strconv.Atoi(stepParts[1])
	if err != nil {
		return fmt.Errorf("invalid step value: %s", stepParts[1])
	}

	if step <= 0 {
		return fmt.Errorf("step must be positive: %d", step)
	}

	base := "*"
	if stepParts[0] != "*" {
		base = stepParts[0]
	}

	var rangeValues []int
	if base == "*" {
		for i := min; i <= max; i++ {
			rangeValues = append(rangeValues, i)
		}
	} else {
		cf, err := parseField(base, min, max)
		if err != nil {
			return err
		}

		for i := min; i <= max; i++ {
			if cf.values[i] {
				rangeValues = append(rangeValues, i)
			}
		}
	}

	for i, val := range rangeValues {
		if i%step == 0 {
			values[val] = true
		}
	}

	return nil
}

func (c *CronExpression) Next(t time.Time) time.Time {
	next := t.Add(time.Minute).Truncate(time.Minute)

	for {
		if c.matches(next) {
			return next
		}

		next = next.Add(time.Minute)

		if next.Year() > t.Year()+10 {
			return time.Time{}
		}
	}
}

func (c *CronExpression) matches(t time.Time) bool {
	return c.minute.matches(t.Minute()) &&
		c.hour.matches(t.Hour()) &&
		c.dayOfMonth.matches(t.Day()) &&
		c.month.matches(int(t.Month())) &&
		c.dayOfWeek.matches(int(t.Weekday()))
}

func (cf *CronField) matches(value int) bool {
	return cf.values[value]
}

func (c *CronExpression) String() string {
	var parts []string

	parts = append(parts, fieldToString(c.minute, 0, 59))
	parts = append(parts, fieldToString(c.hour, 0, 23))
	parts = append(parts, fieldToString(c.dayOfMonth, 1, 31))
	parts = append(parts, fieldToString(c.month, 1, 12))
	parts = append(parts, fieldToString(c.dayOfWeek, 0, 6))

	return strings.Join(parts, " ")
}

func fieldToString(cf CronField, min, max int) string {
	all := true
	for i := min; i <= max; i++ {
		if !cf.values[i] {
			all = false
			break
		}
	}

	if all {
		return "*"
	}

	values := make([]int, 0, len(cf.values))
	for i := min; i <= max; i++ {
		if cf.values[i] {
			values = append(values, i)
		}
	}

	if len(values) == 1 {
		return strconv.Itoa(values[0])
	}

	ranges := make([]string, 0)
	start := values[0]
	end := values[0]

	for i := 1; i < len(values); i++ {
		if values[i] == end+1 {
			end = values[i]
		} else {
			if start == end {
				ranges = append(ranges, strconv.Itoa(start))
			} else {
				ranges = append(ranges, fmt.Sprintf("%d-%d", start, end))
			}
			start = values[i]
			end = values[i]
		}
	}

	if start == end {
		ranges = append(ranges, strconv.Itoa(start))
	} else {
		ranges = append(ranges, fmt.Sprintf("%d-%d", start, end))
	}

	return strings.Join(ranges, ",")
}
