<script lang="ts">
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card/index';
	import * as Tooltip from '$lib/components/ui/tooltip/index';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { Snippet } from 'svelte';

	let {
		title,
		content,
		description,
		footer,
		hint,
		class: className
	}: {
		title: string;
		content: Snippet;
		description?: Snippet;
		footer?: Snippet;
		hint?: Snippet;
		class?: string;
	} = $props();
</script>

<Card.Root class={cn('bg-muted/40 h-full w-full truncate border-none shadow-none', className)}>
	<Card.Header class="h-[150px]">
		<Card.Title class="flex items-start space-x-2">
			<h1 class="whitespace-nowrap text-3xl">{title}</h1>
			{#if hint}
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
			{/if}
		</Card.Title>
		<Card.Description>
			{@render description?.()}
		</Card.Description>
	</Card.Header>
	<Card.Content class="h-[250px]">
		{@render content()}
	</Card.Content>
	<Card.Footer>
		{@render footer?.()}
	</Card.Footer>
</Card.Root>
