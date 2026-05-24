<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Feeding } from '$lib/api';
	import { dayLabel, fmtOz, toLocalInputDate } from '$lib/utils';

	interface DayTotal {
		date: string;
		label: string;
		total: number;
		formula: number;
		breast: number;
		count: number;
	}

	let feedings: Feeding[] = [];
	let loading = true;
	let error = '';

	$: dayTotals = computeTotals(feedings);

	function computeTotals(list: Feeding[]): DayTotal[] {
		const map = new Map<string, DayTotal>();
		for (const f of list) {
			const key = toLocalInputDate(new Date(f.timestamp));
			if (!map.has(key)) {
				map.set(key, { date: key, label: dayLabel(key), total: 0, formula: 0, breast: 0, count: 0 });
			}
			const day = map.get(key)!;
			day.total = Math.round((day.total + f.oz) * 10) / 10;
			day.count++;
			if (f.type === 'formula') day.formula = Math.round((day.formula + f.oz) * 10) / 10;
			else day.breast = Math.round((day.breast + f.oz) * 10) / 10;
		}
		return Array.from(map.values()).sort((a, b) => b.date.localeCompare(a.date));
	}

	onMount(async () => {
		try {
			feedings = await api.feedings.list();
		} catch {
			error = 'Failed to load feedings.';
		} finally {
			loading = false;
		}
	});
</script>

<div class="page">
	<h1>Dashboard</h1>

	{#if loading}
		<p class="muted">Loading…</p>
	{:else if error}
		<p class="error-text">{error}</p>
	{:else if dayTotals.length === 0}
		<p class="muted">No feedings logged yet. <a href="/history">Add one →</a></p>
	{:else}
		<div class="list">
			{#each dayTotals as day (day.date)}
				<div class="day-card card">
					<div class="day-header">
						<span class="day-label">{day.label}</span>
						<span class="day-total">{fmtOz(day.total)}</span>
					</div>
					<div class="day-meta">
						{#if day.formula > 0}
							<span class="tag formula">Formula {fmtOz(day.formula)}</span>
						{/if}
						{#if day.breast > 0}
							<span class="tag breast">Breast {fmtOz(day.breast)}</span>
						{/if}
						<span class="count">{day.count} feeding{day.count !== 1 ? 's' : ''}</span>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	h1 {
		font-size: 22px;
		margin-bottom: 20px;
	}
	.error-text {
		color: var(--danger);
		font-size: 14px;
	}
	.list {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}
	.day-card {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}
	.day-header {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
	}
	.day-label {
		font-weight: 600;
		font-size: 16px;
	}
	.day-total {
		font-size: 20px;
		font-weight: 700;
		color: var(--accent);
	}
	.day-meta {
		display: flex;
		align-items: center;
		gap: 8px;
		flex-wrap: wrap;
	}
	.tag {
		font-size: 12px;
		padding: 2px 8px;
		border-radius: 12px;
		font-weight: 500;
	}
	.formula {
		background: #e8f0fe;
		color: #1a56db;
	}
	.breast {
		background: #fde8f0;
		color: #b5477a;
	}
	.count {
		font-size: 12px;
		color: var(--text-muted);
		margin-left: auto;
	}
</style>
