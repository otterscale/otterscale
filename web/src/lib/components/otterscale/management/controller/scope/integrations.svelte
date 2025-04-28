<!-- <script lang="ts">
	import * as Table from '$lib/components/ui/table';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { StackService, type Integration } from '$gen/api/stack/v1/stack_pb';
	import type { Plugin } from 'svelte-exmarkdown';
	import rehypeHighlight from 'rehype-highlight';
	import { Badge } from '$lib/components/ui/badge';
	import 'highlight.js/styles/github.css';
	import Icon from '@iconify/svelte';
	import {
		Drawer,
		DrawerContent,
		DrawerDescription,
		DrawerHeader,
		DrawerTitle
	} from '$lib/components/ui/drawer';

	const transport: Transport = getContext('transport');
	const client = createClient(StackService, transport);

	const integrationsStore = writable<Integration[]>([]);
	const integrationsIsLoading = writable(true);
	async function fetchIntegrations() {
		try {
			const response = await client.listIntegrations({
				modelUuid: model_uuid
			});
			integrationsStore.set(response.integrations);
		} catch (error) {
			console.error('Error fetching machines:', error);
		} finally {
			integrationsIsLoading.set(false);
		}
	}

	let {
		model_uuid,
		open = $bindable()
	}: {
		model_uuid: string;
		open: boolean;
	} = $props();

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchIntegrations();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<Drawer {open} onOpenChange={() => (open = false)}>
		<DrawerContent class="inset-x-auto inset-y-0 right-0 w-3/5 px-3">
			<DrawerHeader>
				<DrawerTitle class="text-center">Integrations</DrawerTitle>
			</DrawerHeader>
			<div class="overflow-y-auto px-3">
				{#if $integrationsIsLoading}
					<div class="flex items-center justify-center gap-2">
						<Icon icon="ph:spinner" class="size-4 animate-spin" />
						Loading integrations...
					</div>
				{:else}
					<Table.Root>
						<Table.Header>
							<Table.Row class="*:text-xs *:font-light">
								<Table.Head>RELATION PROVIDER</Table.Head>
								<Table.Head>REQUIRER</Table.Head>
								<Table.Head>INTERFACE</Table.Head>
								<Table.Head>ROLE</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each $integrationsStore as integration}
								<Table.Row class="*:text-xs">
									<Table.Cell>{integration.provider}</Table.Cell>
									<Table.Cell>{integration.requirer}</Table.Cell>
									<Table.Cell>{integration.interface}</Table.Cell>
									<Table.Cell
										><Badge variant="outline">
											{integration.role}
										</Badge>
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				{/if}
			</div>
		</DrawerContent>
	</Drawer>
{:else}
	<div class="flex h-full w-full items-center justify-center gap-2 text-sm text-muted-foreground">
		<Icon icon="ph:spinner" class="size-8 animate-spin" />
		Loading...
	</div>
{/if} -->
