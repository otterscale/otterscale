<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { diskTypes, busTypes, dataVolumeSourceTypes } from './dropdown';

	import type {
		VirtualMachine,
		UpdateVirtualMachineRequest,
		VirtualMachineDisk,
		DataVolumeSource,
	} from '$lib/api/kubevirt/v1/kubevirt_pb';
	import {
		KubeVirtService,
		VirtualMachineDisk_type,
		VirtualMachineDisk_bus,
		DataVolumeSource_Type,
	} from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const KubeVirtClient = createClient(KubeVirtService, transport);

	// Form validation state
	let invalidNamespace: boolean | undefined = $state();

	// Label management state
	let labelKey = $state('');
	let labelValue = $state('');

	// ==================== Local Dropdown Options ====================
	const namespaces: Writable<SingleSelect.OptionType[]> = writable([]);

	// ==================== API Functions ====================
	async function loadNamespaces() {
		try {
			const response = await KubeVirtClient.listNamespaces({
				scopeUuid: $currentKubernetes?.scopeUuid,
				facilityName: $currentKubernetes?.name,
			});

			const namespaceOptions = response.namespaces.map((namespace) => ({
				value: namespace,
				label: namespace,
				icon: 'ph:folder',
			}));

			namespaces.set(namespaceOptions);
		} catch (error) {
			toast.error('Failed to load namespaces', {
				description: (error as ConnectError).message.toString(),
			});
		}
	}

	// ==================== Default Values & Constants ====================
	const DEFAULT_DISK_SOURCE = '';
	const DEFAULT_DISK_DATA_VOLUME_SOURCE = {
		type: DataVolumeSource_Type.HTTP,
		source: '',
		sizeBytes: 1n * 1024n * 1024n * 1024n, // 1GiB
	} as DataVolumeSource;
	const DEFAULT_DISK = {
		name: '',
		diskType: VirtualMachineDisk_type.DATAVOLUME,
		busType: VirtualMachineDisk_bus.VIRTIO,
		sourceData: { case: 'source', value: DEFAULT_DISK_SOURCE },
	} as VirtualMachineDisk;

	const defaults = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		name: virtualMachine.metadata?.name,
		namespace: virtualMachine.metadata?.namespace,
		networkName: virtualMachine.networkName,
		labels: virtualMachine.metadata?.labels || {},
		disks: virtualMachine.disks || [],
	} as UpdateVirtualMachineRequest;

	let request = $state(defaults);
	let newDisk: VirtualMachineDisk = $state(DEFAULT_DISK);
	let newDiskSource = $state(DEFAULT_DISK_SOURCE);
	let newDiskSourceDataVolume = $state(DEFAULT_DISK_DATA_VOLUME_SOURCE);

	function reset() {
		request = defaults;
		labelKey = '';
		labelValue = '';
		newDisk = DEFAULT_DISK;
		newDiskSource = DEFAULT_DISK_SOURCE;
		newDiskSourceDataVolume = DEFAULT_DISK_DATA_VOLUME_SOURCE;
	}

	// ==================== Disk Management ====================
	function addDisk() {
		if (newDisk.name.trim()) {
			if (newDisk.diskType === VirtualMachineDisk_type.DATAVOLUME) {
				newDisk.sourceData = {
					case: 'dataVolume',
					value: newDiskSourceDataVolume,
				};
			} else {
				newDisk.sourceData = {
					case: 'source',
					value: newDiskSource,
				};
			}
			request.disks = [...request.disks, { ...newDisk }];
			newDisk = DEFAULT_DISK;
			newDiskSource = DEFAULT_DISK_SOURCE;
			newDiskSourceDataVolume = DEFAULT_DISK_DATA_VOLUME_SOURCE;
		}
	}
	function removeDisk(index: number) {
		request.disks = request.disks.filter((_, i) => i !== index);
	}

	// ==================== Label Management ====================
	function addLabel() {
		if (labelKey.trim() && labelValue.trim()) {
			request.labels = { ...request.labels, [labelKey.trim()]: labelValue.trim() };
			labelKey = '';
			labelValue = '';
		}
	}
	function removeLabel(key: string) {
		const { [key]: _, ...rest } = request.labels;
		request.labels = rest;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	// ==================== Lifecycle Hooks ====================
	onMount(() => {
		loadNamespaces();
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="outline">
		<Icon icon="ph:pencil" />
		{m.edit()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit()} {m.virtual_machine()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Basic Configuration ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.virtual_machine_name()}</Form.Label>
					<SingleInput.General required readonly bind:value={request.name} />
				</Form.Field>
				<Form.Field>
					<Form.Label>Namespace</Form.Label>
					<SingleSelect.Root
						required
						options={namespaces}
						bind:value={request.namespace}
						bind:invalid={invalidNamespace}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $namespaces as namespace}
											<SingleSelect.Item option={namespace}>
												<Icon
													icon={namespace.icon ? namespace.icon : 'ph:empty'}
													class={cn('size-5', namespace.icon ? 'visible' : 'invisible')}
												/>
												{namespace.label}
												<SingleSelect.Check option={namespace} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>
				<Form.Field>
					<Form.Label>Network Name</Form.Label>
					<SingleInput.General bind:value={request.networkName} />
				</Form.Field>
			</Form.Fieldset>

			<!-- ==================== Disk Configuration ==================== -->
			<Form.Fieldset>
				<Form.Legend>Disks</Form.Legend>

				<!-- Display Current Disks -->
				{#if request.disks.length > 0}
					<div class="space-y-2">
						<h4 class="font-medium">Current Disks</h4>
						{#each request.disks as disk, index}
							<div class="bg-muted flex items-center justify-between rounded-md px-3 py-2">
								<div class="flex-1">
									<div class="flex items-center gap-2">
										<Icon icon="ph:hard-drive" class="size-4" />
										<span class="font-medium">{disk.name}</span>
									</div>
									<div class="text-muted-foreground text-sm">
										<span>Bus: {$busTypes.find((b) => b.value === disk.busType)?.label}</span>
										<span class="mx-2">•</span>
										<span>Type: {$diskTypes.find((t) => t.value === disk.diskType)?.label}</span>
										<span class="mx-2">•</span>
										{#if disk.diskType === VirtualMachineDisk_type.DATAVOLUME && disk.sourceData?.case === 'dataVolume'}
											<span
												>Source Type: {$dataVolumeSourceTypes.find(
													(s) => s.value === (disk.sourceData.value as DataVolumeSource).type,
												)?.label}</span
											>
											<span class="mx-2">•</span>
											<span
												>Size: {Math.floor(
													Number(disk.sourceData.value.sizeBytes) / (1024 * 1024 * 1024),
												)}GB</span
											>
											<span class="mx-2">•</span>
											<span>Source: {disk.sourceData.value.source}</span>
										{:else}
											<span
												>Source: {disk.sourceData?.case === 'source'
													? disk.sourceData.value
													: 'Unknown'}</span
											>
										{/if}
									</div>
								</div>
								<Button type="button" variant="ghost" size="sm" onclick={() => removeDisk(index)}>
									<Icon icon="ph:x" class="size-4" />
								</Button>
							</div>
						{/each}
					</div>
				{/if}

				<!-- Add New Disk Form -->
				<div class="space-y-4 border-t pt-4">
					<h4 class="font-medium">Add New Disk</h4>
					<Form.Field>
						<Form.Label>Disk Name</Form.Label>
						<SingleInput.General
							required
							type="text"
							placeholder="Enter disk name"
							bind:value={newDisk.name}
						/>
					</Form.Field>
					<Form.Field>
						<Form.Label>Bus Type</Form.Label>
						<SingleSelect.Root required options={busTypes} bind:value={newDisk.busType}>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.List>
										<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $busTypes as busType}
												<SingleSelect.Item option={busType}>
													<Icon
														icon={busType.icon ? busType.icon : 'ph:empty'}
														class={cn('size-5', busType.icon ? 'visible' : 'invisible')}
													/>
													{busType.label}
													<SingleSelect.Check option={busType} />
												</SingleSelect.Item>
											{/each}
										</SingleSelect.Group>
									</SingleSelect.List>
								</SingleSelect.Options>
							</SingleSelect.Content>
						</SingleSelect.Root>
					</Form.Field>
					<Form.Field>
						<Form.Label>Disk Type</Form.Label>
						<SingleSelect.Root required options={diskTypes} bind:value={newDisk.diskType}>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.List>
										<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $diskTypes as diskType}
												<SingleSelect.Item option={diskType}>
													<Icon
														icon={diskType.icon ? diskType.icon : 'ph:empty'}
														class={cn('size-5', diskType.icon ? 'visible' : 'invisible')}
													/>
													{diskType.label}
													<SingleSelect.Check option={diskType} />
												</SingleSelect.Item>
											{/each}
										</SingleSelect.Group>
									</SingleSelect.List>
								</SingleSelect.Options>
							</SingleSelect.Content>
						</SingleSelect.Root>
					</Form.Field>
					{#if newDisk.diskType === VirtualMachineDisk_type.DATAVOLUME}
						<Form.Field>
							<Form.Label>Source Type</Form.Label>
							<SingleSelect.Root
								required
								options={dataVolumeSourceTypes}
								bind:value={newDiskSourceDataVolume.type}
							>
								<SingleSelect.Trigger />
								<SingleSelect.Content>
									<SingleSelect.Options>
										<SingleSelect.List>
											<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
											<SingleSelect.Group>
												{#each $dataVolumeSourceTypes as sourceType}
													<SingleSelect.Item option={sourceType}>
														<Icon
															icon={sourceType.icon ? sourceType.icon : 'ph:empty'}
															class={cn(
																'size-5',
																sourceType.icon ? 'visible' : 'invisible',
															)}
														/>
														{sourceType.label}
														<SingleSelect.Check option={sourceType} />
													</SingleSelect.Item>
												{/each}
											</SingleSelect.Group>
										</SingleSelect.List>
									</SingleSelect.Options>
								</SingleSelect.Content>
							</SingleSelect.Root>
						</Form.Field>
						<Form.Field>
							<Form.Label>Source</Form.Label>
							<SingleInput.General
								required
								type="text"
								placeholder="Enter source reference"
								bind:value={newDiskSourceDataVolume.source}
							/>
						</Form.Field>
						<Form.Field>
							<Form.Label>Size</Form.Label>
							<SingleInput.Measurement
								required
								bind:value={newDiskSourceDataVolume.sizeBytes}
								transformer={(value) => String(value)}
								units={[{ value: 1024 * 1024 * 1024, label: 'GB' } as SingleInput.UnitType]}
							/>
						</Form.Field>
					{:else}
						<Form.Field>
							<Form.Label>Source</Form.Label>
							<SingleInput.General
								required
								type="text"
								placeholder="Enter source reference"
								bind:value={newDiskSource}
							/>
						</Form.Field>
					{/if}
					<Button
						type="button"
						variant="outline"
						size="sm"
						disabled={!newDisk.name.trim() ||
							(!newDiskSource.trim() && !newDiskSourceDataVolume.source.trim())}
						onclick={addDisk}
					>
						<Icon icon="ph:plus" class="size-4" />
						Add Disk
					</Button>
				</div>
			</Form.Fieldset>

			<!-- ==================== Labels Configuration ==================== -->
			<Form.Fieldset>
				<Form.Legend>Labels</Form.Legend>
				<Form.Field>
					<Form.Label>Labels</Form.Label>
					<div class="space-y-2">
						<div class="flex gap-2">
							<SingleInput.General type="text" placeholder="Key" bind:value={labelKey} class="flex-1" />
							<SingleInput.General
								type="text"
								placeholder="Value"
								bind:value={labelValue}
								class="flex-1"
							/>
							<Button
								type="button"
								variant="outline"
								size="sm"
								disabled={!labelKey.trim() || !labelValue.trim()}
								onclick={addLabel}
							>
								<Icon icon="ph:plus" class="size-4" />
								Add
							</Button>
						</div>
						{#if Object.keys(request.labels).length > 0}
							<div class="space-y-1">
								{#each Object.entries(request.labels) as [key, value]}
									<div class="bg-muted flex items-center justify-between rounded-md px-3 py-2">
										<span class="text-sm">
											<span class="font-medium">{key}</span>: {value}
										</span>
										<Button
											type="button"
											variant="ghost"
											size="sm"
											onclick={() => removeLabel(key)}
										>
											<Icon icon="ph:x" class="size-4" />
										</Button>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalidNamespace}
					onclick={() => {
						toast.promise(() => KubeVirtClient.updateVirtualMachine(request), {
							loading: `Updating ${virtualMachine.metadata?.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully updated ${virtualMachine.metadata?.name}`;
							},
							error: (error) => {
								let message = `Failed to update ${virtualMachine.metadata?.name}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
								});
								return message;
							},
						});
						reset();
						close();
					}}
				>
					{m.save()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
