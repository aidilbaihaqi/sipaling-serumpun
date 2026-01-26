WITH base AS (
  SELECT
    i.id,
    s.name AS status,
    ia.assignee_id,
    SUBSTRING(
      COALESCE(u.display_name,'') || ' ' ||
      COALESCE(u.first_name,'')   || ' ' ||
      COALESCE(u.last_name,''),
      '(21[0-9]{2})'
    ) AS kab_kode
  FROM issues i
  JOIN projects p ON p.id = i.project_id
  JOIN workspaces w ON w.id = p.workspace_id
  JOIN states s ON s.id = i.state_id
  LEFT JOIN issue_assignees ia
    ON ia.issue_id = i.id
   AND ia.deleted_at IS NULL
  LEFT JOIN users u ON u.id = ia.assignee_id
  WHERE w.id = '58f6ec9b-f0ae-4e68-8f05-8f1d9ddf9cac'
    AND p.id = 'cfc12151-e169-4caf-bca9-3eb83ed588ee'
    AND i.deleted_at IS NULL
),
data AS (
  SELECT
    CASE
      WHEN assignee_id IS NULL THEN 'Belum Ditugaskan'
      WHEN kab_kode IS NULL THEN 'Kode Kab/Kota Tidak Terbaca'
      WHEN kab_kode = '2101' THEN 'Karimun'
      WHEN kab_kode = '2102' THEN 'Bintan'
      WHEN kab_kode = '2103' THEN 'Natuna'
      WHEN kab_kode = '2104' THEN 'Lingga'
      WHEN kab_kode = '2105' THEN 'Kep. Anambas'
      WHEN kab_kode = '2171' THEN 'Batam'
      WHEN kab_kode = '2172' THEN 'Tanjung Pinang'
      ELSE 'Lainnya'
    END AS kab_kota,
    status
  FROM base
)
SELECT
  kab_kota,
  COUNT(*) AS total,
  COUNT(*) FILTER (WHERE status IN ('Selesai','Completed')) AS selesai,
  COUNT(*) FILTER (
    WHERE status IN ('Dikerjakan','Menunggu','Revisi','On Progress','Perlu Tindak Lanjut')
  ) AS dikerjakan,
  COUNT(*) FILTER (
    WHERE status IN ('Dibatalkan','Cancel','Tidak Dilanjutkan')
  ) AS dibatalkan,
  (
    ROUND(
      (COUNT(*) FILTER (WHERE status IN ('Selesai','Completed'))::numeric * 100)
      / NULLIF(COUNT(*)::numeric, 0),
      2
    )
  )::float8 AS persen_selesai
FROM data
GROUP BY kab_kota
ORDER BY persen_selesai ASC;
