<script lang="ts">
    import { Button } from '$lib/components/ui/button/index.js';
    import { cn } from '$lib/utils';
    import Icon from '@iconify/svelte';
    import * as Command from '$lib/components/ui/command/index.js';
    import * as Popover from '$lib/components/ui/popover/index.js';
    import type { Machine } from '$gen/api/machine/v1/machine_pb';

    let { selectedMachine = $bindable(), machines }: { selectedMachine: Machine; machines: Machine[] } = $props();
</script>

<span class="flex items-center gap-2">
    <p class="flex h-8 items-center rounded-lg bg-muted p-4">Machine</p>
    <Popover.Root>
        <Popover.Trigger>
            <Button variant="outline" class="h-8 w-full justify-between">
                {(!selectedMachine?.id || selectedMachine.id == '') ? 'All Machine' : selectedMachine.fqdn}
            </Button>
        </Popover.Trigger>
        <Popover.Content class="w-fit p-0">
            <Command.Root>
                <Command.Input placeholder="Search" />
                <Command.List>
                    <Command.Empty class="justify-center grid items-center">No Result</Command.Empty>
                    <Command.Group>
                        {#each machines as machine}
                            <Command.Item
                                class="hover:cursor-pointer"
                                value={machine.fqdn}
                                onSelect={() => {
                                    selectedMachine = machine;
                                }}
                            >
                                <div class="flex w-full items-center justify-between gap-2">
                                    <div class="flex items-center gap-2">
                                        <Icon
                                            icon="ph:check"
                                            class={cn('h-4 w-4', selectedMachine === machine ? 'visible' : 'invisible')}
                                        />
                                        {(!machine?.id || machine.id == '') ? 'All Machine' : machine.fqdn}
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