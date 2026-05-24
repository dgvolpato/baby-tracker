<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Diaper } from '$lib/api';
	import { formatDateTime, timeOptions, roundedTimeValue, toLocalInputDate } from '$lib/utils';

	let diapers: Diaper[] = [];
	let loading = true;
	let error = '';

	let showModal = false;
	let editingId: string | null = null;
	let saving = false;

	interface FormState {
		date: string;
		time: string;
		type: 'wet' | 'poop' | 'both';
	}

	function freshForm(): FormState {
		const now = new Date();
		return { date: toLocalInputDate(now), time: roundedTimeValue(now), type: 'wet' };
	}

	let form: FormState = freshForm();
	const times = timeOptions();

	function openAdd() {
		editingId = null;
		form = freshForm();
		showModal = true;
	}

	function openEdit(d: Diaper) {
		editingId = d.id;
		const dt = new Date(d.timestamp);
		form = {
			date: toLocalInputDate(dt),
			time: `${String(dt.getHours()).padStart(2, '0')}:${String(dt.getMinutes()).padStart(2, '0')}`,
			type: d.type
		};
		showModal = true;
	}

	function closeModal() {
		showModal = false;
	}

	async function save() {
		saving = true;
		const timestamp = new Date(`${form.date}T${form.time}:00`).toISOString();
		const data = { timestamp, type: form.type };
		try {
			if (editingId) {
				const updated = await api.diapers.update(editingId, data);
				diapers = diapers.map((d) => (d.id === editingId ? updated : d));
			} else {
				const created = await api.diapers.create(data);
				diapers = [created, ...diapers];
			}
			diapers = [...diapers].sort((a, b) => b.timestamp.localeCompare(a.timestamp));
			closeModal();
		} catch {
			// keep modal open on error
		} finally {
			saving = false;
		}
	}

	async function remove(id: string) {
		if (!confirm('Delete this diaper change?')) return;
		await api.diapers.delete(id);
		diapers = diapers.filter((d) => d.id !== id);
	}

	const typeLabel: Record<string, string> = { wet: 'Wet', poop: 'Poop', both: 'Wet + Poop' };

	onMount(async () => {
		try {
			diapers = await api.diapers.list();
		} catch {
			error = 'Failed to load diapers.';
		} finally {
			loading = false;
		}
	});
</script>

<div class="page">
	<div class="page-header">
		<h1>Diapers</h1>
		<button class="add-btn btn-primary" on:click={openAdd} aria-label="Add diaper change">+</button>
	</div>

	{#if loading}
		<p class="muted">Loading…</p>
	{:else if error}
		<p class="error-text">{error}</p>
	{:else if diapers.length === 0}
		<p class="muted">No diaper changes yet. Tap + to add one.</p>
	{:else}
		<div class="list">
			{#each diapers as d (d.id)}
				<div class="row card">
					<div class="row-left">
						<span class="row-time">{formatDateTime(d.timestamp)}</span>
						<span class="tag {d.type}">{typeLabel[d.type]}</span>
					</div>
					<div class="row-right">
						<button class="btn-ghost btn-sm" on:click={() => openEdit(d)}>Edit</button>
						<button class="btn-danger btn-sm" on:click={() => remove(d.id)}>Delete</button>
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
			<h2>{editingId ? 'Edit Diaper Change' : 'Add Diaper Change'}</h2>

			<div class="field">
				<label for="d-date">Date</label>
				<input id="d-date" type="date" bind:value={form.date} />
			</div>

			<div class="field">
				<label for="d-time">Time</label>
				<select id="d-time" bind:value={form.time}>
					{#each times as t}
						<option value={t.value}>{t.label}</option>
					{/each}
				</select>
			</div>

			<div class="field">
				<span class="field-label">Type</span>
				<div class="type-group">
					<label class="type-option" class:selected={form.type === 'wet'}>
						<input type="radio" bind:group={form.type} value="wet" />
						Wet
					</label>
					<label class="type-option" class:selected={form.type === 'poop'}>
						<input type="radio" bind:group={form.type} value="poop" />
						Poop
					</label>
					<label class="type-option" class:selected={form.type === 'both'}>
						<input type="radio" bind:group={form.type} value="both" />
						Wet + Poop
					</label>
				</div>
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
	.btn-sm {
		padding: 4px 10px;
		font-size: 13px;
	}
	.tag {
		font-size: 11px;
		padding: 2px 8px;
		border-radius: 10px;
		font-weight: 500;
	}
	.wet {
		background: #dbeafe;
		color: #1d4ed8;
	}
	.poop {
		background: #fef3c7;
		color: #92400e;
	}
	.both {
		background: #ede9fe;
		color: #6d28d9;
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
	.field label,
	.field-label {
		font-size: 13px;
		font-weight: 500;
	}
	.type-group {
		display: flex;
		gap: 8px;
	}
	.type-option {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 5px;
		font-size: 13px;
		font-weight: 500;
		padding: 8px 4px;
		border: 1px solid var(--border);
		border-radius: var(--radius);
		cursor: pointer;
		transition: background 0.15s, border-color 0.15s;
		text-align: center;
	}
	.type-option input {
		display: none;
	}
	.type-option.selected {
		border-color: var(--accent);
		background: var(--accent-light);
		color: var(--accent);
	}
	.modal-actions {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
		padding-top: 4px;
	}
</style>
