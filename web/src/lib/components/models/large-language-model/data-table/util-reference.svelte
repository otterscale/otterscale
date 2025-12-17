<script lang="ts" module>
	import { onMount } from 'svelte';

	import * as Code from '$lib/components/custom/code';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Item from '$lib/components/ui/item/index.js';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tabs from '$lib/components/ui/tabs';
	import { m } from '$lib/paraglide/messages.js';
</script>

<script lang="ts">
	let {
		modelName
	}: {
		modelName: string;
	} = $props();

	let configuration = $state('');
	async function fetchConfiguration() {
		const response = await fetch(`https://huggingface.co/${modelName}/resolve/main/config.json`);
		configuration = await response.text();
	}

	let information = $state('');
	async function fetchInformation() {
		const response = await fetch(`https://huggingface.co/api/models/${modelName}`);
		information = await response.text();
	}

	async function fetchAll() {
		await Promise.all([fetchConfiguration(), fetchInformation()]);
	}

	let isLoaded = $state(false);
	onMount(async () => {
		await fetchAll();
		isLoaded = true;
	});
</script>

<Sheet.Root>
	<Sheet.Trigger class={buttonVariants({ variant: 'outline' })}>
		{m.reference()}
	</Sheet.Trigger>
	<Sheet.Content class="min-w-[38vw]">
		{#if isLoaded}
			<Item.Root class="w-full">
				<Item.Content class="flex flex-col items-start">
					<Item.Title class="text-xl font-bold">
						{modelName}
					</Item.Title>
				</Item.Content>
			</Item.Root>
			<Tabs.Root value="configuration" class="h-full p-4">
				<Tabs.List>
					<Tabs.Trigger value="configuration">{m.configuration()}</Tabs.Trigger>
					<Tabs.Trigger value="information">{m.information()}</Tabs.Trigger>
				</Tabs.List>
				<Tabs.Content value="configuration">
					<Code.Root
						lang="json"
						code={JSON.stringify(JSON.parse(configuration), null, 2)}
						class="h-fit max-h-[77vh] w-full border-none"
					/>
				</Tabs.Content>
				<Tabs.Content value="information">
					<Code.Root
						lang="json"
						code={JSON.stringify(JSON.parse(information), null, 2)}
						class="h-fit max-h-[77vh] w-full border-none"
					/>
				</Tabs.Content>
			</Tabs.Root>
		{/if}
	</Sheet.Content>
</Sheet.Root>
