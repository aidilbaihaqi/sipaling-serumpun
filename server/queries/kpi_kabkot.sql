-- KPI Kabkot: Dynamic from CSV (All Ketua)
-- Placeholders: {{NAMA_CASES}}, {{BIDANG_CASES}}, {{INSTANSI_CASES}}, {{JABATAN_CASES}}, {{EMAILS}}, {{WHERE_CLAUSE}}
WITH 
directory AS (
  SELECT 
    u.id AS user_id,
    LOWER(u.email) AS email,
    CASE
{{NAMA_CASES}}
      ELSE COALESCE(u.display_name, u.first_name || ' ' || u.last_name, '-')
    END AS nama,
    CASE
{{INSTANSI_CASES}}
      ELSE 'BPS Kabupaten/Kota'
    END AS instansi,
    CASE
{{JABATAN_CASES}}
      ELSE 'Ketua'
    END AS jabatan,
    CASE
{{BIDANG_CASES}}
      ELSE NULL
    END AS bidang
  FROM users u
  WHERE LOWER(u.email) IN ({{EMAILS}})
),

issue_agg AS (
  SELECT
    ia.assignee_id AS user_id,
    l.name AS bidang,
    COUNT(*) FILTER (WHERE s."group" = 'backlog') AS backlog,
    COUNT(*) FILTER (WHERE s."group" = 'unstarted') AS todo,
    COUNT(*) FILTER (WHERE s."group" IN ('started', 'triage')) AS in_progress,
    COUNT(*) FILTER (WHERE s."group" = 'completed') AS done
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
)

SELECT
  d.nama,
  d.email,
  d.bidang,
  d.instansi,
  d.jabatan,
  COALESCE(ia.backlog, 0) AS backlog,
  COALESCE(ia.todo, 0) AS todo,
  COALESCE(ia.in_progress, 0) AS in_progress,
  COALESCE(ia.done, 0) AS done,
  CASE 
    WHEN COALESCE(ia.backlog, 0) + COALESCE(ia.todo, 0) + COALESCE(ia.in_progress, 0) + COALESCE(ia.done, 0) = 0 
    THEN 0
    ELSE ROUND(
      (COALESCE(ia.done, 0)::numeric * 100) 
      / NULLIF(
        (COALESCE(ia.backlog, 0) + COALESCE(ia.todo, 0) + COALESCE(ia.in_progress, 0) + COALESCE(ia.done, 0))::numeric, 
        0
      ),
      2
    )
  END::float8 AS percent
FROM directory d
LEFT JOIN issue_agg ia ON ia.user_id = d.user_id AND ia.bidang = d.bidang
WHERE d.bidang IS NOT NULL{{WHERE_CLAUSE}}
ORDER BY d.instansi, d.jabatan, d.bidang;
