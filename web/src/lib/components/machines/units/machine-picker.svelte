<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import type { Machine } from '$lib/api/machine/v1/machine_pb';

	let { selectedMachine = $bindable(), machines }: { selectedMachine: Machine; machines: Machine[] } = $props();

	const ALL_MACHINE_ID = 'All Machine';

	function getMachineDisplayName(machine: Machine | undefined): string {
		if (!machine?.id || machine.id === ALL_MACHINE_ID) {
			return ALL_MACHINE_ID;
		}
		return machine.fqdn;
	}
</script>

<div class="flex items-center gap-2">
	<p class="bg-muted flex h-8 items-center rounded-lg p-4">Machine</p>
	<Popover.Root>
		<Popover.Trigger>
			<Button variant="outline" class="h-8 w-full justify-between">
				{getMachineDisplayName(selectedMachine)}
			</Button>
		</Popover.Trigger>
		<Popover.Content class="w-fit p-0">
			<Command.Root>
				<Command.Input placeholder="Search machines..." />
				<Command.List>
					<Command.Empty class="grid items-center justify-center">No machines found</Command.Empty>
					<Command.Group>
						{#each machines as machine}
							<Command.Item
								class="hover:cursor-pointer"
								value={machine.fqdn}
								onSelect={() => (selectedMachine = machine)}
							>
								<div class="flex w-full items-center justify-between gap-2">
									<div class="flex items-center gap-2">
										<Icon
											icon="ph:check"
											class={cn('h-4 w-4', selectedMachine === machine ? 'visible' : 'invisible')}
										/>
										{getMachineDisplayName(machine)}
									</div>
								</div>
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</Command.Root>
		</Popover.Content>
	</Popover.Root>
</div>
