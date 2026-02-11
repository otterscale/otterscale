<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import type { Column } from '@tanstack/table-core';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	import type { DataSchemaType } from './utils';

	let {
		ref = $bindable(null),
		column,
		dataSchemas,
		children,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		column: Column<Record<string, JsonValue>>;
		dataSchemas: Record<string, DataSchemaType>;
	} = $props();

	const dataSchema = $derived(dataSchemas[column.id]);
</script>

<div class={className}>
	<Tooltip.Provider>
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#if children}
					{@render children()}
				{:else}
					<h3>{column.id}</h3>
				{/if}
			</Tooltip.Trigger>
			{#if dataSchema}
				<Tooltip.Content>
					{dataSchema}
				</Tooltip.Content>
			{/if}
		</Tooltip.Root>
	</Tooltip.Provider>
</div>
