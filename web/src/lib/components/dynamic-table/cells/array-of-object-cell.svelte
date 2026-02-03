<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import Braces from '@lucide/svelte/icons/braces';
	import File from '@lucide/svelte/icons/file';
	import type { Column, Row } from '@tanstack/table-core';
	import lodash from 'lodash';

	import * as Code from '$lib/components/custom/code/index.js';
	import { Button } from '$lib/components/ui/button';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Item from '$lib/components/ui/item';
	import * as Sheet from '$lib/components/ui/sheet/index.js';

	let {
		keys,
		row,
		column,
		fields
	}: {
		keys: {
			title: string;
			description: string;
			actions: string;
		};
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
		fields: Record<string, { description: string; type: string; format: string }>;
	} = $props();

	const data = $derived(row.original[column.id] as JsonValue[]);
</script>

<Sheet.Root>
	<Sheet.Trigger>
		<Button variant="ghost" class="hover:underline">
			{data.length}
		</Button>
	</Sheet.Trigger>
	<Sheet.Content
		side="right"
		class="flex h-full max-w-[50vw] min-w-[38vw] flex-col gap-0 overflow-y-auto p-4"
	>
		<Sheet.Header class="shrink-0 space-y-4">
			<Sheet.Title>{column.id}</Sheet.Title>
			<Sheet.Description>
				{fields[column.id].description}
			</Sheet.Description>
		</Sheet.Header>
		{#if data.length > 0}
			<div class="space-y-0">
				{#each data as datum, index (index)}
					{#if datum}
						<Collapsible.Root class="rounded-lg transition-colors duration-200 hover:bg-muted/50">
							<Collapsible.Trigger class="w-full transition-colors duration-200 hover:underline">
								<Item.Root size="sm">
									<Item.Media variant="icon">
										<File />
									</Item.Media>
									<Item.Content class="min-w-0 flex-1 text-left">
										<Item.Title class="w-full">
											{lodash.get(datum, [keys.title])}
										</Item.Title>
										<Item.Description class="wrap-break-words breaks-all">
											{lodash.get(datum, [keys.description])}
										</Item.Description>
									</Item.Content>
									<Item.Actions>
										{lodash.get(datum, [keys.actions])}
									</Item.Actions>
								</Item.Root>
							</Collapsible.Trigger>
							<Collapsible.Content class="overflow-hidden transition-all duration-300 ease-in-out">
								<Code.Root
									lang="json"
									hideLines
									code={JSON.stringify(datum, null, 4)}
									class="border-none bg-transparent px-8"
								/>
							</Collapsible.Content>
						</Collapsible.Root>
					{/if}
				{/each}
			</div>
		{:else}
			<Empty.Root class="m-4 bg-muted/50">
				<Empty.Header>
					<Empty.Media variant="icon">
						<Braces size={36} />
					</Empty.Media>
					<Empty.Title>No Data</Empty.Title>
					<Empty.Description>
						No data is currently available for this resource.
						<br />
						To populate this resource, please add properties or values through the resource editor.
					</Empty.Description>
				</Empty.Header>
				<Empty.Content></Empty.Content>
			</Empty.Root>
		{/if}
	</Sheet.Content>
</Sheet.Root>
