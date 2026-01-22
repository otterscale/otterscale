<script lang="ts">
	import Icon from '@iconify/svelte';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import SearchAlert from '@lucide/svelte/icons/search-alert';

	import { Button } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Item from '$lib/components/ui/item/index.js';
	import * as Popover from '$lib/components/ui/popover';

	let {
		value = $bindable(),
		resource,
		options,
		onSelect,
		class: className
	}: {
		value?: string;
		resource: string;
		options: {
			icon: string;
			label: string;
			value: string;
			description: string;
		}[];
		onSelect?: (option: {
			icon: string;
			label: string;
			value: string;
			description: string;
		}) => void;
		class?: string;
	} = $props();

	let open = $state(false);

	let searchTerm = $state('');

	const selectedOption = $derived(options.find((option) => option.value === value));

	function handleSelect(currentValue: string) {
		value = currentValue;
		open = false;
	}
	function handleReset() {
		if (searchTerm) searchTerm = '';
	}
</script>

<Popover.Root bind:open>
	<Popover.Trigger class={className}>
		{#snippet child({ props })}
			<Button
				variant="outline"
				role="combobox"
				aria-expanded={open}
				class="w-full justify-between bg-background px-3 font-normal outline-offset-0 hover:bg-background focus-visible:border-ring focus-visible:outline-[3px] focus-visible:outline-ring/20"
				{...props}
			>
				{#if value && selectedOption}
					<span class="flex min-w-0 items-center gap-2">
						<Icon icon={selectedOption.icon} />
						<span class="truncate">{selectedOption.label}</span>
					</span>
				{:else}
					<span class="text-muted-foreground">{resource}</span>
				{/if}
				<ChevronDown size={16} class="shrink-0 text-muted-foreground/80" aria-hidden="true" />
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content class="w-full min-w-(--bits-popover-anchor-width) p-0" align="start">
		<Command.Root>
			<Command.Input placeholder="Search" bind:value={searchTerm} />
			<Command.List>
				<Command.Empty>
					<Empty.Root>
						<Empty.Header>
							<Empty.Media variant="icon">
								<SearchAlert />
							</Empty.Media>
							<Empty.Title>No Selection</Empty.Title>
							<Empty.Description>
								No items found. Clear the filter to see all options.
								<br />
								Try broadening your search terms or check your spelling.
							</Empty.Description>
						</Empty.Header>
						<Empty.Content>
							<Button size="sm" onclick={handleReset}>Reset</Button>
						</Empty.Content>
					</Empty.Root>
				</Command.Empty>
				<Command.Group>
					{#each options as option (option.value)}
						<Command.Item
							value={option.value}
							onSelect={() => {
								handleSelect(option.value);
								onSelect?.(option);
							}}
						>
							<Item.Root size="sm" class="w-full p-0">
								<Item.Media class="p-1">
									<Icon icon={option.icon} class="size-5" />
								</Item.Media>
								<Item.Content class="gap-0.5">
									<Item.Title class="text-xs">{option.label}</Item.Title>
									<Item.Description class="text-xs">{option.description}</Item.Description>
								</Item.Content>
							</Item.Root>
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
