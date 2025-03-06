<script lang="ts">
	import { useId } from 'bits-ui';
	import { tick } from 'svelte';
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';

	import * as Accordion from '$lib/components/ui/accordion';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import { Input } from '$lib/components/ui/input';
	import * as Popover from '$lib/components/ui/popover';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import type { Connector } from './connector';

	let {
		parent = $bindable(),
		sources,
		destinations
	}: {
		parent: boolean;
		sources: Connector[];
		destinations: Connector[];
	} = $props();

	//#endregion

	//#region sources

	let srcOpen = $state(false);
	let srcValue = $state('');
	const srcSelected = $derived(sources?.find((s) => s.name === srcValue) ?? null);

	function closeAndFocusSourceTrigger(triggerId: string) {
		srcOpen = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}

	const srcId = useId();

	//endregion

	//#region destinations

	let dstOpen = $state(false);
	let dstValue = $state('');
	const dstSelected = $derived(destinations?.find((s) => s.name === dstValue) ?? null);

	function closeAndFocusDestinationTrigger(triggerId: string) {
		dstOpen = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}

	const dstId = useId();

	//endregion

	let confirm = $state(false);
</script>

<Popover.Root bind:open={srcOpen}>
	<Popover.Trigger
		class={buttonVariants({
			variant: 'outline',
			class: 'col-span-3 w-full'
		})}
		id={srcId}
	>
		<div class="flex items-center gap-2 text-foreground">
			{#if srcSelected}
				<span>{srcSelected.name}</span>
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
					{#each sources as source}
						<Command.Item
							value={source.name}
							onSelect={() => {
								srcValue = source.name;
								closeAndFocusSourceTrigger(srcId);
							}}
						>
							{source.name}
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>

<Popover.Root bind:open={dstOpen}>
	<Popover.Trigger
		class={buttonVariants({
			variant: 'outline',
			class: 'col-span-3 w-full'
		})}
		id={dstId}
	>
		<div class="flex items-center gap-2 text-foreground">
			{#if dstSelected}
				<span>{dstSelected.name}</span>
			{:else}
				<span> + To </span>
			{/if}
		</div>
	</Popover.Trigger>
	<Popover.Content class="p-0" align="start" side="right">
		<Command.Root>
			<Command.Input placeholder="Find..." />
			<Command.List>
				<Command.Empty>No results found.</Command.Empty>
				<Command.Group>
					{#each destinations as destination}
						<Command.Item
							value={destination.name}
							onSelect={() => {
								dstValue = destination.name;
								closeAndFocusDestinationTrigger(dstId);
							}}
						>
							{destination.name}
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>

<!-- <div class="w-full flex-col space-y-2">
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
					<div class="flex items-center gap-2 text-foreground">
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
								{#each connectors as connector}
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
				<Input class="text-foreground" type="text" id={p.key} value={p.value} />
				<p class="pt-2 text-sm text-muted-foreground">{p.description}</p>
			</fieldset>
		{/each}
		{#if selected.extraParameters}
			<Accordion.Root type="single">
				{#each selected.extraParameters as p, i}
					<Accordion.Item value="item-{i}">
						<Accordion.Trigger>{p.name}</Accordion.Trigger>
						<Accordion.Content>
							<Input class="text-foreground" type="text" id={p.key} value={p.value} />
							<p class="pt-2 text-sm text-muted-foreground">{p.description}</p>
						</Accordion.Content>
					</Accordion.Item>
				{/each}
			</Accordion.Root>
		{/if}
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
			<Popover.Root bind:open={templateOpen}>
				<Popover.Trigger
					class={buttonVariants({
						variant: 'outline',
						class: 'w-full'
					})}
					id={triggerTemplateId}
				>
					<div class="flex items-center gap-2 text-foreground">
						{#if templateSelected}
							<span>{templateSelected.name}</span>
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
												templateValue = template.name;
												selected?.parameters.forEach((p) => {
													templateSelected?.parameters.forEach((t) => {
														if (p.key === t.key) {
															p.value = t.value;
														}
													});
												});
												closeAndFocusTemplateTrigger(triggerTemplateId);
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
	<div class="flex justify-end pt-4">
		<AlertDialog.Root bind:open={confirm}>
			<AlertDialog.Trigger disabled={!selected} class={buttonVariants({})}>
				Continue
			</AlertDialog.Trigger>
			<AlertDialog.Content class="space-y-2">
				<AlertDialog.Header class="space-y-4">
					<AlertDialog.Title>Please confirm the configuration</AlertDialog.Title>
					<AlertDialog.Description class="space-y-2">
						{#if selected}
							<fieldset class="items-center gap-6 rounded-lg border p-4">
								<legend class="-ml-1 px-1 text-sm font-medium"> New </legend>
								<div class="flex items-center space-x-2 text-base text-foreground">
									<Icon icon={selected.icon} class="size-8" />
									<span>{selected.name}</span>
								</div>
							</fieldset>
							{#each selected.parameters as p}
								<fieldset class="items-center gap-6 rounded-lg border p-4">
									<legend class="-ml-1 px-1 text-sm font-medium"> {p.name} </legend>
									<p class="text-foreground">{p.value}</p>
								</fieldset>
							{/each}
						{/if}
					</AlertDialog.Description>
				</AlertDialog.Header>
				<AlertDialog.Footer>
					<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
					<AlertDialog.Action
						onclick={() => {
							confirm = false;
							parent = false;
							toast.success('Created!');
						}}
					>
						Create
					</AlertDialog.Action>
				</AlertDialog.Footer>
			</AlertDialog.Content>
		</AlertDialog.Root>
	</div>
</div> -->
