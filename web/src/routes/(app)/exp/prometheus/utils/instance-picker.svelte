<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';

	let {
		selectedInstance = $bindable(),
		instances
	}: { selectedInstance: string; instances: string[] } = $props();
</script>

<span class="ml-auto flex items-center gap-2">
	<p class="flex h-10 items-center rounded-lg bg-muted p-4">Instance</p>
	<Popover.Root>
		<Popover.Trigger>
			<Button variant="outline" class="w-full justify-between">
				{selectedInstance}
			</Button>
		</Popover.Trigger>
		<Popover.Content class="w-fit p-0">
			<Command.Root>
				<Command.Input placeholder="Search" />
				<Command.List>
					<Command.Empty class="justift-center grid items-center">No Result</Command.Empty>
					<Command.Group>
						{#each instances as instance}
							<Command.Item
								class="hover:cursor-pointer"
								value={instance}
								onSelect={() => {
									selectedInstance = instance;
								}}
							>
								<div class="flex w-full items-center justify-between gap-2">
									<div class="flex items-center gap-2">
										<Icon
											icon="ph:check"
											class={cn('h-4 w-4', selectedInstance === instance ? 'visible' : 'invisible')}
										/>
										{instance}
									</div>
								</div>
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</Command.Root>
		</Popover.Content>
	</Popover.Root>
</span>
