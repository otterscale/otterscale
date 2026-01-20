<script lang="ts">
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	let {
		ref = $bindable(null),
		children,
		field,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		// eslint-disable-next-line
		field: any;
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
