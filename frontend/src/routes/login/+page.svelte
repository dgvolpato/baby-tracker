<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';

	let username = '';
	let password = '';
	let error = '';
	let loading = false;

	async function submit() {
		error = '';
		loading = true;
		try {
			const { token } = await api.login(username, password);
			localStorage.setItem('token', token);
			goto('/dashboard');
		} catch {
			error = 'Invalid username or password.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="login-wrap">
	<div class="login-card card">
		<h1>Baby Tracker</h1>
		<p class="subtitle">Sign in to continue</p>

		<form on:submit|preventDefault={submit}>
			<div class="field">
				<label for="username">Username</label>
				<input id="username" type="text" bind:value={username} autocomplete="username" />
			</div>
			<div class="field">
				<label for="password">Password</label>
				<input id="password" type="password" bind:value={password} autocomplete="current-password" />
			</div>

			{#if error}
				<p class="error">{error}</p>
			{/if}

			<button type="submit" class="btn-primary" disabled={loading}>
				{loading ? 'Signing in…' : 'Sign in'}
			</button>
		</form>
	</div>
</div>

<style>
	.login-wrap {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 16px;
	}
	.login-card {
		width: 100%;
		max-width: 360px;
		padding: 40px 32px;
	}
	h1 {
		font-size: 22px;
		margin-bottom: 4px;
	}
	.subtitle {
		color: var(--text-muted);
		font-size: 14px;
		margin-bottom: 28px;
	}
	.field {
		margin-bottom: 14px;
	}
	label {
		display: block;
		font-size: 13px;
		font-weight: 500;
		margin-bottom: 5px;
	}
	.error {
		font-size: 13px;
		color: var(--danger);
		margin-bottom: 12px;
	}
	button {
		width: 100%;
		padding: 10px;
		font-size: 15px;
		margin-top: 4px;
	}
</style>
