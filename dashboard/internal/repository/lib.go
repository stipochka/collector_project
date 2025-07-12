package repository

import (
	"dashboard/internal/models"
	"strings"
)

func buildLimitOffsetQuery(filter models.LogFilter) (string, []any) {
	var limitCondition strings.Builder
	var args []any

	addLimits := func(stmt string, arg any) {
		limitCondition.WriteString(stmt)
		args = append(args, arg)
	}

	if filter.Limit != 0 {
		addLimits(" LIMIT ? ", filter.Limit)
	}

	if filter.Offset != 0 {
		addLimits(" OFFSET ? ", filter.Offset)
	}

	return limitCondition.String(), args
}

func buildTimeRangeCondition(filter models.TimeRangeFilter) (string, []any) {
	whereNeeded := true
	var condition strings.Builder
	var args []any

	addCondition := func(cond string, arg any) {
		if whereNeeded {
			whereNeeded = false
			condition.WriteString(" WHERE ")
		} else {
			condition.WriteString(" AND ")
		}
		condition.WriteString(cond)
		args = append(args, arg)
	}

	if !filter.TimeFrom.IsZero() {
		addCondition("timestamp >= ?", filter.TimeFrom)
	}

	if !filter.TimeTo.IsZero() {
		addCondition("timestamp <= ?", filter.TimeTo)
	}

	return condition.String(), args
}

func buildCondition(filter models.LogFilter) (string, []any) {
	whereNeeded := true
	var condition strings.Builder
	var args []any

	addCondition := func(cond string, arg any) {
		if whereNeeded {
			whereNeeded = false
			condition.WriteString(" WHERE ")
		} else {
			condition.WriteString(" AND ")
		}
		condition.WriteString(cond)
		args = append(args, arg)
	}

	if filter.Level != "" {
		addCondition("level = ?", filter.Level)
	}

	if filter.ServiceName != "" {
		addCondition("service_name = ?", filter.ServiceName)
	}

	if filter.Op != "" {
		addCondition("op = ?", filter.Op)
	}

	if filter.MessageLike != "" {
		addCondition("message ILIKE ?", "%"+filter.MessageLike+"%")
	}

	if !filter.From.IsZero() {
		addCondition("timestamp >= ?", filter.From)
	}

	if !filter.To.IsZero() {
		addCondition("timestamp <= ?", filter.To)
	}

	return condition.String(), args
}
