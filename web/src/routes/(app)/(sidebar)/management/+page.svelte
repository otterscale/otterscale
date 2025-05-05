<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Nexus, type Error, type Scope } from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Monitor } from '$lib/components/otterscale/index';
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Label } from '$lib/components/ui/label';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.listScopes({});
			scopesStore.set(response.scopes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}
	async function createDefaultScope() {
		try {
			await client.createDefaultScope({});
		} catch (error) {
			console.error('Error creating:', error);
		}
	}
	let selectedScope = $state({} as Scope | undefined);

	// const errorsStore = writable<Error[]>([]);

	// const machinesStore = writable<Machine[]>([]);
	// const machinesLoading = writable(true);
	// async function fetchMachines() {
	// 	try {
	// 		const response = await client.listMachines({});
	// 		machinesStore.set(response.machines);
	// 	} catch (error) {
	// 		console.error('Error fetching:', error);
	// 	} finally {
	// 		machinesLoading.set(false);
	// 	}
	// }

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchScopes();
			if ($scopesStore.length === 0) {
				await createDefaultScope();
				await fetchScopes();
			}
			selectedScope = $scopesStore.find((s) => s.name === 'default');
			// errorsStore.set([
			// 	{ code: 'CEPH_NOT_FOUND' } as Error,
			// 	{ code: 'KUBERNETES_NOT_FOUND' } as Error,
			// 	{ code: 'PROMETHEUS_NOT_FOUND' } as Error
			// ]);
			// await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<div class="flex flex-col gap-2">
		{@render SelectScope()}
		{#key selectedScope}
			{#if selectedScope}
				<Monitor scope={selectedScope} />
			{/if}
		{/key}
	</div>
{:else}
	<PageLoading />
{/if}

{#snippet SelectScope()}
	<div class="ml-auto flex items-center gap-2">
		<Label for="scope">Scope</Label>
		{#if $scopesLoading}
			<p>Loading scopes...</p>
		{:else}
			<Select.Root type="single">
				<Select.Trigger class="w-fit">
					{#if selectedScope}
						{selectedScope.name}
					{:else}
						Select
					{/if}
				</Select.Trigger>
				<Select.Content>
					{#if $scopesStore.length > 0}
						{#each $scopesStore as scope}
							<Select.Item
								value={scope.uuid}
								onclick={() => {
									selectedScope = scope;
								}}
							>
								{scope.name}
							</Select.Item>
						{/each}
					{/if}
				</Select.Content>
			</Select.Root>
		{/if}
	</div>
{/snippet}
