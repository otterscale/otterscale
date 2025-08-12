<script lang="ts" module>
	import { Button } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { FilterManager } from './utils';
</script>

<script lang="ts">
	let { filterManager }: { filterManager: FilterManager } = $props();

	const keywords = $derived(
		[...new Set(filterManager.charts.flatMap((chart) => chart.keywords))].sort()
	);
</script>

<Popover.Root>
	<Popover.Trigger>
		<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
			<Icon icon="ph:funnel" class="h-3 w-3" />
			Keyword
			<Icon icon="ph:caret-down" class="h-3 w-3" />
		</Button>
	</Popover.Trigger>
	<Popover.Content class="p-0">
		<Command.Root>
			<Command.Input placeholder="Search" />
			<Command.List>
				<Command.Empty>No keyword found.</Command.Empty>
				<Command.Group>
					{#each keywords as keyword}
						<Command.Item
							onclick={() => {
								filterManager.toggleKeyword(keyword);
							}}
						>
							<Icon
								icon="ph:check"
								class={cn(
									filterManager.isKeywordSelected(keyword) ? 'visible' : 'invisible',
									'h-4 w-4'
								)}
							/>
							{keyword}
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
