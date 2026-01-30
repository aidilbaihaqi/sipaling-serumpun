-- Timeline: Deadline tracking for Gantt Chart
WITH 
base AS (
  SELECT
    i.id AS issue_id,
    i.name AS issue_title,
    s."group" AS status,
    i.start_date,
    i.target_date,
    i.created_at,

    ia.assignee_id,
    u.email AS assignee_email,

    -- Get nama from CSV
    CASE
{{NAMA_CASES}}
      ELSE COALESCE(u.display_name, u.first_name || ' ' || u.last_name, '-')
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
final AS (
  SELECT
    b.scope,
    b.issue_title,
    b.assignee_name AS ketua_bidang,
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
    b.start_date,
    b.target_date,
    CASE
      WHEN b.target_date IS NULL THEN NULL
      ELSE (b.target_date - CURRENT_DATE)
    END AS days_remaining,
    CASE
      WHEN b.target_date IS NULL THEN 'No Deadline'
      WHEN b.status = 'completed' THEN 'Completed'
      WHEN b.target_date < CURRENT_DATE THEN 'Overdue'
      WHEN (b.target_date - CURRENT_DATE) <= 7 THEN 'Warning'
      ELSE 'On Track'
    END AS deadline_status,
    CASE
      WHEN b.status = 'completed' THEN 100
      WHEN b.status IN ('started', 'triage') THEN 50
      WHEN b.status = 'unstarted' THEN 25
      ELSE 0
    END AS progress_percent
  FROM base b
)
SELECT * FROM final{{WHERE_CLAUSE}}
ORDER BY 
  CASE 
    WHEN deadline_status = 'Overdue' THEN 1
    WHEN deadline_status = 'Warning' THEN 2
    WHEN deadline_status = 'On Track' THEN 3
    WHEN deadline_status = 'Completed' THEN 4
    ELSE 5
  END,
  days_remaining NULLS LAST,
  scope,
  kab_kota;
