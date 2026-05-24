<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Feeding } from '$lib/api';
	import { formatDateTime, timeOptions, roundedTimeValue, toLocalInputDate } from '$lib/utils';

	let feedings: Feeding[] = [];
	let loading = true;
	let error = '';

	// Modal state
	let showModal = false;
	let editingId: string | null = null;
	let saving = false;

	interface FormState {
		date: string;
		time: string;
		type: 'formula' | 'breast';
		oz: string;
	}

	function freshForm(): FormState {
		const now = new Date();
		return {
			date: toLocalInputDate(now),
			time: roundedTimeValue(now),
			type: 'formula',
			oz: ''
		};
	}

	let form: FormState = freshForm();
	const times = timeOptions();

	function openAdd() {
		editingId = null;
		form = freshForm();
		showModal = true;
	}

	function openEdit(f: Feeding) {
		editingId = f.id;
		const d = new Date(f.timestamp);
		form = {
			date: toLocalInputDate(d),
			time: `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`,
			type: f.type,
			oz: String(f.oz)
		};
		showModal = true;
	}

	function closeModal() {
		showModal = false;
	}

	async function save() {
		const ozNum = parseFloat(form.oz);
		if (!form.oz || isNaN(ozNum) || ozNum <= 0) return;
		saving = true;
		const timestamp = new Date(`${form.date}T${form.time}:00`).toISOString();
		const data = { timestamp, type: form.type, oz: ozNum };
		try {
			if (editingId) {
				const updated = await api.feedings.update(editingId, data);
				feedings = feedings.map((f) => (f.id === editingId ? updated : f));
			} else {
				const created = await api.feedings.create(data);
				feedings = [created, ...feedings];
			}
			feedings = [...feedings].sort((a, b) => b.timestamp.localeCompare(a.timestamp));
			closeModal();
		} catch {
			// keep modal open on error
		} finally {
			saving = false;
		}
	}

	async function remove(id: string) {
		if (!confirm('Delete this feeding?')) return;
		await api.feedings.delete(id);
		feedings = feedings.filter((f) => f.id !== id);
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
	<div class="page-header">
		<h1>History</h1>
		<button class="add-btn btn-primary" on:click={openAdd} aria-label="Add feeding">+</button>
	</div>

	{#if loading}
		<p class="muted">Loading…</p>
	{:else if error}
		<p class="error-text">{error}</p>
	{:else if feedings.length === 0}
		<p class="muted">No feedings yet. Tap + to add one.</p>
	{:else}
		<div class="list">
			{#each feedings as f (f.id)}
				<div class="row card">
					<div class="row-left">
						<span class="row-time">{formatDateTime(f.timestamp)}</span>
						<span class="tag {f.type}">{f.type === 'formula' ? 'Formula' : 'Breast milk'}</span>
					</div>
					<div class="row-right">
						<span class="row-oz">{f.oz} oz</span>
						<button class="btn-ghost btn-sm" on:click={() => openEdit(f)}>Edit</button>
						<button class="btn-danger btn-sm" on:click={() => remove(f.id)}>Delete</button>
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
			<h2>{editingId ? 'Edit Feeding' : 'Add Feeding'}</h2>

			<div class="field">
				<label for="f-date">Date</label>
				<input id="f-date" type="date" bind:value={form.date} />
			</div>

			<div class="field">
				<label for="f-time">Time</label>
				<select id="f-time" bind:value={form.time}>
					{#each times as t}
						<option value={t.value}>{t.label}</option>
					{/each}
				</select>
			</div>

			<div class="field">
				<span class="field-label">Type</span>
				<div class="radio-group">
					<label class="radio">
						<input type="radio" bind:group={form.type} value="formula" />
						Formula
					</label>
					<label class="radio">
						<input type="radio" bind:group={form.type} value="breast" />
						Breast milk
					</label>
				</div>
			</div>

			<div class="field">
				<label for="f-oz">Amount (oz)</label>
				<input id="f-oz" type="number" min="0.5" step="0.5" bind:value={form.oz} placeholder="e.g. 4" />
			</div>

			<div class="modal-actions">
				<button class="btn-ghost" on:click={closeModal}>Cancel</button>
				<button class="btn-primary" on:click={save} disabled={saving}>
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
		gap: 4px;
	}
	.row-time {
		font-size: 14px;
		font-weight: 500;
	}
	.row-right {
		display: flex;
		align-items: center;
		gap: 6px;
		flex-shrink: 0;
	}
	.row-oz {
		font-size: 16px;
		font-weight: 600;
		color: var(--accent);
		margin-right: 4px;
	}
	.btn-sm {
		padding: 4px 10px;
		font-size: 13px;
	}
	.tag {
		font-size: 11px;
		padding: 2px 7px;
		border-radius: 10px;
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
		max-width: 380px;
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
	.field label,
	.field-label {
		font-size: 13px;
		font-weight: 500;
	}
	.radio-group {
		display: flex;
		gap: 20px;
		padding: 4px 0;
	}
	.radio {
		display: flex;
		align-items: center;
		gap: 6px;
		font-size: 14px;
		cursor: pointer;
	}
	.radio input {
		width: auto;
	}
	.modal-actions {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
		padding-top: 4px;
	}
</style>
