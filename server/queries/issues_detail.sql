WITH base AS (
  SELECT
    i.id AS issue_id,                -- keep UUID for joins
    i.name AS issue_title,
    s.name AS status,
    i.start_date,
    i.target_date,
    i.created_at,
    i.updated_at,

    ia.assignee_id,

    -- nama ketua bidang (aman)
    COALESCE(
      u.display_name,
      NULLIF(TRIM(COALESCE(u.first_name,'') || ' ' || COALESCE(u.last_name,'')), ''),
      '-'
    ) AS assignee_name,

    u.email AS assignee_email,

    -- ekstraksi kode kab/kota
    SUBSTRING(
      COALESCE(u.display_name,'') || ' ' ||
      COALESCE(u.first_name,'')   || ' ' ||
      COALESCE(u.last_name,''),
      '(21[0-9]{2})'
    ) AS kab_kode,

    l.id AS label_id,
    l.name AS bidang
  FROM issues i
  JOIN projects p ON p.id = i.project_id
  JOIN workspaces w ON w.id = p.workspace_id
  JOIN states s ON s.id = i.state_id
  LEFT JOIN issue_assignees ia ON ia.issue_id = i.id AND ia.deleted_at IS NULL
  LEFT JOIN users u ON u.id = ia.assignee_id
  LEFT JOIN issue_labels il ON il.issue_id = i.id AND il.deleted_at IS NULL
  LEFT JOIN labels l ON l.id = il.label_id
  WHERE w.id = '58f6ec9b-f0ae-4e68-8f05-8f1d9ddf9cac'::uuid
    AND p.id = 'cfc12151-e169-4caf-bca9-3eb83ed588ee'::uuid
    AND i.deleted_at IS NULL
),
latest_comment AS (
  SELECT DISTINCT ON (ic.issue_id)
    ic.issue_id,
    ic.comment_stripped AS last_comment,
    ic.created_at AS comment_time,
    ic.created_by_id AS comment_by_id
  FROM issue_comments ic
  ORDER BY ic.issue_id, ic.created_at DESC
),
comment_user AS (
  SELECT
    u.id,
    COALESCE(
      u.display_name,
      NULLIF(TRIM(COALESCE(u.first_name,'') || ' ' || COALESCE(u.last_name,'')), ''),
      '-'
    ) AS commenter_name
  FROM users u
)
SELECT
  b.issue_id::text AS issue_id,      -- cast ONLY for CSV output
  b.issue_title,
  CASE
    WHEN b.assignee_id IS NULL THEN 'Belum Ditugaskan'
    WHEN b.kab_kode IS NULL THEN 'Kode Kab/Kota Tidak Terbaca'
    WHEN b.kab_kode = '2101' THEN 'Karimun'
    WHEN b.kab_kode = '2102' THEN 'Bintan'
    WHEN b.kab_kode = '2103' THEN 'Natuna'
    WHEN b.kab_kode = '2104' THEN 'Lingga'
    WHEN b.kab_kode = '2105' THEN 'Kep. Anambas'
    WHEN b.kab_kode = '2171' THEN 'Batam'
    WHEN b.kab_kode = '2172' THEN 'Tanjung Pinang'
    ELSE 'Lainnya'
  END AS kab_kota,
  COALESCE(b.bidang, '-') AS bidang,
  b.status,
  b.assignee_name AS ketua_bidang,
  COALESCE(b.assignee_email, '-') AS email_ketua_bidang,
  b.start_date,
  b.target_date,
  lc.last_comment,
  lc.comment_time,
  cu.commenter_name AS comment_by,
  b.created_at,
  b.updated_at
FROM base b
LEFT JOIN latest_comment lc ON lc.issue_id = b.issue_id   -- uuid = uuid ✅
LEFT JOIN comment_user cu ON cu.id = lc.comment_by_id
ORDER BY b.updated_at DESC;
