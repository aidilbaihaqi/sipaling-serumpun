<script lang="ts">
    // Dummy embed data - simulasi data yang akan ditampilkan di canvas
    const dummyEmbeds = [
        {
            id: "1",
            title: "Progres SE2026 per Kabupaten/Kota",
            type: "chart",
            color: "#6366f1",
        },
        {
            id: "2",
            title: "KPI Penugasan",
            type: "kpi",
            color: "#22d3ee",
        },
        {
            id: "3",
            title: "Heatmap Bidang × Wilayah",
            type: "heatmap",
            color: "#a855f7",
        },
        {
            id: "4",
            title: "Tabel Detail Penugasan",
            type: "table",
            color: "#f59e0b",
        },
        {
            id: "5",
            title: "Tren Penyelesaian Mingguan",
            type: "trend",
            color: "#22c55e",
        },
        {
            id: "6",
            title: "Komentar Terbaru",
            type: "comments",
            color: "#ec4899",
        },
    ];

    // Stats dummy
    const stats = [
        { label: "Total Penugasan", value: "1,248", change: "+12%" },
        { label: "Selesai", value: "892", change: "+8%" },
        { label: "Dalam Proses", value: "286", change: "-3%" },
        { label: "Pending", value: "70", change: "-15%" },
    ];
</script>

<svelte:head>
    <title>Dashboard Monitoring - SERUMPUN</title>
    <meta
        name="description"
        content="Dashboard monitoring dan evaluasi Sensus Ekonomi 2026"
    />
</svelte:head>

