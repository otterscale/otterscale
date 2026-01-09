<script lang="ts">
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	import type { FieldSchema } from './types';

	let {
		ref = $bindable(null),
		children,
		field,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		field: FieldSchema;
	} = $props();
</script>

<div class={className}>
	<Tooltip.Provider>
		<Tooltip.Root>
			<Tooltip.Trigger>
				{@render children?.()}
			</Tooltip.Trigger>
			<Tooltip.Content>
				{field.description}
			</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
</div>
