<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Measurement } from '$lib/api';
	import { formatDateTime, timeOptions, roundedTimeValue, toLocalInputDate } from '$lib/utils';

	let measurements: Measurement[] = [];
	let loading = true;
	let error = '';

	let showModal = false;
	let editingId: string | null = null;
	let saving = false;

	interface FormState {
		date: string;
		time: string;
		weight: string;
		height: string;
	}

	function freshForm(): FormState {
		const now = new Date();
		return { date: toLocalInputDate(now), time: roundedTimeValue(now), weight: '', height: '' };
	}

	let form: FormState = freshForm();
	const times = timeOptions();

	function openAdd() {
		editingId = null;
		form = freshForm();
		showModal = true;
	}

	function openEdit(m: Measurement) {
		editingId = m.id;
		const dt = new Date(m.timestamp);
		form = {
			date: toLocalInputDate(dt),
			time: `${String(dt.getHours()).padStart(2, '0')}:${String(dt.getMinutes()).padStart(2, '0')}`,
			weight: m.weightLbs != null ? String(m.weightLbs) : '',
			height: m.heightIn != null ? String(m.heightIn) : ''
		};
		showModal = true;
	}

	function closeModal() {
		showModal = false;
	}

	$: hasValue = form.weight.trim() !== '' || form.height.trim() !== '';

	async function save() {
		if (!hasValue) return;
		saving = true;

		const timestamp = new Date(`${form.date}T${form.time}:00`).toISOString();
		const weightLbs = form.weight.trim() !== '' ? parseFloat(form.weight) : undefined;
		const heightIn = form.height.trim() !== '' ? parseFloat(form.height) : undefined;
		const data = { timestamp, weightLbs, heightIn };

		try {
			if (editingId) {
				const updated = await api.measurements.update(editingId, data);
				measurements = measurements.map((m) => (m.id === editingId ? updated : m));
			} else {
				const created = await api.measurements.create(data);
				measurements = [created, ...measurements];
			}
			measurements = [...measurements].sort((a, b) => b.timestamp.localeCompare(a.timestamp));
			closeModal();
		} catch {
			// keep modal open on error
		} finally {
			saving = false;
		}
	}

	async function remove(id: string) {
		if (!confirm('Delete this measurement?')) return;
		await api.measurements.delete(id);
		measurements = measurements.filter((m) => m.id !== id);
	}

	function fmt(value: number | undefined, unit: string): string {
		if (value == null) return '—';
		return `${value % 1 === 0 ? value : value.toFixed(1)} ${unit}`;
	}

	onMount(async () => {
		try {
			measurements = await api.measurements.list();
		} catch {
			error = 'Failed to load measurements.';
		} finally {
			loading = false;
		}
	});
</script>

<div class="page">
	<div class="page-header">
		<h1>Measurements</h1>
		<button class="add-btn btn-primary" on:click={openAdd} aria-label="Add measurement">+</button>
	</div>

	{#if loading}
		<p class="muted">Loading…</p>
	{:else if error}
		<p class="error-text">{error}</p>
	{:else if measurements.length === 0}
		<p class="muted">No measurements yet. Tap + to add one.</p>
	{:else}
		<div class="list">
			{#each measurements as m (m.id)}
				<div class="row card">
					<div class="row-left">
						<span class="row-time">{formatDateTime(m.timestamp)}</span>
						<div class="row-values">
							{#if m.weightLbs != null}
								<span class="val-chip weight">{fmt(m.weightLbs, 'lbs')}</span>
							{/if}
							{#if m.heightIn != null}
								<span class="val-chip height">{fmt(m.heightIn, 'in')}</span>
							{/if}
						</div>
					</div>
					<div class="row-right">
						<button class="btn-ghost btn-sm" on:click={() => openEdit(m)}>Edit</button>
						<button class="btn-danger btn-sm" on:click={() => remove(m.id)}>Delete</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

{#if showModal}
	<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
	<div class="overlay" on:click|self={closeModal}>
		<div class="modal card" role="dialog" aria-modal="true">
			<h2>{editingId ? 'Edit Measurement' : 'Add Measurement'}</h2>

			<div class="field">
				<label for="m-date">Date</label>
				<input id="m-date" type="date" bind:value={form.date} />
			</div>

			<div class="field">
				<label for="m-time">Time</label>
				<select id="m-time" bind:value={form.time}>
					{#each times as t}
						<option value={t.value}>{t.label}</option>
					{/each}
				</select>
			</div>

			<div class="field-row">
				<div class="field">
					<label for="m-weight">Weight (lbs)</label>
					<input
						id="m-weight"
						type="number"
						min="0"
						step="0.1"
						bind:value={form.weight}
						placeholder="e.g. 8.5"
					/>
				</div>
				<div class="field">
					<label for="m-height">Height (in)</label>
					<input
						id="m-height"
						type="number"
						min="0"
						step="0.25"
						bind:value={form.height}
						placeholder="e.g. 21.5"
					/>
				</div>
			</div>

			{#if !hasValue}
				<p class="hint">Enter at least one value.</p>
			{/if}

			<div class="modal-actions">
				<button class="btn-ghost" on:click={closeModal}>Cancel</button>
				<button class="btn-primary" on:click={save} disabled={saving || !hasValue}>
					{saving ? 'Saving…' : 'Save'}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20px;
	}
	h1 {
		font-size: 22px;
	}
	.add-btn {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		font-size: 24px;
		padding: 0;
		line-height: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.error-text {
		color: var(--danger);
		font-size: 14px;
	}
	.list {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}
	.row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 12px;
		padding: 12px 16px;
	}
	.row-left {
		display: flex;
		flex-direction: column;
		gap: 5px;
	}
	.row-time {
		font-size: 14px;
		font-weight: 500;
	}
	.row-values {
		display: flex;
		gap: 6px;
	}
	.val-chip {
		font-size: 12px;
		font-weight: 600;
		padding: 2px 8px;
		border-radius: 10px;
	}
	.weight {
		background: #d1fae5;
		color: #065f46;
	}
	.height {
		background: #fce7f3;
		color: #9d174d;
	}
	.row-right {
		display: flex;
		align-items: center;
		gap: 6px;
		flex-shrink: 0;
	}
	.btn-sm {
		padding: 4px 10px;
		font-size: 13px;
	}

	/* Modal */
	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.35);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 16px;
		z-index: 100;
	}
	.modal {
		width: 100%;
		max-width: 360px;
		display: flex;
		flex-direction: column;
		gap: 14px;
		padding: 24px;
	}
	.modal h2 {
		font-size: 18px;
	}
	.field {
		display: flex;
		flex-direction: column;
		gap: 5px;
	}
	.field label {
		font-size: 13px;
		font-weight: 500;
	}
	.field-row {
		display: flex;
		gap: 12px;
	}
	.field-row .field {
		flex: 1;
	}
	.hint {
		font-size: 12px;
		color: var(--text-muted);
	}
	.modal-actions {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
		padding-top: 4px;
	}
</style>
