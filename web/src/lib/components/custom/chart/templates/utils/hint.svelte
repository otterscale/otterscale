<script lang="ts">
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip/index';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { Snippet } from 'svelte';

	let {
		hint
	}: {
		hint: Snippet;
	} = $props();
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger class={cn(buttonVariants({ variant: 'ghost', size: 'icon' }))}>
			<Icon icon="ph:info" />
		</Tooltip.Trigger>
		<Tooltip.Content class="bg-popover text-popover-foreground border p-2 shadow">
			{@render hint()}
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
