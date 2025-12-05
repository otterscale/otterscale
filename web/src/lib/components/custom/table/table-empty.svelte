<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Empty from '$lib/components/ui/empty/index.js';
	import { cn } from '$lib/utils';

	let {
		ref = $bindable(null),
		children,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {} = $props();
</script>

<div class={cn('p-4 text-center text-muted-foreground', className)}>
	{#if children}
		{@render children()}
	{:else}
		<div class={cn('flex items-center justify-center p-4', className)}>
			<Empty.Root>
				<Empty.Header>
					<Empty.Media variant="icon">
						<Icon icon="ph:magnifying-glass" />
					</Empty.Media>
					<Empty.Title>No data</Empty.Title>
					<Empty.Description>
						Please adjust your filters or add new entries to see data here.
					</Empty.Description>
				</Empty.Header>
			</Empty.Root>
		</div>
	{/if}
</div>
