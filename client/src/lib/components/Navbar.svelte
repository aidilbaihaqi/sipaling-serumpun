<script lang="ts">
	import { onMount } from "svelte";

	let activeDropdown: string | null = null;
	let isMobileMenuOpen = false;
	let expandedMobileMenu: string | null = null;

	// Grouped navigation menus
	const navMenus = [
		{
			id: "portal",
			label: "Portal",
			items: [
				{
					label: "Portal (All Users)",
					icon: "🌐",
					href: "https://map.gurind.am/spaces/issues/3009eed4dd9b42558207898d5d524fde/?board=list",
				},
				{
					label: "Portal (Members)",
					icon: "👥",
					href: "https://map.gurind.am/serumpun/projects/cfc12151-e169-4caf-bca9-3eb83ed588ee/issues/",
				},
			],
		},
		{
			id: "resources",
			label: "Resources",
			items: [
				{
					label: "Pendaftaran",
					icon: "📝",
					href: "https://docs.google.com/spreadsheets/d/1MElGOLZ-nH368URZ0LmLnwJrQy_2vWgxZjmiSZ5n1n8/edit?gid=0#gid=0",
				},
				{
					label: "Petunjuk",
					icon: "📖",
					href: "https://docs.google.com/presentation/d/17-JMG0H1xrfYNpYv-ZtwjBEF1wqbCvSk/edit?slide=id.p1#slide=id.p1",
				},
			],
		},
	];

	function openDropdown(id: string) {
		activeDropdown = id;
	}

	function closeDropdown() {
		activeDropdown = null;
	}

	function toggleMobileMenu() {
		isMobileMenuOpen = !isMobileMenuOpen;
		// Prevent body scroll when menu open
		document.body.style.overflow = isMobileMenuOpen ? "hidden" : "";
	}

	function toggleMobileAccordion(menuId: string) {
		expandedMobileMenu = expandedMobileMenu === menuId ? null : menuId;
	}

	function closeMobileMenu() {
		isMobileMenuOpen = false;
		document.body.style.overflow = "";
	}

	// Click outside to close desktop dropdown
	onMount(() => {
		function handleClickOutside(event: MouseEvent) {
			const target = event.target as HTMLElement;
			if (!target.closest(".nav-menu")) {
				closeDropdown();
			}
		}

		document.addEventListener("click", handleClickOutside);
		return () => {
			document.removeEventListener("click", handleClickOutside);
			document.body.style.overflow = "";
		};
	});
</script>

<!-- Mobile Backdrop -->
<div
	class="mobile-backdrop"
	class:active={isMobileMenuOpen}
	on:click={closeMobileMenu}
	on:keydown={(e) => e.key === "Escape" && closeMobileMenu()}
	role="button"
	tabindex="-1"
	aria-label="Close menu"
></div>

