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
		items
	}: {
		parent: boolean;
		items: Connector[];
	} = $props();

	let connectors = $state<Connector[]>(items);

	//#region parameters

	let open = $state(false);
	let value = $state('');
	const selected = $derived(connectors.find((s) => s.name === value) ?? null);

	function closeAndFocusTrigger(triggerId: string) {
		open = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}

	const triggerId = useId();

	//endregion

	//#region templates

	const templates = $derived(connectors.find((s) => s.name === value)?.templates);

	let templateOpen = $state(false);
	let templateValue = $state('');
	const templateSelected = $derived(templates?.find((s) => s.name === templateValue) ?? null);

	function closeAndFocusTemplateTrigger(triggerId: string) {
		templateOpen = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}

	const triggerTemplateId = useId();

	//endregion

	let confirm = $state(false);
</script>

<div class="w-full flex-col space-y-2">
	<fieldset class="items-center rounded-lg border p-4">
		<legend class="-ml-2 px-1 text-sm font-medium"> New </legend>
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
							<span> + Connector </span>
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
										<Icon icon={connector.icon} />
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
			<fieldset class="items-center rounded-lg border p-4">
				<legend class="-ml-2 px-1 text-sm font-medium"> {p.name} </legend>
				<Input
					class="text-foreground placeholder:italic"
					type="text"
					id={p.key}
					bind:value={p.value}
					placeholder={p.description}
				/>
			</fieldset>
		{/each}
		{#if selected.extraParameters}
			<fieldset class="items-center rounded-lg border p-4">
				<legend class="-ml-2 px-1 text-sm font-medium"> Extra </legend>
				<Accordion.Root type="single">
					{#each selected.extraParameters as p, i}
						<Accordion.Item class="border-none" value="item-{i}">
							<Accordion.Trigger class="py-2">{p.name}</Accordion.Trigger>
							<Accordion.Content>
								<Input
									class="text-foreground placeholder:italic"
									type="text"
									id={p.key}
									bind:value={p.value}
									placeholder={p.description}
								/>
							</Accordion.Content>
						</Accordion.Item>
					{/each}
				</Accordion.Root>
			</fieldset>
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
		<fieldset class="grid w-full gap-6 rounded-lg border p-6">
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
							<fieldset class="items-center gap-6 rounded-lg border p-3">
								<legend class="-ml-1 px-1 text-sm font-medium"> New </legend>
								<div class="flex items-center space-x-2 text-sm text-foreground">
									<Icon icon={selected.icon} class="size-8" />
									<span>{selected.name}</span>
								</div>
							</fieldset>
							{#each selected.parameters as p}
								<fieldset class="items-center gap-6 rounded-lg border p-3">
									<legend class="-ml-1 px-1 text-sm font-medium"> {p.name} </legend>
									{#if p.value}
										<p class="text-foreground">{p.value}</p>
									{:else}
										<p class="text-muted-foreground">(empty)</p>
									{/if}
								</fieldset>
							{/each}
							{#if selected.extraParameters}
								<div class="grid grid-cols-2 gap-2">
									{#each selected.extraParameters as p}
										<fieldset class="items-center gap-6 rounded-lg border p-3">
											<legend class="-ml-1 px-1 text-sm font-medium"> {p.name} </legend>
											{#if p.value}
												<p class="text-foreground">{p.value}</p>
											{:else}
												<p class="text-muted-foreground">(empty)</p>
											{/if}
										</fieldset>
									{/each}
								</div>
							{/if}
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
</div>
