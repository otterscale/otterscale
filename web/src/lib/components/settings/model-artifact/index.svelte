<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	// import Icon from '@iconify/svelte';
	// import { mode } from 'mode-watcher';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import Actions from './cell-actions.svelte';
	import Create from './create.svelte';

	// import { DataVolume_Source_Type } from '$lib/api/instance/v1/instance_pb';
	import type { ModelArtifact } from '$lib/api/model/v1/model_pb';
	import { ModelService } from '$lib/api/model/v1/model_pb';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/layout';
	import { Badge } from '$lib/components/ui/badge';
	// import * as HoverCard from '$lib/components/ui/hover-card';
	// import * as Tooltip from '$lib/components/ui/tooltip';
	// import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { scope, facility, namespace }: { scope: string; facility: string; namespace: string } = $props();

	const transport: Transport = getContext('transport');
	const modelClient = createClient(ModelService, transport);

	const modelArtifacts = writable<ModelArtifact[]>([]);
	let isLoading = $state(true);

	async function fetch() {
		modelClient
			.listModelArtifacts({
				scope: scope,
				facility: facility,
				namespace: namespace,
			})
			.then((response) => {
				modelArtifacts.set(response.modelArtifacts);
			})
			.catch((error) => {
				console.error('Error reloading model artifacts:', error);
			});
	}

	const reloadManager = new ReloadManager(fetch);
	setContext('reloadManager', reloadManager);

	onMount(async () => {
		try {
			await fetch();
			isLoading = false;
			reloadManager.start();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if !isLoading}
	<Layout.Root>
		<Layout.Title>{m.virtual_machine_data_volume()}</Layout.Title>
		<Layout.Description>
			{m.setting_data_volume_description()}
		</Layout.Description>
		<Layout.Controller>
			<Create />
			<Reloader
				bind:checked={reloadManager.state}
				onCheckedChange={() => {
					if (reloadManager.state) {
						reloadManager.restart();
					} else {
						reloadManager.stop();
					}
				}}
			/>
		</Layout.Controller>
		<Layout.Viewer>
			<div class="w-full rounded-lg border shadow-sm">
				<Table.Root>
					<Table.Header>
						<Table.Row class="[&_th]:bg-muted *:px-4 [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg">
							<Table.Head>{m.name()}</Table.Head>
							<Table.Head>{m.namespace()}</Table.Head>
							<Table.Head>{m.model_name()}</Table.Head>
							<Table.Head>{m.phase()}</Table.Head>
							<Table.Head class="text-right">{m.size()}</Table.Head>
							<Table.Head>{m.volume()}</Table.Head>
							<Table.Head class="text-right">{m.time()}</Table.Head>
							<Table.Head></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $modelArtifacts as modelArtifact}
							<Table.Row class="*:px-4">
								<Table.Cell>{modelArtifact.name}</Table.Cell>
								<Table.Cell><Badge variant="outline">{modelArtifact.namespace}</Badge></Table.Cell>
								<Table.Cell>{modelArtifact.modelName}</Table.Cell>
								<Table.Cell><Badge variant="outline">{modelArtifact.phase}</Badge></Table.Cell>
								<Table.Cell class="text-right">
									{modelArtifact.size}
								</Table.Cell>
								<Table.Cell>
									{modelArtifact.volumeName}
								</Table.Cell>
								<Table.Cell class="text-right">
									{modelArtifact.createdAt}
								</Table.Cell>
								<Table.Cell>
									<Actions {modelArtifact} />
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</Layout.Viewer>
	</Layout.Root>
{/if}
