WITH base AS (
  SELECT
    s.name AS status,
    l.name AS bidang
  FROM issues i
  JOIN projects p ON p.id = i.project_id
  JOIN workspaces w ON w.id = p.workspace_id
  JOIN states s ON s.id = i.state_id
  JOIN issue_labels il ON il.issue_id = i.id AND il.deleted_at IS NULL
  JOIN labels l ON l.id = il.label_id
  WHERE w.id = '58f6ec9b-f0ae-4e68-8f05-8f1d9ddf9cac'
    AND p.id = 'cfc12151-e169-4caf-bca9-3eb83ed588ee'
    AND i.deleted_at IS NULL
)
SELECT
  bidang,
  COUNT(*) AS total,
  COUNT(*) FILTER (WHERE status IN ('Selesai','Completed')) AS selesai,
  COUNT(*) FILTER (
    WHERE status IN ('Dikerjakan','Menunggu','Revisi','On Progress','Perlu Tindak Lanjut')
  ) AS dikerjakan,
  COUNT(*) FILTER (
    WHERE status IN ('Dibatalkan','Cancel','Tidak Dilanjutkan')
  ) AS dibatalkan,
  ROUND(
    COUNT(*) FILTER (WHERE status IN ('Selesai','Completed')) * 100.0 / NULLIF(COUNT(*),0),
    2
  ) AS persen_selesai
FROM base
GROUP BY bidang
ORDER BY persen_selesai ASC, total DESC;
