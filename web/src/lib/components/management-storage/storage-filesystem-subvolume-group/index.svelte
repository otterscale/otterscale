<script lang="ts">
	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';

	let group = $state('');
</script>

<Pickers bind:group />

<DataTable {group} />
<!-- {#if !$data}
	<PageLoading />
{:else}
	<DataTable bind:data />
{/if} -->
