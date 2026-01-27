<script lang="ts">
    import { page } from "$app/stores";

    let isMenuOpen = false;

    const navLinks = [
        { href: "/", label: "Beranda" },
        { href: "/dashboard", label: "Dashboard" },
    ];

    function toggleMenu() {
        isMenuOpen = !isMenuOpen;
    }
</script>

<nav class="navbar glass">
    <div class="container navbar-content">
        <a href="/" class="logo">
            <span class="logo-icon">📊</span>
            <span class="logo-text">SERUMPUN</span>
        </a>

        <button
            class="menu-toggle"
            on:click={toggleMenu}
            aria-label="Toggle menu"
        >
            <span class="hamburger" class:active={isMenuOpen}></span>
        </button>

        <div class="nav-links" class:open={isMenuOpen}>
            {#each navLinks as link}
                <a
                    href={link.href}
                    class="nav-link"
                    class:active={$page.url.pathname === link.href}
                    on:click={() => (isMenuOpen = false)}
                >
                    {link.label}
                </a>
            {/each}

            <a
                href="https://plane.serumpun.id"
                target="_blank"
                rel="noopener"
                class="btn btn-primary nav-cta"
            >
                Portal SERUMPUN
            </a>
        </div>
    </div>
</nav>

<style>
    .navbar {
        position: sticky;
        top: 0;
        z-index: 100;
        padding: var(--spacing-md) 0;
    }

    .navbar-content {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: var(--spacing-lg);
    }

    .logo {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        font-size: 1.25rem;
        font-weight: 800;
        color: var(--color-text);
        transition: transform var(--transition-fast);
    }

    .logo:hover {
        transform: scale(1.02);
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

    .nav-links {
        display: flex;
        align-items: center;
        gap: var(--spacing-xl);
    }

    .nav-link {
        font-weight: 500;
        color: var(--color-text-secondary);
        transition: color var(--transition-fast);
        position: relative;
    }

    .nav-link:hover,
    .nav-link.active {
        color: var(--color-text);
    }

    .nav-link.active::after {
        content: "";
        position: absolute;
        bottom: -4px;
        left: 0;
        right: 0;
        height: 2px;
        background: var(--gradient-primary);
        border-radius: var(--radius-full);
    }

    .nav-cta {
        padding: var(--spacing-sm) var(--spacing-lg);
    }

    .menu-toggle {
        display: none;
        padding: var(--spacing-sm);
        z-index: 101;
    }

    .hamburger {
        display: block;
        width: 24px;
        height: 2px;
        background: var(--color-text);
        position: relative;
        transition: background var(--transition-fast);
    }

    .hamburger::before,
    .hamburger::after {
        content: "";
        position: absolute;
        left: 0;
        width: 100%;
        height: 2px;
        background: var(--color-text);
        transition: all var(--transition-fast);
    }

    .hamburger::before {
        top: -8px;
    }
    .hamburger::after {
        bottom: -8px;
    }

    .hamburger.active {
        background: transparent;
    }

    .hamburger.active::before {
        top: 0;
        transform: rotate(45deg);
    }

    .hamburger.active::after {
        bottom: 0;
        transform: rotate(-45deg);
    }

    @media (max-width: 768px) {
        .menu-toggle {
            display: block;
        }

        .nav-links {
            position: fixed;
            top: 0;
            right: 0;
            bottom: 0;
            width: 280px;
            flex-direction: column;
            justify-content: center;
            gap: var(--spacing-xl);
            padding: var(--spacing-2xl);
            background: var(--color-bg-secondary);
            border-left: 1px solid var(--color-border);
            transform: translateX(100%);
            transition: transform var(--transition-normal);
        }

        .nav-links.open {
            transform: translateX(0);
        }

        .nav-cta {
            width: 100%;
            text-align: center;
        }
    }
</style>
