package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
		LEFT JOIN stamp_descriptions sd ON s.id = sd.stamp_id
		LEFT JOIN users u ON s.creator_id = u.id
		LEFT JOIN stamp_tags st ON s.id = st.stamp_id
		LEFT JOIN tags t ON st.tag_id = t.id
	`
	var whereClauses []string
	var havingClauses []string
	var args []interface{}

	if params.CreatedSince != nil {
		whereClauses = append(whereClauses, "s.created_at >= ?")
		args = append(args, params.CreatedSince)
	}
	if params.CreatedUntil != nil {
		whereClauses = append(whereClauses, "s.created_at <= ?")
		args = append(args, params.CreatedUntil)
	}
	if params.UpdatedSince != nil {
		whereClauses = append(whereClauses, "s.updated_at >= ?")
		args = append(args, params.UpdatedSince)
	}
	if params.UpdatedUntil != nil {
		whereClauses = append(whereClauses, "s.updated_at <= ?")
		args = append(args, params.UpdatedUntil)
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
		whereClauses = append(whereClauses, "s.count_monthly >= ?")
		args = append(args, *params.CountMonthlyMin)
	}
	if params.CountMonthlyMax != nil {
		whereClauses = append(whereClauses, "s.count_monthly <= ?")
		args = append(args, *params.CountMonthlyMax)
	}

	addHavingOrClause := func(query string, field string) {
		terms := strings.Fields(query)
		if len(terms) > 0 {
			var clauses []string
			for _, term := range terms {
				clauses = append(clauses, fmt.Sprintf("%s COLLATE utf8mb4_unicode_ci LIKE ?", field))
				args = append(args, "%"+term+"%")
			}
			havingClauses = append(havingClauses, "("+strings.Join(clauses, " OR ")+")")
		}
	}

	if params.Name != "" {
		addHavingOrClause(params.Name, "s.name")
	}
	if params.Description != "" {
		addHavingOrClause(params.Description, "descriptions")
	}
	if params.Creator != "" {
		addHavingOrClause(params.Creator, "creator_name")
	}
	if len(params.Tags) > 0 {
		addHavingOrClause(strings.Join(params.Tags, " "), "tags")
	}
	if params.Query != "" {
		terms := strings.Fields(params.Query)
		if len(terms) > 0 {
			var qClauses []string
			for _, term := range terms {
				qClauses = append(qClauses, "s.name COLLATE utf8mb4_unicode_ci LIKE ? OR descriptions COLLATE utf8mb4_unicode_ci LIKE ? OR tags COLLATE utf8mb4_unicode_ci LIKE ? OR creator_name COLLATE utf8mb4_unicode_ci LIKE ?")
				args = append(args, "%"+term+"%", "%"+term+"%", "%"+term+"%", "%"+term+"%")
			}
			havingClauses = append(havingClauses, "("+strings.Join(qClauses, " OR ")+")")
		}
	}

	orderByClause := ""
	switch params.SortBy {
	case "created_at_asc":
		orderByClause = "ORDER BY s.created_at ASC, s.name ASC"
	case "created_at_desc":
		orderByClause = "ORDER BY s.created_at DESC, s.name ASC"
	case "count_monthly_asc":
		orderByClause = "ORDER BY s.count_monthly ASC, s.name ASC"
	case "count_monthly_desc":
		orderByClause = "ORDER BY s.count_monthly DESC, s.name ASC"
	default: 
		orderByClause = "ORDER BY s.name ASC" 
	}

	finalQuery := baseQuery
	if len(whereClauses) > 0 {
		finalQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	finalQuery += " GROUP BY s.id"
	if len(havingClauses) > 0 {
		finalQuery += " HAVING " + strings.Join(havingClauses, " AND ")
	}
	finalQuery += " " + orderByClause

	query := r.db.Rebind(finalQuery)
	var result []StampForSearch
	if err := r.db.SelectContext(ctx, &result, query, args...); err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w\nQuery: %s\nArgs: %v", err, query, args)
	}

	if result == nil {
		return []StampForSearch{}, nil
	}

	return result, nil
}
