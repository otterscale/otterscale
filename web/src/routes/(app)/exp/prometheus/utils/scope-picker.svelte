<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import type { Scope } from '$gen/api/nexus/v1/nexus_pb';

	let { selectedScope = $bindable(), scopes }: { selectedScope: Scope; scopes: Scope[] } = $props();
</script>

<span class="flex items-center gap-2">
	<p class="flex h-8 items-center rounded-lg bg-muted p-4">Scope</p>
	<Popover.Root>
		<Popover.Trigger>
			<Button variant="outline" class="h-8 w-full justify-between">
				{selectedScope.name}
			</Button>
		</Popover.Trigger>
		<Popover.Content class="w-fit p-0">
			<Command.Root>
				<Command.Input placeholder="Search" />
				<Command.List>
					<Command.Empty class="justift-center grid items-center">No Result</Command.Empty>
					<Command.Group>
						{#each scopes as scope}
							<Command.Item
								class="hover:cursor-pointer"
								value={scope.uuid}
								onSelect={() => {
									selectedScope = scope;
								}}
							>
								<div class="flex w-full items-center justify-between gap-2">
									<div class="flex items-center gap-2">
										<Icon
											icon="ph:check"
											class={cn('h-4 w-4', selectedScope === scope ? 'visible' : 'invisible')}
										/>
										{scope.name}
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
