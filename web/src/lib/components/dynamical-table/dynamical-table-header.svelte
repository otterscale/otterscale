<script lang="ts">
	import type { JsonValue } from '@openfeature/server-sdk';
	import type { Column } from '@tanstack/table-core';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	let {
		ref = $bindable(null),
		column,
		children,
		fields,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		column: Column<Record<string, JsonValue>>;
		fields: Record<string, { description: string; type: string; format: string }>;
	} = $props();
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
			{#if fields[column.id].description}
				<Tooltip.Content>
					<p class="max-w-3xl truncate">{fields[column.id].description}</p>
				</Tooltip.Content>
			{/if}
		</Tooltip.Root>
	</Tooltip.Provider>
</div>
