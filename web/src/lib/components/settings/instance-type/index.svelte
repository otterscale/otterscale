<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import Actions from './cell-actions.svelte';
	import Create from './create.svelte';

	import { InstanceService, type InstanceType } from '$lib/api/instance/v1/instance_pb';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { scope, facility, namespace }: { scope: string; facility: string; namespace: string } = $props();

	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(InstanceService, transport);

	const instanceTypes = writable<InstanceType[]>();
	let isInstanceTypesLoading = $state(true);

	const reloadManager = new ReloadManager(() => {
		virtualMachineClient
			.listInstanceTypes({
				scope: scope,
				facility: facility,
				namespace: namespace,
				includeClusterWide: false,
			})
			.then((response) => {
				instanceTypes.set(response.instanceTypes);
			});
	});
	setContext('reloadManager', reloadManager);

	onMount(async () => {
		try {
			const response = await virtualMachineClient.listInstanceTypes({
				scope: scope,
				facility: facility,
				namespace: namespace,
				includeClusterWide: false,
			});
			instanceTypes.set(response.instanceTypes);
			isInstanceTypesLoading = false;

			reloadManager.start();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});

	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if !isInstanceTypesLoading}
	<Layout.Root>
		<Layout.Title>{m.instance_type()}</Layout.Title>
		<Layout.Description>
			{m.setting_instance_type_description()}
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
			<div class="rounded-lg border shadow-sm w-full">
				<Table.Root>
					<Table.Header>
						<Table.Row class="[&_th]:bg-muted *:px-4 [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg">
							<Table.Head>{m.name()}</Table.Head>
							<Table.Head>{m.namespace()}</Table.Head>
							<Table.Head class="text-right">{m.cpu_cores()}</Table.Head>
							<Table.Head class="text-right">{m.memory()}</Table.Head>
							<Table.Head class="text-right">{m.create_time()}</Table.Head>
							<Table.Head></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $instanceTypes as instanceType}
							<Table.Row class="*:px-4">
								<Table.Cell>{instanceType.name}</Table.Cell>
								<Table.Cell>
									{#if instanceType.namespace}
										<Badge variant="outline">{instanceType.namespace}</Badge>
									{/if}
								</Table.Cell>
								<Table.Cell class="text-right">
									<span>{instanceType.cpuCores}</span>
								</Table.Cell>
								<Table.Cell class="text-right">
									{@const { value: memory, unit } = formatCapacity(instanceType.memoryBytes)}
									<span>{memory} {unit}</span>
								</Table.Cell>
								<Table.Cell class="text-right">
									{#if instanceType.createdAt}
										<Tooltip.Provider>
											<Tooltip.Root>
												<Tooltip.Trigger>
													{formatTimeAgo(timestampDate(instanceType.createdAt))}
												</Tooltip.Trigger>
												<Tooltip.Content>
													{timestampDate(instanceType.createdAt)}
												</Tooltip.Content>
											</Tooltip.Root>
										</Tooltip.Provider>
									{/if}
								</Table.Cell>
								<Table.Cell class="p-0">
									<Actions {instanceType} />
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</Layout.Viewer>
	</Layout.Root>
{/if}
