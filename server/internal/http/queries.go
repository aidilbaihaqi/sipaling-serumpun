package httpx

import (
	"context"
)

type AggCounts struct {
	UserID     string
	Bidang     string
	Backlog    int
	Todo       int
	InProgress int
	Done       int
}

func (s *Server) MapUsersByEmail(ctx context.Context, emails []string) (map[string]string, error) {
	sqlText, err := s.Queries.Load("map_users_by_email.sql")
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(ctx, sqlText, emails)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[string]string, len(emails))
	for rows.Next() {
		var userID, email string
		if err := rows.Scan(&userID, &email); err != nil {
			return nil, err
		}
		m[email] = userID
	}
	return m, rows.Err()
}

func (s *Server) AggIssuesByUserBidang(ctx context.Context, bidangList []string, workspaceID, projectID string) ([]AggCounts, error) {
	sqlText, err := s.Queries.Load("agg_issues_by_user_bidang.sql")
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(ctx, sqlText, bidangList, workspaceID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []AggCounts{}
	for rows.Next() {
		var a AggCounts
		if err := rows.Scan(&a.UserID, &a.Bidang, &a.Backlog, &a.Todo, &a.InProgress, &a.Done); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