<div class="dashboard">
    <header class="dashboard-header">
        <div class="container">
            <div class="header-content">
                <div class="header-text">
                    <h1 class="page-title">Dashboard Monitoring</h1>
                    <p class="page-desc">
                        Pantau progres Sensus Ekonomi 2026 secara real-time
                    </p>
                </div>

                <div class="header-meta">
                    <span class="update-badge">
                        <span class="update-dot"></span>
                        Live Update
                    </span>
                    <span class="last-update"
                        >Terakhir: {new Date().toLocaleString("id-ID")}</span
                    >
                </div>
            </div>
        </div>
    </header>

    <!-- Stats Overview -->
    <section class="stats-section">
        <div class="container">
            <div class="stats-grid">
                {#each stats as stat}
                    <div class="stat-card card">
                        <span class="stat-label">{stat.label}</span>
                        <span class="stat-value">{stat.value}</span>
                        <span
                            class="stat-change"
                            class:positive={stat.change.startsWith("+")}
                            class:negative={stat.change.startsWith("-")}
                        >
                            {stat.change}
                        </span>
                    </div>
                {/each}
            </div>
        </div>
    </section>

    <!-- Canvas Area -->
    <section class="canvas-section">
        <div class="container">
            <div class="canvas-header">
                <h2 class="canvas-title">📊 Canvas Visualisasi</h2>
                <p class="canvas-desc">
                    Area untuk menampilkan embed dari berbagai platform
                </p>
            </div>

            <div class="canvas-area">
                <div class="embed-grid">
                    {#each dummyEmbeds as embed, i}
                        <div
                            class="embed-card card"
                            style="--accent-color: {embed.color}; animation-delay: {i *
                                100}ms"
                        >
                            <div class="embed-header">
                                <span class="embed-type">{embed.type}</span>
                                <button class="embed-menu" aria-label="Options">
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        width="16"
                                        height="16"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    >
                                        <circle cx="12" cy="12" r="1" />
                                        <circle cx="12" cy="5" r="1" />
                                        <circle cx="12" cy="19" r="1" />
                                    </svg>
                                </button>
                            </div>

                            <div class="embed-content">
                                <div class="embed-placeholder">
                                    <div class="placeholder-icon">
                                        {#if embed.type === "chart"}
                                            📊
                                        {:else if embed.type === "kpi"}
                                            📈
                                        {:else if embed.type === "heatmap"}
                                            🗺️
                                        {:else if embed.type === "table"}
                                            📋
                                        {:else if embed.type === "trend"}
                                            📉
                                        {:else}
                                            💬
                                        {/if}
                                    </div>
                                    <p class="placeholder-text">
                                        Embed Preview
                                    </p>
                                    <span class="placeholder-hint"
                                        >Flourish / Data Studio</span
                                    >
                                </div>
                            </div>

                            <div class="embed-footer">
                                <h3 class="embed-title">{embed.title}</h3>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>

            <!-- Empty State - ditampilkan jika tidak ada embed -->
            <!-- <div class="canvas-empty">
        <div class="empty-icon">📊</div>
        <h3 class="empty-title">Canvas Kosong</h3>
        <p class="empty-desc">Tambahkan embed dari Flourish atau platform visualisasi lainnya</p>
        <button class="btn btn-primary">
          + Tambah Embed
        </button>
      </div> -->
        </div>
    </section>
</div>

<style>
    .dashboard {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xl);
        padding-bottom: var(--spacing-2xl);
    }

    /* Header */
    .dashboard-header {
        padding: var(--spacing-xl) 0;
        border-bottom: 1px solid var(--color-border);
        background: var(--color-bg-secondary);
    }

    .header-content {
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        align-items: flex-start;
        gap: var(--spacing-lg);
    }

    .page-title {
        font-size: 1.75rem;
        margin-bottom: var(--spacing-xs);
    }

    .page-desc {
        color: var(--color-text-muted);
        font-size: 0.9375rem;
    }

    .header-meta {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        gap: var(--spacing-xs);
    }

    .update-badge {
        display: inline-flex;
        align-items: center;
        gap: var(--spacing-xs);
        padding: var(--spacing-xs) var(--spacing-sm);
        background: rgba(34, 197, 94, 0.15);
        border: 1px solid rgba(34, 197, 94, 0.3);
        border-radius: var(--radius-full);
        font-size: 0.75rem;
        font-weight: 600;
        color: var(--color-success);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .update-dot {
        width: 6px;
        height: 6px;
        background: var(--color-success);
        border-radius: 50%;
        animation: pulse 2s infinite;
    }

    .last-update {
        font-size: 0.8125rem;
        color: var(--color-text-muted);
    }

    /* Stats Section */
    .stats-section {
        padding: var(--spacing-lg) 0 0;
    }

    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
        gap: var(--spacing-md);
    }

    .stat-card {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xs);
        padding: var(--spacing-lg);
    }

    .stat-label {
        font-size: 0.8125rem;
        color: var(--color-text-muted);
        font-weight: 500;
    }

    .stat-value {
        font-size: 1.75rem;
        font-weight: 700;
        color: var(--color-text);
    }

    .stat-change {
        font-size: 0.8125rem;
        font-weight: 600;
    }

    .stat-change.positive {
        color: var(--color-success);
    }

    .stat-change.negative {
        color: var(--color-error);
    }

    /* Canvas Section */
    .canvas-section {
        flex: 1;
    }

    .canvas-header {
        margin-bottom: var(--spacing-lg);
    }

    .canvas-title {
        font-size: 1.25rem;
        margin-bottom: var(--spacing-xs);
    }

    .canvas-desc {
        color: var(--color-text-muted);
        font-size: 0.875rem;
    }

    .canvas-area {
        background: var(--color-bg-secondary);
        border: 2px dashed var(--color-border);
        border-radius: var(--radius-xl);
        padding: var(--spacing-lg);
        min-height: 500px;
    }

    .embed-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
        gap: var(--spacing-lg);
    }

    /* Embed Card */
    .embed-card {
        display: flex;
        flex-direction: column;
        overflow: hidden;
        padding: 0;
        border-left: 3px solid var(--accent-color);
        animation: fadeIn 0.5s ease forwards;
        opacity: 0;
    }

    .embed-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: var(--spacing-sm) var(--spacing-md);
        border-bottom: 1px solid var(--color-border);
        background: rgba(0, 0, 0, 0.2);
    }

    .embed-type {
        font-size: 0.6875rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: var(--accent-color);
        padding: 2px 8px;
        background: rgba(255, 255, 255, 0.05);
        border-radius: var(--radius-sm);
    }

    .embed-menu {
        padding: var(--spacing-xs);
        color: var(--color-text-muted);
        border-radius: var(--radius-sm);
        transition: all var(--transition-fast);
    }

    .embed-menu:hover {
        color: var(--color-text);
        background: var(--color-bg-card);
    }

    .embed-content {
        flex: 1;
        min-height: 180px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(
            135deg,
            rgba(0, 0, 0, 0.1) 0%,
            transparent 100%
        );
    }

    .embed-placeholder {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: var(--spacing-sm);
        text-align: center;
        padding: var(--spacing-lg);
    }

    .placeholder-icon {
        font-size: 2.5rem;
        opacity: 0.7;
    }

    .placeholder-text {
        font-size: 0.9375rem;
        font-weight: 500;
        color: var(--color-text-secondary);
    }

    .placeholder-hint {
        font-size: 0.75rem;
        color: var(--color-text-muted);
    }

    .embed-footer {
        padding: var(--spacing-md);
        border-top: 1px solid var(--color-border);
    }

    .embed-title {
        font-size: 0.9375rem;
        font-weight: 600;
        color: var(--color-text);
    }

    /* Empty State */
    .canvas-empty {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 400px;
        text-align: center;
        gap: var(--spacing-md);
    }

    .empty-icon {
        font-size: 4rem;
        opacity: 0.5;
    }

    .empty-title {
        font-size: 1.25rem;
        color: var(--color-text);
    }

    .empty-desc {
        color: var(--color-text-muted);
        max-width: 300px;
    }

    @media (max-width: 768px) {
        .header-meta {
            align-items: flex-start;
        }

        .embed-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
