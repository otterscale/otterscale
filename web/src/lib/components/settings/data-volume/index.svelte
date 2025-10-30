<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import Actions from './cell-actions.svelte';
	import Create from './create.svelte';

	import { DataVolume_Source_Type, InstanceService, type DataVolume } from '$lib/api/instance/v1/instance_pb';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { scope, facility, namespace }: { scope: string; facility: string; namespace: string } = $props();

	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(InstanceService, transport);

	const dataVolumes = writable<DataVolume[]>();
	let isConfigurationLoading = $state(true);

	const reloadManager = new ReloadManager(() => {
		virtualMachineClient
			.listDataVolumes({
				scope: scope,
				facility: facility,
				namespace: namespace,
				bootImage: true,
			})
			.then((response) => {
				dataVolumes.set(response.dataVolumes);
			});
	});
	setContext('reloadManager', reloadManager);

	onMount(async () => {
		try {
			await virtualMachineClient
				.listDataVolumes({
					scope: scope,
					facility: facility,
					namespace: namespace,
					bootImage: true,
				})
				.then((response) => {
					dataVolumes.set(response.dataVolumes);
					isConfigurationLoading = false;
				});

			reloadManager.start();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});

	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if !isConfigurationLoading}
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
			<div class="rounded-lg border shadow-sm">
				<Table.Root>
					<Table.Header>
						<Table.Row class="[&_th]:bg-muted *:px-4 [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg">
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
						{#each $dataVolumes as dataVolume}
							<Table.Row class="*:px-4">
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
												{dataVolume.phase}
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
																	<Table.Head class="text-right"
																		>{m.type()}</Table.Head
																	>
																	<Table.Cell
																		>{DataVolume_Source_Type[
																			dataVolume.source.type
																		]}</Table.Cell
																	>
																</Table.Row>
																<Table.Row>
																	<Table.Head class="text-right"
																		>{m.source()}</Table.Head
																	>
																	<Table.Cell class="break-all"
																		>{dataVolume.source.data}</Table.Cell
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
								<Table.Cell class="p-0">
									<Actions {dataVolume} />
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</Layout.Viewer>
	</Layout.Root>
{/if}