<nav class="navbar">
	<div class="container navbar-content">
		<!-- Logo -->
		<a href="/" class="logo">
			<img
				src="/images/logo-se2026.jfif"
				alt="Sensus Ekonomi 2026"
				class="logo-image"
			/>
			<span class="logo-text">SERUMPUN</span>
		</a>

		<!-- Desktop Dropdown Menus -->
		<div class="nav-menus">
			{#each navMenus as menu}
				<div
					class="nav-menu"
					on:mouseenter={() => openDropdown(menu.id)}
					on:mouseleave={closeDropdown}
					on:keydown={(e) => {
						if (e.key === "Enter" || e.key === " ") {
							e.preventDefault();
							openDropdown(menu.id);
						}
						if (e.key === "Escape") {
							closeDropdown();
						}
					}}
					role="button"
					tabindex="0"
					aria-haspopup="true"
					aria-expanded={activeDropdown === menu.id}
				>
					<button class="menu-trigger">
						<span>{menu.label}</span>
						<svg
							class="chevron"
							class:rotate={activeDropdown === menu.id}
							width="16"
							height="16"
							viewBox="0 0 16 16"
							fill="none"
							xmlns="http://www.w3.org/2000/svg"
						>
							<path
								d="M4 6L8 10L12 6"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
							/>
						</svg>
					</button>

					{#if activeDropdown === menu.id}
						<div class="dropdown-menu" role="menu">
							{#each menu.items as item}
								<a
									href={item.href}
									target="_blank"
									rel="noopener noreferrer"
									class="dropdown-item"
									role="menuitem"
									aria-label="Open {item.label} in new tab"
								>
									<span class="item-icon">{item.icon}</span>
									<span class="item-label">{item.label}</span>
								</a>
							{/each}
						</div>
					{/if}
				</div>
			{/each}
		</div>

		<!-- Mobile Toggle -->
		<button
			class="menu-toggle"
			on:click={toggleMobileMenu}
			aria-label="Toggle menu"
		>
			<span class="hamburger" class:active={isMobileMenuOpen}></span>
		</button>
	</div>
</nav>

<!-- Mobile Menu -->
<div class="mobile-menu" class:open={isMobileMenuOpen}>
	{#each navMenus as menu}
		<div class="mobile-menu-group">
			<button
				class="mobile-menu-trigger"
				on:click={() => toggleMobileAccordion(menu.id)}
			>
				<span>{menu.label}</span>
				<svg
					class="chevron"
					class:rotate={expandedMobileMenu === menu.id}
					width="16"
					height="16"
					viewBox="0 0 16 16"
					fill="none"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path
						d="M4 6L8 10L12 6"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
					/>
				</svg>
			</button>

			{#if expandedMobileMenu === menu.id}
				<div class="mobile-submenu">
					{#each menu.items as item}
						<a
							href={item.href}
							target="_blank"
							rel="noopener noreferrer"
							class="mobile-submenu-item"
							on:click={closeMobileMenu}
						>
							<span class="item-icon">{item.icon}</span>
							<span>{item.label}</span>
						</a>
					{/each}
				</div>
			{/if}
		</div>
	{/each}
</div>

<style>
	/* Navbar with Glassmorphism */
	.navbar {
		position: sticky;
		top: 0;
		z-index: 1001; /* ✅ FIX: Higher than mobile-menu (1000) to ensure hamburger stays visible */
		padding: var(--spacing-md) 0;
		background: rgba(255, 255, 255, 0.9);
		backdrop-filter: blur(12px) saturate(180%);
		-webkit-backdrop-filter: blur(12px) saturate(180%);
		border-bottom: 1px solid rgba(248, 132, 45, 0.15);
		box-shadow: 0 2px 16px rgba(0, 0, 0, 0.04);
	}

	.navbar-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--spacing-lg);
	}

	/* Logo */
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

	.logo-image {
		height: 40px;
		width: auto;
		object-fit: contain;
	}

	.logo-text {
		background: linear-gradient(
			to right,
			var(--color-primary),
			var(--color-accent)
		);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}

	/* Desktop Navigation Menus */
	.nav-menus {
		display: flex;
		align-items: center;
		gap: var(--spacing-lg);
	}

	.nav-menu {
		position: relative;
	}

	/* Menu Trigger */
	.menu-trigger {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		font-weight: 500;
		font-size: 1rem;
		color: var(--color-text-secondary);
		border-radius: var(--radius-md);
		transition: all var(--transition-fast);
		cursor: pointer;
		background: transparent;
		border: none;
	}

	.menu-trigger:hover {
		color: var(--color-text);
		background: rgba(248, 132, 45, 0.05);
	}

	.menu-trigger:active {
		transform: scale(0.98);
	}

	/* Chevron Icon */
	.chevron {
		transition: transform var(--transition-fast);
	}

	.chevron.rotate {
		transform: rotate(180deg);
	}

	/* Dropdown Menu */
	.dropdown-menu {
		position: absolute;
		top: calc(100% + 0.5rem);
		left: 0;
		min-width: 220px;
		background: var(--color-bg);
		border-radius: var(--radius-lg);
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
		padding: 0.5rem;
		animation: slideDown 0.2s ease;
		z-index: 1001;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-8px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	/* Dropdown Item */
	.dropdown-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		transition: all var(--transition-fast);
		text-decoration: none;
	}

	.dropdown-item:hover {
		background: linear-gradient(
			135deg,
			rgba(248, 132, 45, 0.1),
			rgba(250, 178, 40, 0.1)
		);
		color: var(--color-primary);
	}

	.dropdown-item:hover .item-icon {
		transform: scale(1.1);
	}

	.item-icon {
		font-size: 1.25rem;
		transition: transform var(--transition-fast);
	}

	.item-label {
		font-weight: 500;
	}

	/* Mobile Toggle */
	.menu-toggle {
		display: none;
		padding: var(--spacing-sm);
		background: transparent;
		border: none;
		cursor: pointer;
		z-index: 1002; /* ✅ FIX: Higher than mobile-menu (1000) and backdrop (999) */
		position: relative; /* Ensure z-index works */
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

	/* Mobile Menu */
	.mobile-menu {
		position: fixed;
		top: 0;
		right: -100%;
		width: 85%;
		max-width: 320px;
		height: 100vh;
		background: var(--color-bg);
		box-shadow: -4px 0 24px rgba(0, 0, 0, 0.1);
		padding: 5rem var(--spacing-xl); /* ✅ Increased top padding for breathing room */
		transition: right 0.3s ease;
		z-index: 1000;
		overflow-y: auto;
	}

	.mobile-menu.open {
		right: 0; /* ✅ FIX: Use right instead of transform */
	}

	/* Mobile Backdrop */
	.mobile-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		opacity: 0;
		pointer-events: none;
		transition: opacity 0.3s ease;
		z-index: 999;
	}

	.mobile-backdrop.active {
		opacity: 1;
		pointer-events: all;
	}

	/* Mobile Menu Group */
	.mobile-menu-group {
		margin-bottom: 0.5rem;
	}

	/* Mobile Menu Trigger */
	.mobile-menu-trigger {
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
		padding: 1rem;
		font-size: 1.125rem;
		font-weight: 600;
		color: var(--color-text);
		background: var(--color-bg-secondary);
		border: none;
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: background var(--transition-fast);
	}

	.mobile-menu-trigger:hover {
		background: rgba(248, 132, 45, 0.08);
	}

	.mobile-menu-trigger:active {
		transform: scale(0.98);
	}

	/* Mobile Submenu */
	.mobile-submenu {
		padding-left: 1rem;
		padding-top: 0.5rem;
		animation: slideDown 0.2s ease;
	}

	.mobile-submenu-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.875rem 1rem;
		color: var(--color-text-secondary);
		border-radius: var(--radius-md);
		margin-bottom: 0.25rem;
		text-decoration: none;
		transition: all var(--transition-fast);
	}

	.mobile-submenu-item:hover {
		background: rgba(248, 132, 45, 0.1);
		color: var(--color-primary);
	}

	/* Mobile Responsive */
	@media (max-width: 768px) {
		.navbar-content {
			padding: 0 var(--spacing-md);
		}

		.logo-image {
			height: 32px;
		}

		.nav-menus {
			display: none;
		}

		.menu-toggle {
			display: block;
		}
	}

	/* Focus Indicators for Accessibility */
	.menu-trigger:focus-visible,
	.mobile-menu-trigger:focus-visible {
		outline: 2px solid var(--color-primary);
		outline-offset: 2px;
	}

	.dropdown-item:focus-visible,
	.mobile-submenu-item:focus-visible {
		outline: 2px solid var(--color-primary);
		outline-offset: -2px;
	}
</style>
