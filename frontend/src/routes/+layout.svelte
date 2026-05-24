<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		if (!localStorage.getItem('token') && $page.url.pathname !== '/login') {
			goto('/login');
		}
	});

	function logout() {
		localStorage.removeItem('token');
		goto('/login');
	}
</script>

{#if $page.url.pathname !== '/login'}
	<header>
		<nav>
			<span class="brand">Baby Tracker</span>
			<div class="links">
				<a href="/dashboard" class:active={$page.url.pathname === '/dashboard'}>Dashboard</a>
				<a href="/history" class:active={$page.url.pathname === '/history'}>Feedings</a>
				<a href="/diapers" class:active={$page.url.pathname === '/diapers'}>Diapers</a>
				<a href="/measurements" class:active={$page.url.pathname === '/measurements'}>Measurements</a>
				<button class="btn-ghost logout" on:click={logout}>Log out</button>
			</div>
		</nav>
	</header>
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
		z-index: 10;
	}
	nav {
		max-width: 640px;
		margin: 0 auto;
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 56px;
	}
	.brand {
		font-weight: 600;
		font-size: 17px;
	}
	.links {
		display: flex;
		align-items: center;
		gap: 4px;
	}
	.links a {
		color: var(--text-muted);
		font-size: 14px;
		padding: 6px 10px;
		border-radius: var(--radius);
		transition: background 0.15s, color 0.15s;
	}
	.links a:hover:not(.active) {
		background: var(--border);
		color: var(--text);
	}
	.links a.active {
		color: var(--accent);
		background: var(--accent-light);
	}
	.logout {
		font-size: 13px;
		padding: 5px 10px;
		margin-left: 4px;
	}
</style>
