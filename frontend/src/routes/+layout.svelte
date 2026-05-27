<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let menuOpen = false;

	onMount(() => {
		if (!localStorage.getItem('token') && $page.url.pathname !== '/login') {
			goto('/login');
		}
	});

	function logout() {
		localStorage.removeItem('token');
		goto('/login');
	}

	function close() {
		menuOpen = false;
	}

	function navigate(path: string) {
		close();
		goto(path);
	}
</script>

{#if $page.url.pathname !== '/login'}
	<header>
		<span class="brand">Baby Tracker</span>
		<div class="header-right">
			<a href="/dashboard" class="shortcut" class:active={$page.url.pathname === '/dashboard'} aria-label="Dashboard">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
					<rect x="3" y="3" width="7" height="7" rx="1"/>
					<rect x="14" y="3" width="7" height="7" rx="1"/>
					<rect x="3" y="14" width="7" height="7" rx="1"/>
					<rect x="14" y="14" width="7" height="7" rx="1"/>
				</svg>
			</a>
			<a href="/history" class="shortcut" class:active={$page.url.pathname === '/history'} aria-label="Feedings">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
					<g transform="rotate(45, 12, 12)">
						<rect x="10.5" y="1" width="3" height="3" rx="1.5"/>
						<rect x="8.5" y="4" width="7" height="2.5" rx="1"/>
						<path d="M8.5 6.5 C6 8 5.5 10 5.5 12.5 L5.5 18.5 Q5.5 22.5 12 22.5 Q18.5 22.5 18.5 18.5 L18.5 12.5 C18.5 10 18 8 15.5 6.5 Z"/>
					</g>
				</svg>
			</a>
			<button
				class="hamburger"
				class:open={menuOpen}
				aria-label="Open menu"
				aria-expanded={menuOpen}
				on:click={() => (menuOpen = !menuOpen)}
			>
				<span></span>
				<span></span>
				<span></span>
			</button>
		</div>
	</header>

	{#if menuOpen}
		<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
		<div class="drawer-overlay" on:click|self={close}>
			<div class="drawer-backdrop" on:click={close}></div>
			<nav class="drawer">
				<a
					href="/dashboard"
					class:active={$page.url.pathname === '/dashboard'}
					on:click|preventDefault={() => navigate('/dashboard')}
				>
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
						<rect x="3" y="3" width="7" height="7" rx="1"/>
						<rect x="14" y="3" width="7" height="7" rx="1"/>
						<rect x="3" y="14" width="7" height="7" rx="1"/>
						<rect x="14" y="14" width="7" height="7" rx="1"/>
					</svg>
					Dashboard
				</a>
				<a
					href="/history"
					class:active={$page.url.pathname === '/history'}
					on:click|preventDefault={() => navigate('/history')}
				>
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
						<g transform="rotate(45, 12, 12)">
							<rect x="10.5" y="1" width="3" height="3" rx="1.5"/>
							<rect x="8.5" y="4" width="7" height="2.5" rx="1"/>
							<path d="M8.5 6.5 C6 8 5.5 10 5.5 12.5 L5.5 18.5 Q5.5 22.5 12 22.5 Q18.5 22.5 18.5 18.5 L18.5 12.5 C18.5 10 18 8 15.5 6.5 Z"/>
						</g>
					</svg>
					Feedings
				</a>
				<a
					href="/diapers"
					class:active={$page.url.pathname === '/diapers'}
					on:click|preventDefault={() => navigate('/diapers')}
				>
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
						<path d="M4 8h16l-2 9H6L4 8z"/>
						<path d="M4 8c0-2 1.5-4 4-4h8c2.5 0 4 2 4 4"/>
						<path d="M9 8c0 2 1.5 3 3 3s3-1 3-3"/>
					</svg>
					Diapers
				</a>
				<a
					href="/measurements"
					class:active={$page.url.pathname === '/measurements'}
					on:click|preventDefault={() => navigate('/measurements')}
				>
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
						<line x1="3" y1="12" x2="21" y2="12"/>
						<line x1="3" y1="6" x2="6" y2="9"/>
						<line x1="3" y1="18" x2="6" y2="15"/>
						<line x1="21" y1="6" x2="18" y2="9"/>
						<line x1="21" y1="18" x2="18" y2="15"/>
						<line x1="9" y1="9" x2="9" y2="15"/>
						<line x1="12" y1="9" x2="12" y2="15"/>
						<line x1="15" y1="9" x2="15" y2="15"/>
					</svg>
					Measurements
				</a>
				<div class="drawer-divider"></div>
				<button class="logout-btn" on:click={logout}>
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
						<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
						<polyline points="16 17 21 12 16 7"/>
						<line x1="21" y1="12" x2="9" y2="12"/>
					</svg>
					Log out
				</button>
			</nav>
		</div>
	{/if}
{/if}

<main>
	<slot />
</main>

<style>
	header {
		background: var(--surface);
		border-bottom: 1px solid var(--border);
		padding: 0 16px;
		position: sticky;
		top: 0;
		z-index: 20;
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 52px;
	}
	.brand {
		font-weight: 600;
		font-size: 17px;
	}
	.header-right {
		display: flex;
		align-items: center;
		gap: 2px;
	}
	.shortcut {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 36px;
		height: 36px;
		border-radius: var(--radius);
		color: var(--text-muted);
		transition: background 0.15s, color 0.15s;
	}
	.shortcut :global(svg) {
		width: 20px;
		height: 20px;
	}
	.shortcut:hover {
		background: var(--border);
		color: var(--text);
	}
	.shortcut.active {
		color: var(--accent);
		background: var(--accent-light);
	}

	/* Hamburger button */
	.hamburger {
		background: none;
		border: none;
		cursor: pointer;
		padding: 6px;
		border-radius: var(--radius);
		color: var(--text);
		display: flex;
		flex-direction: column;
		gap: 5px;
		transition: background 0.15s;
	}
	.hamburger:hover {
		background: var(--border);
	}
	.hamburger span {
		display: block;
		width: 22px;
		height: 2px;
		background: currentColor;
		border-radius: 2px;
		transition: transform 0.2s, opacity 0.2s;
		transform-origin: center;
	}
	.hamburger.open span:nth-child(1) { transform: translateY(7px) rotate(45deg); }
	.hamburger.open span:nth-child(2) { opacity: 0; transform: scaleX(0); }
	.hamburger.open span:nth-child(3) { transform: translateY(-7px) rotate(-45deg); }

	/* Drawer overlay */
	.drawer-overlay {
		position: fixed;
		inset: 52px 0 0 0;
		z-index: 15;
	}
	.drawer-backdrop {
		position: absolute;
		inset: 0;
		background: rgba(0, 0, 0, 0.25);
	}
	.drawer {
		position: absolute;
		top: 0;
		right: 0;
		width: 220px;
		background: var(--surface);
		border-left: 1px solid var(--border);
		border-bottom: 1px solid var(--border);
		border-radius: 0 0 0 var(--radius);
		padding: 8px 0;
		box-shadow: -4px 4px 16px rgba(0, 0, 0, 0.1);
	}
	.drawer a,
	.drawer button {
		display: flex;
		align-items: center;
		gap: 12px;
		width: 100%;
		padding: 13px 20px;
		text-decoration: none;
		color: var(--text);
		font-size: 15px;
		background: none;
		border: none;
		cursor: pointer;
		text-align: left;
		transition: background 0.12s;
	}
	.drawer a:hover,
	.drawer button:hover {
		background: var(--bg);
	}
	.drawer a.active {
		color: var(--accent);
		background: var(--accent-light);
		font-weight: 500;
	}
	.drawer a :global(svg),
	.drawer button :global(svg) {
		width: 18px;
		height: 18px;
		flex-shrink: 0;
	}
	.drawer-divider {
		height: 1px;
		background: var(--border);
		margin: 6px 0;
	}
	.logout-btn {
		color: var(--text-muted);
		font-size: 14px !important;
	}
</style>
