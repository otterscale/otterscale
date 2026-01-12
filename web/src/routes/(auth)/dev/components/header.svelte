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
		field: any;
	} = $props();
</script>

<div class={className}>
	<Tooltip.Provider>
		<Tooltip.Root>
			<Tooltip.Trigger>
				{@render children?.()}
			</Tooltip.Trigger>
			<Tooltip.Content class="max-w-xl p-2">
				{field.description}
			</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
</div>
