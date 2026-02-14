<script lang="ts">
	import Icon from '@iconify/svelte';

	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Item from '$lib/components/ui/item/index.js';

	import type { Item as ItemType } from './types';

	let { items, limit = 1 }: { items: ItemType[]; limit?: number } = $props();
</script>

{#if limit >= 0}
	<HoverCard.Root>
		<HoverCard.Trigger>
			<div class="flex flex-wrap items-center gap-1">
				{#each items.slice(0, limit) as item, index (index)}
					<Badge variant="outline" class="flex items-center gap-1">
						{#if item.icon}
							<Icon icon={item.icon} class="size-5" />
						{/if}
						<p class="text-xs">{item.title}</p>
					</Badge>
					{#if items.length > limit}
						<p class="cursor-default text-xs font-light text-muted-foreground">
							+ {items.length - limit}
						</p>
					{/if}
				{/each}
			</div>
		</HoverCard.Trigger>
		<HoverCard.Content class="max-h-xl w-fit max-w-sm overflow-y-auto">
			<div class="flex flex-wrap items-center gap-2">
				{#each items as item, index (index)}
					{#if item.description}
						<Item.Root size="sm" class="p-0">
							{#if item.icon}
								<Item.Media variant="icon">
									<Icon icon={item.icon} class="size-3" />
								</Item.Media>
							{/if}
							<Item.Content>
								<Item.Title class="text-xs">
									{item.title}
								</Item.Title>
								<Item.Description class="text-xs">
									{item.description}
								</Item.Description>
							</Item.Content>
							<Item.Actions></Item.Actions>
						</Item.Root>
					{:else}
						<span class="flex items-center gap-1 py-1">
							{#if item.icon}
								<Icon icon={item.icon} />
							{/if}
							<p class="text-xs">{item.title}</p>
						</span>
					{/if}
				{/each}
			</div>
		</HoverCard.Content>
	</HoverCard.Root>
{/if}
