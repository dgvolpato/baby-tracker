<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		goto(localStorage.getItem('token') ? '/dashboard' : '/login');
	});
</script>
