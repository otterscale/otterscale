<script lang="ts" module>
	export function getColumnType(type: JsonValue, format: JsonValue) {
		if (type === 'boolean') {
			return 'boolean';
		} else if (type === 'number' || type === 'integer') {
			return 'number';
		} else if (type === 'string' && (format === 'date' || format === 'date-time')) {
			return 'time';
		} else if (type === 'string') {
			return 'string';
		} else if (type === 'array') {
			return 'array';
		} else if (type === 'object') {
			return 'object';
		} else {
			return undefined;
		}
	}
</script>

<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import type { Column } from '@tanstack/table-core';
	import { type WithElementRef } from 'bits-ui';
	import lodash from 'lodash';
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

	const type = $derived(lodash.get(fields, `${column.id}.type`));
	const format = $derived(lodash.get(fields, `${column.id}.format`));
	const columnType = $derived(getColumnType(type, format));
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
			{#if columnType}
				<Tooltip.Content>
					{columnType}
				</Tooltip.Content>
			{/if}
		</Tooltip.Root>
	</Tooltip.Provider>
</div>
