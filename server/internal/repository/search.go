package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SearchStampsParams struct {
	Query              string
	Name               string
	Tags               []string
	Description        string
	Creator            string
	CreatedSince       *time.Time
	CreatedUntil       *time.Time
	UpdatedSince       *time.Time
	UpdatedUntil       *time.Time
	StampTypeUnicode   string
	StampTypeAnimation string
	CountMonthlyMin    *int
	CountMonthlyMax    *int
	SortBy             string
}

type StampForSearch struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	FileID       uuid.UUID `db:"file_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	CountMonthly int       `db:"count_monthly"`
	Tags         string    `db:"tags"`
	Descriptions string    `db:"descriptions"`
	CreatorName  string    `db:"creator_name"`
}

func (r *Repository) SearchStamps(ctx context.Context, params SearchStampsParams) ([]StampForSearch, error) {
	baseQuery := `
		SELECT
			s.id, s.name, s.file_id, s.created_at, s.updated_at, s.count_monthly,
			COALESCE(GROUP_CONCAT(DISTINCT t.name SEPARATOR ' '), '') AS tags,
			COALESCE(GROUP_CONCAT(DISTINCT sd.description SEPARATOR ' '), '') AS descriptions,
			COALESCE(u.name, '') AS creator_name
		FROM stamps s
	`
	joins := make(map[string]struct{})
	var whereClauses []string
	args := make(map[string]interface{})

	createOrClauses := func(field, query, argPrefix string) string {
		terms := strings.Fields(query)
		if len(terms) == 0 {
			return ""
		}
		var clauses []string
		for i, term := range terms {
			argName := fmt.Sprintf("%s%d", argPrefix, i)
			clauses = append(clauses, fmt.Sprintf("%s LIKE :%s", field, argName))
			args[argName] = "%" + term + "%"
		}
		return "(" + strings.Join(clauses, " OR ") + ")"
	}

	if params.Description != "" || (params.Query != "" && strings.Contains(params.Query, "description")) {
		joins["LEFT JOIN stamp_descriptions sd ON s.id = sd.stamp_id"] = struct{}{}
	}
	if params.Creator != "" || (params.Query != "" && strings.Contains(params.Query, "creator")) {
		joins["LEFT JOIN users u ON s.creator_id = u.id"] = struct{}{}
	}
	if len(params.Tags) > 0 || (params.Query != "" && strings.Contains(params.Query, "tag")) {
		joins["LEFT JOIN stamp_tags st ON s.id = st.stamp_id"] = struct{}{}
		joins["LEFT JOIN tags t ON st.tag_id = t.id"] = struct{}{}
	}

	if params.Name != "" {
		whereClauses = append(whereClauses, createOrClauses("s.name", params.Name, "name"))
	}
	if params.Description != "" {
		whereClauses = append(whereClauses, createOrClauses("sd.description", params.Description, "desc"))
	}
	if params.Creator != "" {
		whereClauses = append(whereClauses, createOrClauses("u.name", params.Creator, "creator"))
	}
	if len(params.Tags) > 0 {
		whereClauses = append(whereClauses, "s.id IN (SELECT st.stamp_id FROM stamp_tags st JOIN tags t ON st.tag_id = t.id WHERE t.name IN (:tags))")
		args["tags"] = params.Tags
	}
	if params.Query != "" {
		qClauses := []string{
			createOrClauses("s.name", params.Query, "qName"),
			createOrClauses("sd.description", params.Query, "qDesc"),
		}
		whereClauses = append(whereClauses, "("+strings.Join(qClauses, " OR ")+")")
	}

	if params.CreatedSince != nil {
		whereClauses = append(whereClauses, "s.created_at >= :created_since")
		args["created_since"] = params.CreatedSince
	}
	if params.CreatedUntil != nil {
		whereClauses = append(whereClauses, "s.created_at <= :created_until")
		args["created_until"] = params.CreatedUntil
	}
	if params.UpdatedSince != nil {
		whereClauses = append(whereClauses, "s.updated_at >= :updated_since")
		args["updated_since"] = params.UpdatedSince
	}
	if params.UpdatedUntil != nil {
		whereClauses = append(whereClauses, "s.updated_at <= :updated_until")
		args["updated_until"] = params.UpdatedUntil
	}
	switch params.StampTypeUnicode {
	case "only_unicode":
		whereClauses = append(whereClauses, "s.is_unicode = TRUE")
	case "only_not_unicode":
		whereClauses = append(whereClauses, "s.is_unicode = FALSE")
	}
	switch params.StampTypeAnimation {
	case "only_animated":
		whereClauses = append(whereClauses, "s.is_animated = TRUE")
	case "only_not_animated":
		whereClauses = append(whereClauses, "s.is_animated = FALSE")
	}
	if params.CountMonthlyMin != nil {
		whereClauses = append(whereClauses, "s.count_monthly >= :count_min")
		args["count_min"] = *params.CountMonthlyMin
	}
	if params.CountMonthlyMax != nil {
		whereClauses = append(whereClauses, "s.count_monthly <= :count_max")
		args["count_max"] = *params.CountMonthlyMax
	}

	finalQuery := baseQuery
	for join := range joins {
		finalQuery += " " + join
	}
	if len(whereClauses) > 0 {
		finalQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	finalQuery += " GROUP BY s.id"

	orderByClause := ""
	switch params.SortBy {
	case "created_at_asc":
		orderByClause = " ORDER BY s.created_at ASC, s.name ASC"
	case "created_at_desc":
		orderByClause = " ORDER BY s.created_at DESC, s.name ASC"
	case "count_monthly_asc":
		orderByClause = " ORDER BY s.count_monthly ASC, s.name ASC"
	case "count_monthly_desc":
		orderByClause = " ORDER BY s.count_monthly DESC, s.name ASC"
	}
	finalQuery += orderByClause

	query, boundArgs, err := sqlx.Named(finalQuery, args)
	if err != nil {
		return nil, fmt.Errorf("failed to bind named query: %w", err)
	}
	query, boundArgs, err = sqlx.In(query, boundArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to expand IN clause: `tags`: %w", err)
	}
	query = r.db.Rebind(query)

	var result []StampForSearch
	if err := r.db.SelectContext(ctx, &result, query, boundArgs...); err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}

	if result == nil {
		return []StampForSearch{}, nil
	}

	return result, nil
}

