# SERUMPUN – Frontend (Client)

Frontend SERUMPUN merupakan **portal utama** yang menyediakan:
- Landing Page informatif
- Ringkasan dashboard (overview)
- Halaman dashboard monitoring lengkap
- Akses ke seluruh layanan SERUMPUN

Frontend dibangun menggunakan **Next.js + TypeScript (TSX)** dan
menggunakan **JokoUI components** untuk konsistensi UI.

---

## 🎯 Tujuan Frontend

- Menjadi **gerbang utama** platform SERUMPUN
- Memberikan **informasi ringkas** kondisi SE2026
- Menyediakan UX yang mudah untuk:
  - pimpinan
  - koordinator bidang
  - pengguna umum
- Menyajikan visualisasi data melalui **embed Flourish**

---

## 🧭 Konsep Halaman

### 1️⃣ Landing Page (`/`)
Fungsi:
- Pengantar platform SERUMPUN
- Ringkasan progres (KPI & overview)
- Akses cepat ke:
  - Portal SERUMPUN (All)
  - Portal SERUMPUN (Member)
  - Pendaftaran Pengguna
  - Petunjuk Penggunaan
- CTA menuju Dashboard Lengkap

Landing page **tidak menampilkan data detail**, hanya ringkasan.

---

### 2️⃣ Dashboard Page (`/dashboard`)
Fungsi:
- Monitoring & evaluasi mendalam
- Visualisasi lengkap:
  - KPI
  - Progres per kab/kota
  - Progres per bidang
  - Heatmap
  - Tabel detail + komentar
- Filter interaktif (melalui Flourish)

---

## 📁 Struktur Folder
```
client/
├── app/ # Next.js App Router
│ ├── page.tsx # Landing Page
│ ├── dashboard/
│ │ └── page.tsx # Dashboard Page
│ └── layout.tsx
├── components/
│ ├── Navbar.tsx
│ ├── Footer.tsx
│ ├── OverviewCards.tsx
│ └── FlourishEmbed.tsx
├── lib/
│ └── config.ts # Link & konfigurasi
├── styles/
├── public/
└── README.md
```

---

## 🧱 Teknologi

- Next.js (App Router)
- TypeScript (TSX)
- JokoUI Components
- CSS / Tailwind (sesuai setup)
- Flourish Embed (iframe)

---

## 📊 Integrasi Dashboard (Flourish)

Visualisasi tidak dibuat di frontend,
melainkan di **Flourish** dan di-*embed* menggunakan iframe.

Contoh komponen embed:

```tsx
<iframe
  src="https://public.flourish.studio/story/XXXXX/"
  width="100%"
  height="800"
  frameBorder="0"
  loading="lazy"
/>

// lib/config.ts
export const LINKS = {
  portalAll: "...",
  portalMember: "...",
  pendaftaran: "...",
  petunjuk: "...",
  dashboardEmbed: "https://public.flourish.studio/..."
};
```

🧠 Prinsip UX

- Informasi singkat di landing page
- Data detail hanya di dashboard
- Mobile-friendly
- Minim scroll berlebihan
- Fokus pada keterbacaan data

🚀 Pengembangan Selanjutnya
- Auth / role-based access
- Mode tampilan khusus pimpinan
- Integrasi API backend untuk KPI di landing page
- Dark mode (opsional)

© 2025 – BPS Provinsi Kepulauan Riau