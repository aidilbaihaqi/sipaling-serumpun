-- Workload: Distribution and balance check
WITH 
user_stats AS (
  SELECT
    ia.assignee_id,
    l.name AS bidang,
    COUNT(*) FILTER (WHERE s."group" != 'completed') AS active_issues,
    COUNT(*) FILTER (WHERE s."group" = 'completed') AS completed_issues,
    COUNT(*) AS total_issues,
    AVG(
      CASE 
        WHEN s."group" = 'completed' AND i.start_date IS NOT NULL AND i.completed_at IS NOT NULL
        THEN EXTRACT(EPOCH FROM (i.completed_at - i.start_date)) / 86400.0
        ELSE NULL
      END
    ) AS avg_days_to_complete
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
avg_workload AS (
  SELECT
    AVG(active_issues) AS avg_active_issues
  FROM user_stats
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
    us.active_issues,
    us.completed_issues,
    us.total_issues,
    ROUND(aw.avg_active_issues::numeric, 1)::float8 AS avg_issues_per_person,
    CASE
      WHEN us.total_issues = 0 THEN 0
      ELSE ROUND((us.completed_issues::numeric * 100) / us.total_issues::numeric, 2)
    END::float8 AS completion_rate,
    COALESCE(ROUND(us.avg_days_to_complete::numeric, 1), 0)::float8 AS avg_days_to_complete,
    CASE
      WHEN us.active_issues > aw.avg_active_issues * 1.5 THEN 'Overloaded'
      WHEN us.active_issues < aw.avg_active_issues * 0.5 THEN 'Underutilized'
      ELSE 'Balanced'
    END AS workload_status
  FROM user_stats us
  CROSS JOIN avg_workload aw
  JOIN users u ON u.id = us.assignee_id
)
SELECT * FROM final{{WHERE_CLAUSE}}
ORDER BY 
  CASE workload_status
    WHEN 'Overloaded' THEN 1
    WHEN 'Balanced' THEN 2
    WHEN 'Underutilized' THEN 3
  END,
  active_issues DESC;
