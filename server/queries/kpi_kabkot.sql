-- KPI Kabupaten/Kota: Hanya Ketua (Kepala Kab/Kot)
-- Output: nama, email, bidang, instansi, jabatan, backlog, todo, in_progress, done, percent
-- Query ini akan dipanggil dari Go dengan parameter user_ids yang sudah difilter
-- NOTE: Issues dengan group 'cancelled' TIDAK dihitung dalam total dan persentase

WITH 
-- 1) Agregasi issues per user-bidang berdasarkan states.group
-- EXCLUDE cancelled issues dari semua perhitungan
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
    AND s."group" != 'cancelled'  -- EXCLUDE cancelled issues
  GROUP BY ia.assignee_id, l.name
),

-- 2) Ambil semua bidang yang ada (exclude Sekretariat)
bidang_list AS (
  SELECT DISTINCT l.name AS bidang
  FROM labels l
  JOIN issue_labels il ON il.label_id = l.id AND il.deleted_at IS NULL
  JOIN issues i ON i.id = il.issue_id AND i.deleted_at IS NULL
  JOIN projects p ON p.id = i.project_id
  WHERE p.id = 'cfc12151-e169-4caf-bca9-3eb83ed588ee'
    AND l.name NOT ILIKE '%sekretariat%'
),

-- 3) User info dari database dengan mapping instansi
user_info AS (
  SELECT 
    u.id AS user_id,
    LOWER(u.email) AS email,
    COALESCE(
      u.display_name,
      NULLIF(TRIM(COALESCE(u.first_name,'') || ' ' || COALESCE(u.last_name,'')), ''),
      '-'
    ) AS nama,
    CASE
      WHEN LOWER(u.email) = 'yulia.tri.m@gmail.com' THEN 'BPS Kota Tanjungpinang'
      WHEN LOWER(u.email) = 'ekoaprianto1@gmail.com' THEN 'BPS Kota Batam'
      -- Mapping akan dilengkapi dari CSV
      ELSE 'BPS Kabupaten/Kota'
    END AS instansi
  FROM users u
),

-- 4) Cross join user dengan bidang untuk Kepala Kab/Kot
user_bidang AS (
  SELECT 
    ui.user_id,
    ui.email,
    ui.nama,
    ui.instansi,
    bl.bidang
  FROM user_info ui
  CROSS JOIN bidang_list bl
)

-- 5) Final output
SELECT
  ub.nama,
  ub.email,
  ub.bidang,
  ub.instansi,
  'Ketua' AS jabatan,
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
FROM user_bidang ub
LEFT JOIN issue_agg ia ON ia.user_id = ub.user_id AND ia.bidang = ub.bidang
ORDER BY ub.instansi, ub.nama, ub.bidang;
