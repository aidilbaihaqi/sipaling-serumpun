<script lang="ts">
    const features = [
        {
            icon: "📊",
            title: "Monitoring Real-time",
            desc: "Pantau progres SE2026 secara langsung",
        },
        {
            icon: "🗺️",
            title: "Per Kabupaten/Kota",
            desc: "Detail per wilayah Kepulauan Riau",
        },
        {
            icon: "📈",
            title: "Visualisasi Data",
            desc: "Grafik dan chart interaktif",
        },
        {
            icon: "👥",
            title: "Kolaborasi Tim",
            desc: "Koordinasi antar bidang mudah",
        },
    ];

    const quickLinks = [
        {
            href: "https://plane.serumpun.id",
            label: "Portal SERUMPUN (All)",
            icon: "🌐",
        },
        {
            href: "https://plane.serumpun.id/profile",
            label: "Portal Member",
            icon: "👤",
        },
        {
            href: "https://docs.serumpun.id",
            label: "Petunjuk Penggunaan",
            icon: "📖",
        },
    ];
</script>

<svelte:head>
    <title>SERUMPUN - Platform Monitoring SE2026 Kepulauan Riau</title>
    <meta
        name="description"
        content="Platform kolaborasi dan monitoring kegiatan Sensus Ekonomi 2026 di Provinsi Kepulauan Riau"
    />
</svelte:head>

<section class="hero">
    <div class="container hero-content">
        <div class="hero-badge animate-fade-in">
            <span class="badge-dot"></span>
            Sensus Ekonomi 2026
        </div>

        <h1 class="hero-title animate-fade-in">
            Platform Monitoring<br />
            <span class="text-gradient">SERUMPUN</span>
        </h1>

        <p class="hero-desc animate-fade-in">
            Kolaborasi dan monitoring kegiatan Sensus Ekonomi 2026<br />
            di Provinsi Kepulauan Riau
        </p>

        <div class="hero-actions animate-fade-in">
            <a href="/dashboard" class="btn btn-primary btn-lg">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <rect x="3" y="3" width="7" height="7" />
                    <rect x="14" y="3" width="7" height="7" />
                    <rect x="14" y="14" width="7" height="7" />
                    <rect x="3" y="14" width="7" height="7" />
                </svg>
                Lihat Dashboard
            </a>
            <a
                href="https://plane.serumpun.id"
                target="_blank"
                rel="noopener"
                class="btn btn-outline btn-lg"
            >
                Buka Portal
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
                    <path
                        d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"
                    />
                    <polyline points="15 3 21 3 21 9" />
                    <line x1="10" y1="14" x2="21" y2="3" />
                </svg>
            </a>
        </div>
    </div>

    <div class="hero-glow"></div>
</section>

<section class="features">
    <div class="container">
        <h2 class="section-title">Fitur Utama</h2>

        <div class="features-grid">
            {#each features as feature, i}
                <div
                    class="feature-card card"
                    style="animation-delay: {i * 100}ms"
                >
                    <span class="feature-icon">{feature.icon}</span>
                    <h3 class="feature-title">{feature.title}</h3>
                    <p class="feature-desc">{feature.desc}</p>
                </div>
            {/each}
        </div>
    </div>
</section>

<section class="quick-access">
    <div class="container">
        <h2 class="section-title">Akses Cepat</h2>

        <div class="links-grid">
            {#each quickLinks as link}
                <a
                    href={link.href}
                    target="_blank"
                    rel="noopener"
                    class="quick-link card"
                >
                    <span class="link-icon">{link.icon}</span>
                    <span class="link-label">{link.label}</span>
                    <svg
                        class="link-arrow"
                        xmlns="http://www.w3.org/2000/svg"
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <line x1="5" y1="12" x2="19" y2="12" />
                        <polyline points="12 5 19 12 12 19" />
                    </svg>
                </a>
            {/each}
        </div>
    </div>
</section>

<style>
    /* Hero Section */
    .hero {
        position: relative;
        padding: var(--spacing-3xl) 0;
        min-height: 70vh;
        display: flex;
        align-items: center;
        overflow: hidden;
    }

    .hero-content {
        display: flex;
        flex-direction: column;
        align-items: center;
        text-align: center;
        gap: var(--spacing-lg);
        position: relative;
        z-index: 1;
    }

    .hero-badge {
        display: inline-flex;
        align-items: center;
        gap: var(--spacing-sm);
        padding: var(--spacing-xs) var(--spacing-md);
        background: var(--color-bg-card);
        border: 1px solid var(--color-border);
        border-radius: var(--radius-full);
        font-size: 0.875rem;
        font-weight: 500;
        color: var(--color-text-secondary);
    }

    .badge-dot {
        width: 8px;
        height: 8px;
        background: var(--color-success);
        border-radius: 50%;
        animation: pulse 2s infinite;
    }

    .hero-title {
        font-size: clamp(2.5rem, 6vw, 4rem);
        font-weight: 800;
        line-height: 1.1;
        letter-spacing: -0.02em;
    }

    .hero-desc {
        font-size: 1.125rem;
        color: var(--color-text-secondary);
        max-width: 500px;
        line-height: 1.6;
    }

    .hero-actions {
        display: flex;
        gap: var(--spacing-md);
        flex-wrap: wrap;
        justify-content: center;
        margin-top: var(--spacing-md);
    }

    .btn-lg {
        padding: var(--spacing-md) var(--spacing-xl);
        font-size: 1rem;
    }

    .hero-glow {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 600px;
        height: 600px;
        background: radial-gradient(
            circle,
            rgba(99, 102, 241, 0.2) 0%,
            transparent 70%
        );
        pointer-events: none;
    }

    /* Features Section */
    .features {
        padding: var(--spacing-3xl) 0;
    }

    .section-title {
        text-align: center;
        margin-bottom: var(--spacing-2xl);
        font-size: 1.75rem;
    }

    .features-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
        gap: var(--spacing-lg);
    }

    .feature-card {
        display: flex;
        flex-direction: column;
        align-items: center;
        text-align: center;
        padding: var(--spacing-xl);
        animation: fadeIn 0.5s ease forwards;
        opacity: 0;
    }

    .feature-icon {
        font-size: 2.5rem;
        margin-bottom: var(--spacing-md);
    }

    .feature-title {
        font-size: 1.125rem;
        font-weight: 600;
        margin-bottom: var(--spacing-sm);
    }

    .feature-desc {
        font-size: 0.875rem;
        color: var(--color-text-muted);
    }

    /* Quick Access Section */
    .quick-access {
        padding: var(--spacing-2xl) 0 var(--spacing-3xl);
    }

    .links-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: var(--spacing-md);
    }

    .quick-link {
        display: flex;
        align-items: center;
        gap: var(--spacing-md);
        padding: var(--spacing-lg);
    }

    .link-icon {
        font-size: 1.5rem;
    }

    .link-label {
        flex: 1;
        font-weight: 500;
    }

    .link-arrow {
        color: var(--color-text-muted);
        transition: transform var(--transition-fast);
    }

    .quick-link:hover .link-arrow {
        transform: translateX(4px);
        color: var(--color-primary-light);
    }
</style>
