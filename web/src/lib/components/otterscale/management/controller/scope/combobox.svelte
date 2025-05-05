<script lang="ts">
	import Icon from '@iconify/svelte';
	import { tick } from 'svelte';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';
	import type { Scope } from '$gen/api/nexus/v1/nexus_pb';

	let {
		scopes,
		selected = $bindable(),
		onSelect
	}: { scopes: Scope[]; selected: string; onSelect: () => {} } = $props();

	let open = $state(false);
	let triggerRef = $state<HTMLButtonElement>(null!);

	const selectedValue = $derived(scopes.find((f) => f.uuid === selected)?.name ?? 'Select a scope');

	// We want to refocus the trigger button when the user selects
	// an item from the list so users can continue navigating the
	// rest of the form with the keyboard.
	function closeAndFocusTrigger() {
		open = false;
		tick().then(() => {
			triggerRef.focus();
		});
	}
</script>

<Popover.Root bind:open>
	<Popover.Trigger bind:ref={triggerRef}>
		{#snippet child({ props })}
			<Button
				variant="outline"
				class="w-[200px] justify-between"
				{...props}
				role="combobox"
				aria-expanded={open}
			>
				{selectedValue || 'Select a scope'}
				<Icon icon="ph:caret-up-down" class="opacity-50" />
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content class="w-[200px] p-0">
		<Command.Root>
			<Command.Input placeholder="Search scope..." class="h-9" />
			<Command.List>
				<Command.Empty>No scope found.</Command.Empty>
				<Command.Group value="scopes">
					{#each scopes as scope (scope.uuid)}
						<Command.Item
							value={scope.uuid}
							onSelect={() => {
								selected = scope.uuid;
								closeAndFocusTrigger();
								onSelect();
							}}
						>
							<Icon icon="ph:check" class={cn(selected !== scope.uuid && 'text-transparent')} />
							{scope.name}
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
