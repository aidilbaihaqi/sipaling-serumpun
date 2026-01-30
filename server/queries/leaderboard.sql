-- Leaderboard: Top performers ranking
-- Placeholders: {{NAMA_CASES}}, {{INSTANSI_CASES}}, {{SCOPE_CASES}}, {{WHERE_CLAUSE}}
WITH 
user_stats AS (
  SELECT
    ia.assignee_id,
    l.name AS bidang,
    COUNT(*) AS total_issues,
    COUNT(*) FILTER (WHERE s."group" = 'completed') AS total_completed,
    COUNT(*) FILTER (WHERE s."group" IN ('started', 'triage')) AS in_progress,
    AVG(
      CASE 
        WHEN s."group" = 'completed' AND i.start_date IS NOT NULL AND i.completed_at IS NOT NULL
        THEN EXTRACT(EPOCH FROM (i.completed_at - i.start_date)) / 86400.0
        ELSE NULL
      END
    ) AS avg_completion_days
  FROM issues i
  JOIN projects p ON p.id = i.project_id
  JOIN workspaces w ON w.id = p.workspace_id
  JOIN states s ON s.id = i.state_id
  JOIN issue_assignees ia ON ia.issue_id = i.id AND ia.deleted_at IS NULL
  JOIN issue_labels il ON il.issue_id = i.id AND il.deleted_at IS NULL
  JOIN labels l ON l.id = il.label_id
  WHERE w.id = '58f6ec9b-f0ae-4e68-8f05-8f1d9ddf9cac'
    AND p.id = 'cfc12151-e169-4caf-bca9-3eb83ed588ee'
    AND i.deleted_at IS NULL
    AND s."group" != 'cancelled'
  GROUP BY ia.assignee_id, l.name
),
final AS (
  SELECT
    CASE
{{NAMA_CASES}}
      ELSE COALESCE(u.display_name, u.first_name || ' ' || u.last_name, '-')
    END AS nama,
    CASE
{{INSTANSI_CASES}}
      ELSE '-'
    END AS instansi,
    CASE
{{SCOPE_CASES}}
      ELSE 'unknown'
    END AS scope,
    us.bidang,
    us.total_completed,
    us.total_issues,
    CASE
      WHEN us.total_issues = 0 THEN 0
      ELSE ROUND((us.total_completed::numeric * 100) / us.total_issues::numeric, 2)
    END AS completion_rate,
    COALESCE(ROUND(us.avg_completion_days::numeric, 1), 0) AS avg_completion_days,
    us.in_progress
  FROM user_stats us
  JOIN users u ON u.id = us.assignee_id
  WHERE us.total_issues > 0
)
SELECT
  ROW_NUMBER() OVER (ORDER BY completion_rate DESC, total_completed DESC, avg_completion_days ASC) AS rank,
  nama,
  instansi,
  scope,
  bidang,
  total_completed,
  total_issues,
  completion_rate,
  avg_completion_days,
  in_progress,
  CASE
    WHEN ROW_NUMBER() OVER (ORDER BY completion_rate DESC, total_completed DESC, avg_completion_days ASC) <= 3 THEN 'Gold'
    WHEN ROW_NUMBER() OVER (ORDER BY completion_rate DESC, total_completed DESC, avg_completion_days ASC) <= 10 THEN 'Silver'
    WHEN ROW_NUMBER() OVER (ORDER BY completion_rate DESC, total_completed DESC, avg_completion_days ASC) <= 20 THEN 'Bronze'
    ELSE 'Participant'
  END AS badge
FROM final{{WHERE_CLAUSE}}
ORDER BY rank;
