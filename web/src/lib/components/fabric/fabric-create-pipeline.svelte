<script lang="ts">
	import { useId } from 'bits-ui';
	import { tick } from 'svelte';
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';

	import { buttonVariants } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import type { pbConnector } from '$lib/pb';
	import PipelineConnector from './pipeline_connector.svelte';
	import { connectorIcon } from '$lib/connector';

	let {
		parent = $bindable(),
		sources,
		destinations,
		source = $bindable(),
		destination = $bindable()
	}: {
		parent: boolean;
		sources?: pbConnector[];
		destinations?: pbConnector[];
		source?: pbConnector;
		destination?: pbConnector;
	} = $props();

	//#endregion

	//#region sources

	let srcOpen = $state(false);
	let srcValue = $state('');
	const srcSelected = $derived(source ?? sources?.find((s) => s.name === srcValue) ?? null);

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
	const dstSelected = $derived(
		destination ?? destinations?.find((s) => s.name === dstValue) ?? null
	);

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

<div class="w-full flex-col space-y-4">
	{#if sources && destinations}
		<fieldset class="flex items-center rounded-lg border p-4">
			<legend class="-ml-2 px-1 text-sm font-medium"> New </legend>
			<Popover.Root bind:open={srcOpen}>
				<Popover.Trigger
					class={buttonVariants({ variant: 'outline', class: 'flex w-full flex-auto' })}
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
										<Icon icon={connectorIcon(source.type)} />
										{source.name}
									</Command.Item>
								{/each}
							</Command.Group>
						</Command.List>
					</Command.Root>
				</Popover.Content>
			</Popover.Root>
			<Icon icon="line-md:chevron-small-triple-right" class="size-10 flex-none animate-pulse" />
			<Popover.Root bind:open={dstOpen}>
				<Popover.Trigger
					class={buttonVariants({ variant: 'outline', class: 'flex w-full flex-auto' })}
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
										<Icon icon={connectorIcon(destination.type)} />
										{destination.name}
									</Command.Item>
								{/each}
							</Command.Group>
						</Command.List>
					</Command.Root>
				</Popover.Content>
			</Popover.Root>
		</fieldset>
	{/if}

	<div class="flex w-full space-x-4">
		<fieldset class="flex w-full items-center rounded-lg border p-4">
			<legend class="-ml-2 px-1 text-sm font-medium"> Source </legend>
			<PipelineConnector selected={srcSelected} />
		</fieldset>

		<fieldset class="flex w-full items-center rounded-lg border p-4">
			<legend class="-ml-2 px-1 text-sm font-medium"> Destination </legend>
			<PipelineConnector selected={dstSelected} />
		</fieldset>
	</div>

	<div class="flex justify-end pt-4">
		<AlertDialog.Root bind:open={confirm}>
			<AlertDialog.Trigger disabled={!srcSelected || !dstSelected} class={buttonVariants({})}>
				Continue
			</AlertDialog.Trigger>
			<AlertDialog.Content class="space-y-2">
				<AlertDialog.Header class="space-y-4">
					<AlertDialog.Title>Confirm Pipeline Creation</AlertDialog.Title>
					<AlertDialog.Description class="space-y-2">
						Are you sure you want to create a data pipeline connecting "{srcSelected?.name}" to "{dstSelected?.name}"?
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
