<script lang="ts" module>
	import type { JsonValue } from '@bufbuild/protobuf';

	export type LinkMetadata = {
		hyperlink: string;
	};
</script>

<script lang="ts">
	import { type Column, type Row } from '@tanstack/table-core';
	import { onMount } from 'svelte';

	let {
		row,
		column,
		metadata
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
		metadata: LinkMetadata;
	} = $props();

	const data: JsonValue = $derived(row.original[column.id]);

	onMount(() => {
		if (metadata === undefined) {
			console.warn(`Expected metadata of ${column.id} for LinkCell, but got metadata:`, metadata);
		}
	});
</script>

<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
<a href={metadata.hyperlink} class="hover:underline">
	<p class="max-w-3xs truncate">{data}</p>
</a>
