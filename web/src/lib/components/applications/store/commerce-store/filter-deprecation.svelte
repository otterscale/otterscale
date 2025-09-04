<script lang="ts" module>
	import { Button } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { FilterManager } from './utils';
</script>

<script lang="ts">
	let { filterManager }: { filterManager: FilterManager } = $props();
</script>

<Popover.Root>
	<Popover.Trigger>
		<Button variant="outline" size="sm" class="flex h-8 items-center gap-2">
			<Icon icon="ph:funnel" class="h-3 w-3" />
			{m.deprecation()}
			<Icon icon="ph:caret-down" class="h-3 w-3" />
		</Button>
	</Popover.Trigger>
	<Popover.Content class="p-0">
		<Command.Root>
			<Command.List>
				<Command.Group>
					{#each [null, true, false] as deprecation}
						<Command.Item
							onclick={() => {
								filterManager.toggleDeprecation(deprecation);
							}}
						>
							<Icon
								icon="ph:check"
								class={cn(
									filterManager.isDeprecationSelected(deprecation) ? 'visible' : 'invisible',
									'h-4 w-4',
								)}
							/>
							{#if deprecation === null}
								{m.all()}
							{:else if deprecation === true}
								{m.only_deprecated()}
							{:else if deprecation === false}
								{m.only_available()}
							{/if}
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
