<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { resourcesCase, diskTypes, busTypes, dataVolumeSourceTypes } from './dropdown';

	import type {
		CreateVirtualMachineRequest,
		VirtualMachineResources,
		VirtualMachineDisk,
		DataVolumeSource,
	} from '$lib/api/kubevirt/v1/kubevirt_pb';
	import {
		KubeVirtService,
		VirtualMachineDisk_type,
		VirtualMachineDisk_bus,
		DataVolumeSource_Type,
	} from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Code from '$lib/components/custom/code';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';

	// Context dependencies
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const kubevirtClient = createClient(KubeVirtService, transport);

	// ==================== State Variables ====================

	// UI state
	let isAdvancedOpen = $state(false);
	let open = $state(false);

	// Form validation state
	let invalidName: boolean | undefined = $state();
	let invalidNamespace: boolean | undefined = $state();
	let invalidResourceCase: boolean | undefined = $state();

	// Label management state
	let labelKey = $state('');
	let labelValue = $state('');

	// ==================== Local Dropdown Options ====================
	const namespaces: Writable<SingleSelect.OptionType[]> = writable([]);

	// ==================== API Functions ====================
	async function loadNamespaces() {
		try {
			const response = await kubevirtClient.listNamespaces({
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

	// Default request structure for creating a virtual machine
	const DEFAULT_REQUEST = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		name: '',
		namespace: '',
		networkName: '',
		startupScript: '',
		labels: {},
		disks: [] as VirtualMachineDisk[],
		resources: { case: 'instancetypeName', value: '' },
	} as CreateVirtualMachineRequest;
	const DEFAULT_RESOURCES_CUSTOM = {
		cpuCores: 1,
		memoryBytes: 1n * 1024n * 1024n * 1024n, // 1GiB
	} as VirtualMachineResources;
	const DEFAULT_RESOURCES_INSTANCE = '';
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

	// ==================== Form State ====================
	let request: CreateVirtualMachineRequest = $state(DEFAULT_REQUEST);
	let resourcesCustom = $state(DEFAULT_RESOURCES_CUSTOM);
	let resourcesInstance = $state(DEFAULT_RESOURCES_INSTANCE);
	let newDisk: VirtualMachineDisk = $state(DEFAULT_DISK);
	let newDiskSource = $state(DEFAULT_DISK_SOURCE);
	let newDiskSourceDataVolume = $state(DEFAULT_DISK_DATA_VOLUME_SOURCE);

	// ==================== Reactive Statements ====================
	// Automatically sync request.resources.value with the form state
	$effect(() => {
		if (request.resources.case === 'instancetypeName') {
			request.resources.value = resourcesInstance;
		} else if (request.resources.case === 'custom') {
			request.resources.value = resourcesCustom;
		}
	});

	// ==================== Utility Functions ====================
	function reset() {
		request = DEFAULT_REQUEST;
		resourcesCustom = DEFAULT_RESOURCES_CUSTOM;
		resourcesInstance = DEFAULT_RESOURCES_INSTANCE;
		isAdvancedOpen = false;
		labelKey = '';
		labelValue = '';
		newDisk = DEFAULT_DISK;
		newDiskSource = DEFAULT_DISK_SOURCE;
		newDiskSourceDataVolume = DEFAULT_DISK_DATA_VOLUME_SOURCE;
	}
	function close() {
		open = false;
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

	// ==================== Lifecycle Hooks ====================
	onMount(() => {
		loadNamespaces();
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_virtual_machine()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Basic Configuration ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.name} bind:invalid={invalidName} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>
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
			</Form.Fieldset>

			<!-- ==================== Resource Configuration ==================== -->
			<Form.Fieldset>
				<Form.Legend>{m.resources()}</Form.Legend>
				<Form.Field>
					<Form.Label>{m.type()}</Form.Label>
					<SingleSelect.Root
						required
						options={resourcesCase}
						bind:value={request.resources.case}
						bind:invalid={invalidResourceCase}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $resourcesCase as type}
											<SingleSelect.Item option={type}>
												<Icon
													icon={type.icon ? type.icon : 'ph:empty'}
													class={cn('size-5', type.icon ? 'visible' : 'invisible')}
												/>
												{type.label}
												<SingleSelect.Check option={type} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>
				{#if request.resources.case === 'custom'}
					<Form.Field>
						<Form.Label>{m.cpu_cores()}</Form.Label>
						<SingleInput.General required type="number" bind:value={resourcesCustom.cpuCores} />
					</Form.Field>
					<Form.Field>
						<Form.Label>{m.memory()}</Form.Label>
						<SingleInput.Measurement
							required
							bind:value={resourcesCustom.memoryBytes}
							transformer={(value) => String(value)}
							units={[{ value: 1024 * 1024 * 1024, label: 'GB' } as SingleInput.UnitType]}
						/>
					</Form.Field>
				{:else if request.resources.case === 'instancetypeName'}
					<Form.Field>
						<Form.Label>{m.instance_name()}</Form.Label>
						<SingleInput.General required type="text" bind:value={resourcesInstance} />
					</Form.Field>
				{/if}
			</Form.Fieldset>

			<!-- ==================== Disk Configuration ==================== -->
			<Form.Fieldset>
				<Form.Legend>{m.disk()}</Form.Legend>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General required type="text" placeholder="Enter disk name" bind:value={newDisk.name} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.bus_type()}</Form.Label>
					<Form.Help>
						{m.vm_bus_type_direction()}
					</Form.Help>
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
					<Form.Label>{m.disk_type()}</Form.Label>
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
						<Form.Label>{m.source_type()}</Form.Label>
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
														class={cn('size-5', sourceType.icon ? 'visible' : 'invisible')}
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
						<Form.Label>{m.source()}</Form.Label>
						<SingleInput.General
							required
							type="text"
							placeholder="Enter source reference"
							bind:value={newDiskSourceDataVolume.source}
						/>
					</Form.Field>
					<Form.Field>
						<Form.Label>{m.size()}</Form.Label>
						<SingleInput.Measurement
							required
							bind:value={newDiskSourceDataVolume.sizeBytes}
							transformer={(value) => String(value)}
							units={[{ value: 1024 * 1024 * 1024, label: 'GB' } as SingleInput.UnitType]}
						/>
					</Form.Field>
				{:else}
					<Form.Field>
						<Form.Label>{m.source()}</Form.Label>
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
					disabled={!newDisk.name.trim() || (!newDiskSource.trim() && !newDiskSourceDataVolume.source.trim())}
					onclick={addDisk}
				>
					<Icon icon="ph:plus" class="size-4" />
					Add Disk
				</Button>
				<!-- Display Configured Disks -->
				{#if request.disks.length > 0}
					<div class="space-y-2">
						<h4 class="font-medium">Configured Disks</h4>
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
			</Form.Fieldset>

			<!-- ==================== Advanced Configuration ==================== -->
			<Collapsible.Root bind:open={isAdvancedOpen} class="py-4">
				<div class="flex items-center justify-between gap-2">
					<p class={cn('text-base font-bold', isAdvancedOpen ? 'invisible' : 'visible')}>{m.advance()}</p>
					<Collapsible.Trigger class="bg-muted rounded-full p-1 ">
						<Icon
							icon="ph:caret-left"
							class={cn('transition-all duration-300', isAdvancedOpen ? '-rotate-90' : 'rotate-0')}
						/>
					</Collapsible.Trigger>
				</div>
				<Collapsible.Content>
					<Form.Fieldset>
						<Form.Legend>{m.advance()}</Form.Legend>
						<Form.Field>
							<Form.Label>{m.network_name()}</Form.Label>
							<SingleInput.General type="text" bind:value={request.networkName} />
						</Form.Field>
						<Form.Field>
							<Form.Label>{m.labels()}</Form.Label>
							<div class="space-y-2">
								<div class="flex gap-2">
									<SingleInput.General
										type="text"
										placeholder="Key"
										bind:value={labelKey}
										class="flex-1"
									/>
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
											<div
												class="bg-muted flex items-center justify-between rounded-md px-3 py-2"
											>
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
						<Form.Field>
							<Form.Label>{m.startup_script()}</Form.Label>
							<Code.Root lang="bash" class="w-full" hideLines code={request.startupScript}>
								<Code.CopyButton />
							</Code.Root>
							<SingleInput.Structure preview={false} bind:value={request.startupScript} language="bash" />
						</Form.Field>
					</Form.Fieldset>
				</Collapsible.Content>
			</Collapsible.Root>
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
					disabled={invalidName || invalidNamespace}
					onclick={() => {
						toast.promise(() => kubevirtClient.createVirtualMachine(request), {
							loading: `Creating ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to create ${request.name}`;
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
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
