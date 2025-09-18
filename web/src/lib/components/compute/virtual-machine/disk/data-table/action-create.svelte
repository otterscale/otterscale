<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		DataVolumeSource_Type,
		KubeVirtService,
		VirtualMachineDisk_bus,
		VirtualMachineDisk_type,
	} from '$lib/api/kubevirt/v1/kubevirt_pb';
	import type {
		DataVolumeSource,
		PersistentVolumeClaim,
		VirtualMachine,
		VirtualMachineDisk,
		CreateVirtualMachineDiskRequest,
	} from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { busTypes, dataVolumeSourceTypes, diskTypes } from '$lib/components/compute/virtual-machine/units/dropdown';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';

	let {
		virtualMachine,
	}: {
		virtualMachine: VirtualMachine;
	} = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const kubevirtClient = createClient(KubeVirtService, transport);

	// ==================== State Variables ====================

	// UI state
	let open = $state(false);

	// Form validation state
	let invalidName: boolean | undefined = $state();
	let invalidSource: boolean | undefined = $state();

	// ==================== Local Dropdown Options ====================
	const namespaces: Writable<SingleSelect.OptionType[]> = writable([]);
	const bootablePVCs: Writable<SingleSelect.OptionType[]> = writable([]);

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

	async function loadBootablePVCs() {
		try {
			if (!request.vmNamespace) return;

			const response = await kubevirtClient.listBootablePersistentVolumeClaims({
				scopeUuid: $currentKubernetes?.scopeUuid,
				facilityName: $currentKubernetes?.name,
				namespace: request.vmNamespace,
			});

			const pvcOptions = response.persistentVolumeClaims.map((pvc: PersistentVolumeClaim) => ({
				value: pvc.name,
				label: pvc.name,
				icon: 'ph:hard-drive',
			}));

			bootablePVCs.set(pvcOptions);
		} catch (error) {
			toast.error('Failed to load bootable PVCs', {
				description: (error as ConnectError).message.toString(),
			});
		}
	}

	// ==================== Default Values & Constants ====================

	// Default request structure for creating a virtual machine
	const DEFAULT_REQUEST = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		vmName: virtualMachine.metadata?.name,
		vmNamespace: virtualMachine.metadata?.namespace,
		disk: {
			name: '',
			diskType: VirtualMachineDisk_type.DATAVOLUME,
			busType: VirtualMachineDisk_bus.VIRTIO,
			sourceData: { case: 'source', value: '' },
			isBootable: true,
		} as VirtualMachineDisk,
	} as CreateVirtualMachineDiskRequest;
	const DEFAULT_DISK_DATA_VOLUME_SOURCE = {
		type: DataVolumeSource_Type.HTTP,
		source: '',
		sizeBytes: 1n * 1024n * 1024n * 1024n, // 1GiB
	} as DataVolumeSource;

	// ==================== Form State ====================
	let request: CreateVirtualMachineDiskRequest = $state(DEFAULT_REQUEST);
	let diskSource = $state('');
	let diskSourceDataVolume = $state(DEFAULT_DISK_DATA_VOLUME_SOURCE);

	// ==================== Reactive Statements ====================
	// Load bootable PVCs when namespace changes
	$effect(() => {
		if (request.vmNamespace) {
			loadBootablePVCs();
		}
	});

	// Auto-set bootable based on data volume source type
	$effect(() => {
		if (request.disk.diskType === VirtualMachineDisk_type.DATAVOLUME) {
			if (
				diskSourceDataVolume.type === DataVolumeSource_Type.HTTP ||
				diskSourceDataVolume.type === DataVolumeSource_Type.PVC
			) {
				request.disk.isBootable = true;
			} else if (diskSourceDataVolume.type === DataVolumeSource_Type.BLANK) {
				request.disk.isBootable = false;
			}
			request.disk.sourceData = {
				case: 'dataVolume',
				value: diskSourceDataVolume,
			};
		} else {
			request.disk.sourceData = {
				case: 'source',
				value: diskSource,
			};
		}
	});

	// ==================== Utility Functions ====================
	function reset() {
		request = DEFAULT_REQUEST;
		diskSource = '';
		diskSourceDataVolume = DEFAULT_DISK_DATA_VOLUME_SOURCE;
		bootablePVCs.set([]);
	}
	function close() {
		open = false;
	}

	// ==================== Lifecycle Hooks ====================
	onMount(() => {
		loadNamespaces();
		if (request.vmNamespace) {
			loadBootablePVCs();
		}
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
			<!-- ==================== Disk Configuration ==================== -->
			<Form.Fieldset>
				<Form.Legend>{m.disk()}</Form.Legend>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.disk.name}
						bind:invalid={invalidName}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.bus_type()}</Form.Label>
					<Form.Help>
						{m.vm_bus_type_direction()}
					</Form.Help>
					<SingleSelect.Root required options={busTypes} bind:value={request.disk.busType}>
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
					<SingleSelect.Root required options={diskTypes} bind:value={request.disk.diskType}>
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
				{#if request.disk.diskType === VirtualMachineDisk_type.DATAVOLUME}
					<Form.Field>
						<Form.Label>{m.source_type()}</Form.Label>
						<SingleSelect.Root
							required
							options={dataVolumeSourceTypes}
							bind:value={diskSourceDataVolume.type}
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
						<SingleInput.Boolean
							descriptor={() => m.boot_disk()}
							bind:value={request.disk.isBootable}
							disabled={request.disk.diskType === VirtualMachineDisk_type.DATAVOLUME &&
								diskSourceDataVolume.type === DataVolumeSource_Type.BLANK}
						/>
					</Form.Field>
					<Form.Field>
						<Form.Label>{m.size()}</Form.Label>
						<SingleInput.Measurement
							required
							bind:value={diskSourceDataVolume.sizeBytes}
							transformer={(value) => String(value)}
							units={[{ value: 1024 * 1024 * 1024, label: 'GB' } as SingleInput.UnitType]}
						/>
					</Form.Field>
					<Form.Field>
						<Form.Label>{m.source()}</Form.Label>
						{#if diskSourceDataVolume.type === DataVolumeSource_Type.PVC}
							<SingleSelect.Root
								required
								options={bootablePVCs}
								bind:value={diskSourceDataVolume.source}
								bind:invalid={invalidSource}
							>
								<SingleSelect.Trigger />
								<SingleSelect.Content>
									<SingleSelect.Options>
										<SingleSelect.Input />
										<SingleSelect.List>
											<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
											<SingleSelect.Group>
												{#each $bootablePVCs as pvc}
													<SingleSelect.Item option={pvc}>
														<Icon
															icon={pvc.icon ? pvc.icon : 'ph:empty'}
															class={cn('size-5', pvc.icon ? 'visible' : 'invisible')}
														/>
														{pvc.label}
														<SingleSelect.Check option={pvc} />
													</SingleSelect.Item>
												{/each}
											</SingleSelect.Group>
										</SingleSelect.List>
									</SingleSelect.Options>
								</SingleSelect.Content>
							</SingleSelect.Root>
						{:else}
							<SingleInput.General
								required={diskSourceDataVolume.type === DataVolumeSource_Type.HTTP}
								disabled={diskSourceDataVolume.type === DataVolumeSource_Type.BLANK}
								type="text"
								bind:value={diskSourceDataVolume.source}
								bind:invalid={invalidSource}
								placeholder={diskSourceDataVolume.type === DataVolumeSource_Type.HTTP
									? 'https://cloud-images.ubuntu.com/xxx/xxx/xxx.img'
									: ''}
							/>
						{/if}
					</Form.Field>
				{:else}
					<Form.Field>
						<SingleInput.Boolean descriptor={() => m.boot_disk()} bind:value={request.disk.isBootable} />
					</Form.Field>
					<Form.Field>
						<Form.Label>{m.source()}</Form.Label>
						<SingleInput.General
							required
							type="text"
							bind:value={diskSource}
							bind:invalid={invalidSource}
						/>
					</Form.Field>
				{/if}
				{#if diskSourceDataVolume.type === DataVolumeSource_Type.HTTP}
					<div class="flex justify-end gap-2">
						<Button
							variant="outline"
							size="sm"
							href="https://cloud-images.ubuntu.com/"
							target="_blank"
							class="flex items-center gap-1"
						>
							<Icon icon="ph:arrow-square-out" />
							{m.cloud_image()}
						</Button>
					</div>
				{/if}
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
					disabled={invalidName || invalidSource}
					onclick={() => {
						toast.promise(() => kubevirtClient.createVirtualMachineDisk(request), {
							loading: `Creating ${request.vmName}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created ${request.vmName}`;
							},
							error: (error) => {
								let message = `Failed to create ${request.vmName}`;
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
