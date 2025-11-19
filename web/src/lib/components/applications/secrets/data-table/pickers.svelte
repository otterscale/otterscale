<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';

	import { ApplicationService, type Namespace } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select/index.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { selectedNamespace = $bindable(), scope }: { selectedNamespace: string; scope: string } =
		$props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	const namespaces: Writable<Namespace[]> = writable([]);
	async function fetch() {
		applicationClient
			.listNamespaces({ scope })
			.then((response) => {
				namespaces.set(response.namespaces);
			})
			.catch((error) => {
				console.debug('Failed to fetch namespaces:', error);
			});
	}

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
</script>

{#if isLoaded}
	<Select.Root type="single" bind:value={selectedNamespace}>
		<Select.Trigger
			class={cn('border-none shadow-none', buttonVariants({ variant: 'ghost', size: 'default' }))}
		>
			<span class="flex items-center gap-1">
				<Icon icon="ph:cube" />
				{selectedNamespace}
			</span>
		</Select.Trigger>
		<Select.Content class="rounded-xl">
			{#each $namespaces as namespace (namespace.name)}
				<Select.Item value={namespace.name} class="rounded-lg">{namespace.name}</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
{:else}
	<Loading.Selection />
{/if}
