<script lang="ts">
    const currentYear = new Date().getFullYear();

    const links = [
        { href: "https://kepri.bps.go.id", label: "BPS Kepulauan Riau" },
        { href: "https://plane.serumpun.id", label: "Portal SERUMPUN" },
    ];
</script>

<footer class="footer">
    <div class="container footer-content">
        <div class="footer-brand">
            <div class="brand-logo">
                <span class="logo-icon">📊</span>
                <span class="logo-text">SERUMPUN</span>
            </div>
            <p class="brand-desc">
                Platform Kolaborasi & Monitoring<br />
                Sensus Ekonomi 2026
            </p>
        </div>

        <div class="footer-links">
            {#each links as link}
                <a
                    href={link.href}
                    target="_blank"
                    rel="noopener"
                    class="footer-link"
                >
                    {link.label}
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="12"
                        height="12"
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
            {/each}
        </div>

        <div class="footer-bottom">
            <div class="footer-divider"></div>
            <p class="copyright">
                &copy; {currentYear} BPS Provinsi Kepulauan Riau — Sensus Ekonomi
                2026
            </p>
        </div>
    </div>
</footer>

<style>
    .footer {
        margin-top: auto;
        padding: var(--spacing-2xl) 0 var(--spacing-lg);
        border-top: 1px solid var(--color-border);
        background: var(--color-bg-secondary);
    }

    .footer-content {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xl);
    }

    .footer-brand {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-sm);
    }

    .brand-logo {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        font-size: 1.25rem;
        font-weight: 800;
    }

    .logo-icon {
        font-size: 1.5rem;
    }

    .logo-text {
        background: var(--gradient-primary);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
    }

    .brand-desc {
        color: var(--color-text-muted);
        font-size: 0.875rem;
        line-height: 1.5;
    }

    .footer-links {
        display: flex;
        flex-wrap: wrap;
        gap: var(--spacing-lg);
    }

    .footer-link {
        display: inline-flex;
        align-items: center;
        gap: var(--spacing-xs);
        color: var(--color-text-secondary);
        font-size: 0.875rem;
        transition: color var(--transition-fast);
    }

    .footer-link:hover {
        color: var(--color-primary-light);
    }

    .footer-bottom {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);
    }

    .footer-divider {
        height: 1px;
        background: var(--color-border);
    }

    .copyright {
        color: var(--color-text-muted);
        font-size: 0.8125rem;
        text-align: center;
    }

    @media (min-width: 768px) {
        .footer-content {
            flex-direction: row;
            flex-wrap: wrap;
            justify-content: space-between;
            align-items: flex-start;
        }

        .footer-brand {
            flex: 1;
        }

        .footer-bottom {
            width: 100%;
        }

        .copyright {
            text-align: left;
        }
    }
</style>
