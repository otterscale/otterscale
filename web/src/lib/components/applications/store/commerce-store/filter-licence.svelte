<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { FilterManager } from './utils';

	import { Button } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { filterManager }: { filterManager: FilterManager } = $props();

	const licences = $derived([...new Set(filterManager.charts.flatMap((chart) => chart.license))].sort());
</script>

<Popover.Root>
	<Popover.Trigger>
		<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
			<Icon icon="ph:funnel" class="h-3 w-3" />
			{m.licence()}
			<Icon icon="ph:caret-down" class="h-3 w-3" />
		</Button>
	</Popover.Trigger>
	<Popover.Content class="p-0">
		<Command.Root>
			<Command.Input placeholder="Search" />
			<Command.List>
				<Command.Empty>{m.no_result()}</Command.Empty>
				<Command.Group>
					{#each licences as licence}
						{#if licence}
							<Command.Item
								onclick={() => {
									filterManager.toggleLicence(licence);
								}}
							>
								<Icon
									icon="ph:check"
									class={cn(
										filterManager.isLicenceSelected(licence) ? 'visible' : 'invisible',
										'h-4 w-4',
									)}
								/>
								{licence}
							</Command.Item>
						{/if}
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
