<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import type { Column } from '@tanstack/table-core';
	import { type WithElementRef } from 'bits-ui';
	import lodash from 'lodash';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	import { getColumnType } from './utils';
	import type { FieldsType, ValuesType } from '../kind-viewer/type';

	let {
		ref = $bindable(null),
		column,
		children,
		fields,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		column: Column<ValuesType>;
		fields: FieldsType;
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
