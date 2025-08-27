<script lang="ts" module>
	import * as Label from '$lib/components/ui/label';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { Label as LabelPrimitive } from 'bits-ui';
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		information,
		...restProps
	}: LabelPrimitive.RootProps & { information?: Snippet } = $props();
</script>

<div class="flex items-center justify-between gap-2 capitalize select-none">
	<Label.Root bind:ref data-slot="form-label" class={cn(className)} {...restProps}>
		{@render children?.()}
	</Label.Root>
	{#if information}
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger class="hover:cursor-pointer"><Icon icon="ph:info" /></Tooltip.Trigger>
				<Tooltip.Content class="select-none">{@render information()}</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	{/if}
</div>
