<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import {
		type DataVolume,
		DataVolume_Source_Type,
		InstanceService
	} from '$lib/api/instance/v1/instance_pb';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	import Actions from './cell-actions.svelte';
	import Create from './create.svelte';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(InstanceService, transport);
	const dataVolumes = writable<DataVolume[]>();

	async function fetch() {
		try {
			const response = await virtualMachineClient.listDataVolumes({
				scope: scope,
				namespace: '',
				bootImage: true
			});
			dataVolumes.set(response.dataVolumes);
		} catch (error) {
			console.error('Error fetching data volumes:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isMounted = $state(false);
	onMount(async () => {
		await fetch();
		isMounted = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if isMounted}
	<Layout.Root>
		<Layout.Title>{m.virtual_machine_data_volume()}</Layout.Title>
		<Layout.Description>
			{m.setting_data_volume_description()}
		</Layout.Description>
		<Layout.Controller>
			<Create {scope} {reloadManager} />
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
						<Table.Row>
							<Table.Head>{m.name()}</Table.Head>
							<Table.Head>{m.namespace()}</Table.Head>
							<Table.Head>{m.phase()}</Table.Head>
							<Table.Head>{m.source()}</Table.Head>
							<Table.Head class="text-right">{m.progress()}</Table.Head>
							<Table.Head class="text-right">{m.size()}</Table.Head>
							<Table.Head></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $dataVolumes as dataVolume (`${dataVolume.namespace}/${dataVolume.name}`)}
							<Table.Row>
								<Table.Cell>{dataVolume.name}</Table.Cell>
								<Table.Cell><Badge variant="outline">{dataVolume.namespace}</Badge></Table.Cell>
								<Table.Cell>
									<Tooltip.Provider>
										<Tooltip.Root>
											<Tooltip.Trigger>
												{#if dataVolume.phase === 'Succeeded'}
													<Icon icon="ph:check" class="text-green-600" />
												{:else if dataVolume.phase === 'ImportInProgress'}
													<Icon icon="ph:spinner" class="animate-spin text-blue-600" />
												{:else}
													<Icon icon="ph:x" class="text-gray-400" />
												{/if}
											</Tooltip.Trigger>
											<Tooltip.Content>
												{dataVolume.phase.replace(/(?<!^)([A-Z])/g, ' $1')}
											</Tooltip.Content>
										</Tooltip.Root>
									</Tooltip.Provider>
								</Table.Cell>
								<Table.Cell class="items-start">
									{#if dataVolume.source}
										<div class="flex items-center gap-1">
											{#if dataVolume.source.data}
												<HoverCard.Root>
													<HoverCard.Trigger>
														<!-- <Icon icon="ph:info" /> -->
														<Badge variant="outline">
															{#if dataVolume.source.type === DataVolume_Source_Type.BLANK_IMAGE}
																<Icon icon="ph:file-blank" class="mr-1" />
																BLANK IMAGE
															{:else if dataVolume.source.type === DataVolume_Source_Type.HTTP_URL}
																<Icon icon="ph:file-cloud" class="mr-1" />
																HTTP URL
															{:else if dataVolume.source.type === DataVolume_Source_Type.EXISTING_PERSISTENT_VOLUME_CLAIM}
																<Icon icon="ph:hard-drive" class="mr-1" />
																PVC
															{:else}
																<Icon icon="ph:file-blank" class="mr-1" />
																{DataVolume_Source_Type[dataVolume.source.type]}
															{/if}
														</Badge>
													</HoverCard.Trigger>
													<HoverCard.Content class="min-w-[600px]">
														<Table.Root>
															<Table.Body class="text-xs">
																<Table.Row>
																	<Table.Head class="text-right">{m.type()}</Table.Head>
																	<Table.Cell
																		>{DataVolume_Source_Type[dataVolume.source.type]}</Table.Cell
																	>
																</Table.Row>
																<Table.Row>
																	<Table.Head class="text-right">{m.source()}</Table.Head>
																	<Table.Cell class="break-all">{dataVolume.source.data}</Table.Cell
																	>
																</Table.Row>
															</Table.Body>
														</Table.Root>
													</HoverCard.Content>
												</HoverCard.Root>
											{/if}
										</div>
									{/if}
								</Table.Cell>
								<Table.Cell>
									<span class="flex items-center justify-end gap-1">{dataVolume.progress}</span>
								</Table.Cell>
								<Table.Cell>
									{@const { value: capacity, unit } = formatCapacity(dataVolume.sizeBytes)}
									<span class="flex items-center justify-end gap-1">{capacity} {unit}</span>
								</Table.Cell>
								<Table.Cell>
									<Actions {dataVolume} {scope} {reloadManager} />
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</Layout.Viewer>
	</Layout.Root>
{/if}
