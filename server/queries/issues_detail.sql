-- Issues Detail: Get ketua_bidang nama from CSV (Provinsi + Kabkot)
WITH 
base AS (
  SELECT
    i.id AS issue_id,
    i.name AS issue_title,
    s."group" AS status,
    i.start_date,
    i.target_date,

    ia.assignee_id,
    u.email AS assignee_email,

    -- Get nama from CSV, fallback to users table
    CASE
{{NAMA_CASES}}
      ELSE COALESCE(
        u.display_name,
        NULLIF(TRIM(COALESCE(u.first_name,'') || ' ' || COALESCE(u.last_name,'')), ''),
        '-'
      )
    END AS assignee_name,

    -- Get scope from CSV
    CASE
{{SCOPE_CASES}}
      ELSE 'unknown'
    END AS scope,

    -- ekstraksi kode kab/kota
    SUBSTRING(
      COALESCE(u.display_name,'') || ' ' ||
      COALESCE(u.first_name,'')   || ' ' ||
      COALESCE(u.last_name,''),
      '(21[0-9]{2})'
    ) AS kab_kode,

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
    AND s."group" != 'cancelled'
),
latest_comment AS (
  SELECT DISTINCT ON (ic.issue_id)
    ic.issue_id,
    ic.comment_stripped AS last_comment,
    ic.created_at AS comment_time
  FROM issue_comments ic
  ORDER BY ic.issue_id, ic.created_at DESC
),
final AS (
  SELECT
    b.scope,
    b.issue_title,
    CASE
      WHEN b.assignee_id IS NULL THEN 'Belum Ditugaskan'
      WHEN b.scope = 'provinsi' THEN 'BPS Provinsi Kepulauan Riau'
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
    lc.comment_time
  FROM base b
  LEFT JOIN latest_comment lc ON lc.issue_id = b.issue_id
)
SELECT * FROM final{{WHERE_CLAUSE}}
ORDER BY scope, issue_title;
