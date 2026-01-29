WITH bidang_set AS (
  SELECT UNNEST($1::text[]) AS bidang
),
issue_base AS (
  SELECT
    i.id AS issue_id,
    ia.assignee_id AS user_id,
    l.name AS bidang,
    s."group" AS state_group
  FROM issues i
  JOIN projects p   ON p.id = i.project_id
  JOIN workspaces w ON w.id = p.workspace_id
  JOIN states s     ON s.id = i.state_id
  LEFT JOIN issue_assignees ia
    ON ia.issue_id = i.id AND ia.deleted_at IS NULL
  JOIN issue_labels il
    ON il.issue_id = i.id AND il.deleted_at IS NULL
  JOIN labels l
    ON l.id = il.label_id
  JOIN bidang_set bs
    ON bs.bidang = l.name
  WHERE w.id = $2
    AND p.id = $3
    AND i.deleted_at IS NULL
)
SELECT
  user_id,
  bidang,
  COUNT(*) FILTER (WHERE state_group = 'backlog') AS backlog,
  COUNT(*) FILTER (WHERE state_group = 'unstarted') AS todo,
  COUNT(*) FILTER (WHERE state_group IN ('started','triage')) AS in_progress,
  COUNT(*) FILTER (WHERE state_group = 'completed') AS done
FROM issue_base
WHERE user_id IS NOT NULL
GROUP BY user_id, bidang;
