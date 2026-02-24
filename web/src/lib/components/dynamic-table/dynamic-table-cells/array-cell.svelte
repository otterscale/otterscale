<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import { type Column, type Row } from '@tanstack/table-core';
	import { stringify } from 'yaml';

	import * as CodeBlock from '$lib/components/custom/code/index.js';
	import { Button } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet/index.js';

	let {
		row,
		column
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
	} = $props();

	const data = $derived(row.original[column.id] as JsonValue[]);
</script>

<Sheet.Root>
	<Sheet.Trigger disabled={!data}>
		<Button variant="ghost" class="hover:underline">
			{data.length ?? 0}
		</Button>
	</Sheet.Trigger>
	{#if data && data.length > 0}
		<Sheet.Content
			side="right"
			class="flex h-full max-w-[50vw] min-w-[38vw] flex-col gap-0 overflow-y-auto p-4"
		>
			<Sheet.Header class="shrink-0 space-y-4">
				<Sheet.Title>{column.id}</Sheet.Title>
			</Sheet.Header>
			<CodeBlock.Root
				lang="yaml"
				hideLines
				code={stringify(data, null, 4)}
				class="m-4 rounded-lg border-none bg-muted"
			/>
		</Sheet.Content>
	{/if}
</Sheet.Root>
