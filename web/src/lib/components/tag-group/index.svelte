<script lang="ts">
	import Icon from '@iconify/svelte';

	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { cn } from '$lib/utils';

	import type { Item } from './types';

	let { items, limit = 1 }: { items: Item[]; limit?: number } = $props();
</script>

{#if limit >= 0}
	{#if items.length > limit}
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					<div class={cn('flex flex-wrap items-center gap-1')}>
						{#each items.slice(0, limit) as item, index (index)}
							<Badge variant="outline" class="flex items-center gap-1">
								{#if item.icon}
									<Icon icon={item.icon} class="size-5" />
								{/if}
								<p class="text-xs">{item.title}</p>
							</Badge>
							<p class="text-xs font-light text-muted-foreground">
								+ {items.length - limit}
							</p>
						{/each}
					</div>
				</Tooltip.Trigger>
				<Tooltip.Content class="max-h-xl h-fit w-fit max-w-xs overflow-y-auto">
					<div class="my-2 flex flex-wrap gap-2">
						{#each items as item, index (index)}
							<Badge variant="outline" class="flex items-center gap-1 overflow-auto">
								{#if item.icon}
									<Icon icon={item.icon} class="size-5 text-card" />
								{/if}
								<p class="text-xs text-card">{item.title}</p>
							</Badge>
						{/each}
					</div>
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	{/if}
{/if}
