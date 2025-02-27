<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Drawer from '$lib/components/ui/drawer';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Accordion from '$lib/components/ui/accordion';
	import { formatTimeAgo } from '$lib/formatter';
	import pb from '$lib/pb';
	import * as Table from '$lib/components/ui/table';
	import * as Dialog from '$lib/components/ui/dialog';

	import { useId } from 'bits-ui';

	import CircleAlert from 'lucide-svelte/icons/circle-alert';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import * as Carousel from '$lib/components/ui/carousel';
	import * as Card from '$lib/components/ui/card';

	import * as Popover from '$lib/components/ui/popover/index.js';
	import Ellipsis from 'lucide-svelte/icons/ellipsis';
	import { tick } from 'svelte';
	import * as Command from '$lib/components/ui/command';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Label } from '$lib/components/ui/label/index.js';
	// let { items }: { items: { name: string; icon: string; active: boolean }[] } = $props();
	let open = $state(false);

	type Status = {
		value: string;
		label: string;
	};

	const statuses: Status[] = [
		{
			value: 'backlog',
			label: 'Backlog'
		},
		{
			value: 'todo',
			label: 'Todo'
		},
		{
			value: 'in progress',
			label: 'In Progress'
		},
		{
			value: 'done',
			label: 'Done'
		},
		{
			value: 'canceled',
			label: 'Canceled'
		}
	];

	let value = $state('');

	const selectedStatus = $derived(statuses.find((s) => s.value === value) ?? null);

	// We want to refocus the trigger button when the user selects
	// an item from the list so users can continue navigating the
	// rest of the form with the keyboard.
	function closeAndFocusTrigger(triggerId: string) {
		open = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}

	const triggerId = useId();
</script>

<div class="w-full flex-col space-y-2">
	<fieldset class="grid w-full grid-cols-7 items-center gap-6 rounded-lg border p-4">
		<legend class="-ml-1 px-1 text-sm font-medium"> New </legend>
		<Popover.Root>
			<Popover.Trigger
				class={buttonVariants({
					variant: 'outline',
					class: 'col-span-3 w-full'
				})}
				id={triggerId}
			>
				<div class="flex items-center gap-2">
					{#if selectedStatus}
						<span>{selectedStatus.label}</span>
					{:else}
						<span> + From </span>
					{/if}
				</div>
			</Popover.Trigger>
			<Popover.Content class="p-0" align="start" side="right">
				<Command.Root>
					<Command.Input placeholder="Change status..." />
					<Command.List>
						<Command.Empty>No results found.</Command.Empty>
						<Command.Group>
							{#each statuses as status}
								<Command.Item
									value={status.value}
									onSelect={() => {
										value = status.value;
										closeAndFocusTrigger(triggerId);
									}}
								>
									{status.label}
								</Command.Item>
							{/each}
						</Command.Group>
					</Command.List>
				</Command.Root>
			</Popover.Content>
		</Popover.Root>
		<div class="col-span-1 flex justify-center">
			<Icon icon="line-md:chevron-small-triple-right" class="size-8" />
		</div>
		<Popover.Root>
			<Popover.Trigger
				class={buttonVariants({
					variant: 'outline',
					class: 'col-span-3 w-full'
				})}
				id={triggerId}
			>
				<div class="flex items-center gap-2">
					{#if selectedStatus}
						<span>{selectedStatus.label}</span>
					{:else}
						<span> + To </span>
					{/if}
				</div>
			</Popover.Trigger>
			<Popover.Content class="p-0" align="start" side="right">
				<Command.Root>
					<Command.Input placeholder="Change status..." />
					<Command.List>
						<Command.Empty>No results found.</Command.Empty>
						<Command.Group>
							{#each statuses as status}
								<Command.Item
									value={status.value}
									onSelect={() => {
										value = status.value;
										closeAndFocusTrigger(triggerId);
									}}
								>
									{status.label}
								</Command.Item>
							{/each}
						</Command.Group>
					</Command.List>
				</Command.Root>
			</Popover.Content>
		</Popover.Root>
	</fieldset>
	<div class="relative">
		<div class="absolute inset-0 flex items-center">
			<span class="w-full border-t"></span>
		</div>
		<div class="relative flex justify-center text-xs uppercase">
			<span class="bg-background px-2 text-muted-foreground"> Or continue with </span>
		</div>
	</div>
	<fieldset class="grid w-full gap-6 rounded-lg border p-4">
		<legend class="-ml-1 px-1 text-sm font-medium"> Template </legend>
		<Popover.Root>
			<Popover.Trigger
				class={buttonVariants({
					variant: 'outline',
					class: 'w-full'
				})}
				id={triggerId}
			>
				<div class="flex items-center gap-2">
					{#if selectedStatus}
						<span>{selectedStatus.label}</span>
					{:else}
						<span>+ Select</span>
					{/if}
				</div>
			</Popover.Trigger>
			<Popover.Content class="p-0" align="start" side="right">
				<Command.Root>
					<Command.Input placeholder="Change status..." />
					<Command.List>
						<Command.Empty>No results found.</Command.Empty>
						<Command.Group>
							{#each statuses as status}
								<Command.Item
									value={status.value}
									onSelect={() => {
										value = status.value;
										closeAndFocusTrigger(triggerId);
									}}
								>
									{status.label}
								</Command.Item>
							{/each}
						</Command.Group>
					</Command.List>
				</Command.Root>
			</Popover.Content>
		</Popover.Root>
	</fieldset>
	<div class="flex justify-around pt-4">
		<Button
			size="lg"
			variant="outline"
			onclick={() => {
				// connectorOpens = [false, false, false, false];
			}}>Back</Button
		>
		<Button
			size="lg"
			onclick={() => {
				// connectorOpens = [false, true, false, false];
			}}>Next</Button
		>
	</div>
</div>
