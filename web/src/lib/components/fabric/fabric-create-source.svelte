<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as Drawer from '$lib/components/ui/drawer';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Accordion from '$lib/components/ui/accordion';
	import { formatTimeAgo } from '$lib/formatter';
	import pb from '$lib/pb';
	import * as Table from '$lib/components/ui/table';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';

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

	interface Parameter {
		key: string;
		name?: string;
		value: string;
		description?: string;
	}

	interface Template {
		name: string;
		parameters: Parameter[];
	}

	let {
		item
	}: {
		item: {
			steps: boolean[];
			connectors: {
				name: string;
				icon: string;
				parameters: Parameter[];
				templates?: Template[];
			}[];
		};
	} = $props();

	let open = $state(false);
	let value = $state('');
	const selected = $derived(item.connectors.find((s) => s.name === value) ?? null);

	function closeAndFocusTrigger(triggerId: string) {
		open = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}

	const triggerId = useId();

	// TODO: BETTER
	const templates = $derived(item.connectors.find((s) => s.name === value)?.templates);

	const triggerId1 = useId();
	let open1 = $state(false);
	let value1 = $state('');
	const selected1 = $derived(
		item.connectors.find((s) => s.name === value)?.templates?.find((s) => s.name === value1) ?? null
	);

	function closeAndFocusTrigger1(triggerId: string) {
		open1 = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}
</script>

<div class="w-full flex-col space-y-2">
	<fieldset class="items-center gap-6 rounded-lg border p-4">
		<legend class="-ml-1 px-1 text-sm font-medium"> New </legend>
		<div class="flex space-x-4">
			<div class="flex-col">
				{#if selected}
					<Icon icon={selected.icon} class="size-8" />
				{:else}
					<Skeleton class="size-8" />
				{/if}
			</div>
			<Popover.Root bind:open>
				<Popover.Trigger
					class={buttonVariants({
						variant: 'outline',
						class: 'col-span-3 w-full'
					})}
					id={triggerId}
				>
					<div class="flex items-center gap-2">
						{#if selected}
							<span>{selected.name}</span>
						{:else}
							<span> + From </span>
						{/if}
					</div>
				</Popover.Trigger>
				<Popover.Content class="p-0" align="start" side="right">
					<Command.Root>
						<Command.Input placeholder="Find..." />
						<Command.List>
							<Command.Empty>No results found.</Command.Empty>
							<Command.Group>
								{#each item.connectors as connector}
									<Command.Item
										value={connector.name}
										onSelect={() => {
											value = connector.name;
											closeAndFocusTrigger(triggerId);
										}}
									>
										{connector.name}
									</Command.Item>
								{/each}
							</Command.Group>
						</Command.List>
					</Command.Root>
				</Popover.Content>
			</Popover.Root>
		</div>
	</fieldset>
	{#if selected}
		{#each selected.parameters as p}
			<fieldset class="items-center gap-6 rounded-lg border p-4">
				<legend class="-ml-1 px-1 text-sm font-medium"> {p.name} </legend>
				<Input type="text" id={p.key} value={p.value} />
				<p class="pt-2 text-sm text-muted-foreground">{p.description}</p>
			</fieldset>
		{/each}
	{/if}
	{#if templates && templates.length > 0}
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
			<Popover.Root bind:open={open1}>
				<Popover.Trigger
					class={buttonVariants({
						variant: 'outline',
						class: 'w-full'
					})}
					id={triggerId1}
				>
					<div class="flex items-center gap-2">
						{#if selected1}
							<span>{selected1.name}</span>
						{:else}
							<span>+ Select</span>
						{/if}
					</div>
				</Popover.Trigger>
				<Popover.Content class="p-0" align="start" side="right">
					<Command.Root>
						<Command.Input placeholder="Find..." />
						<Command.List>
							<Command.Empty>No results found.</Command.Empty>
							<Command.Group>
								{#if templates}
									{#each templates as template}
										<Command.Item
											value={template.name}
											onSelect={() => {
												value1 = template.name;
												closeAndFocusTrigger1(triggerId1);
											}}
										>
											{template.name}
										</Command.Item>
									{/each}
								{/if}
							</Command.Group>
						</Command.List>
					</Command.Root>
				</Popover.Content>
			</Popover.Root>
		</fieldset>
	{/if}
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
