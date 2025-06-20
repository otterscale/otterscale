<script lang="ts">
	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';

	let group = $state('');
	let volume = $state('');
</script>

<Pickers bind:group bind:volume />

<DataTable {group} {volume} />
