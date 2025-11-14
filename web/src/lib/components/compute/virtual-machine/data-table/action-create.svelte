<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { CreateVirtualMachineRequest, DataVolume } from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Code from '$lib/components/custom/code';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';

	// Context dependencies
	const transport: Transport = getContext('transport');

	const virtualMachineClient = createClient(InstanceService, transport);

	// ==================== State Variables ====================

	// UI state
	let isAdvancedOpen = $state(false);
	let open = $state(false);

	// Form validation state
	let invalidName: boolean | undefined = $state();
	let invalidInstanceTypeName: boolean | undefined = $state();
	let invalidBootDataVolumeName: boolean | undefined = $state();

	// ==================== Local Dropdown Options ====================
	const bootDataVolumes: Writable<SingleSelect.OptionType[]> = writable([]);
	const instanceTypes: Writable<SingleSelect.OptionType[]> = writable([]);

	// Instance type with CPU and memory information
	type InstanceTypeOption = SingleSelect.OptionType & {
		cpuCores?: number;
		memoryBytes?: bigint;
	};

	type BootDataVolumesOption = SingleSelect.OptionType & {
		sizeBytes?: bigint;
	};

	// ==================== API Functions ====================
	async function loadInstanceTypes() {
		try {
			// Request both namespace-specific and cluster-wide instance types
			const response = await virtualMachineClient.listInstanceTypes({
				scope: scope,
				,
				namespace: request.namespace,
				includeClusterWide: true
			});

			const instanceTypeOptions: InstanceTypeOption[] = response.instanceTypes.map(
				(instanceType) => {
					const memory = formatCapacity(instanceType.memoryBytes);
					return {
						value: instanceType.name,
						label: `${instanceType.name} (CPU: ${instanceType.cpuCores} Core, RAM: ${memory.value} ${memory.unit})`,
						icon: instanceType.clusterWide ? 'ph:graph' : 'ph:layout',
						cpuCores: instanceType.cpuCores,
						memoryBytes: instanceType.memoryBytes
					};
				}
			);

			instanceTypes.set(instanceTypeOptions);
		} catch (error) {
			toast.error('Failed to load instance types', {
				description: (error as ConnectError).message.toString()
			});
		}
	}

	async function loadBootDataVolumes() {
		try {
			if (!request.namespace) return;

			const response = await virtualMachineClient.listDataVolumes({
				scope: scope,
				,
				namespace: request.namespace,
				bootImage: true
			});

			const dvOptions: BootDataVolumesOption[] = response.dataVolumes.map((dv: DataVolume) => ({
				value: dv.name,
				label: dv.name,
				icon: 'ph:hard-drive',
				sizeBytes: dv.sizeBytes
			}));

			bootDataVolumes.set(dvOptions);
		} catch (error) {
			toast.error('Failed to load bootable PVCs', {
				description: (error as ConnectError).message.toString()
			});
		}
	}

	// ==================== Default Values & Constants ====================
	const DEFAULT_REQUEST = {
		scope: scope,
		,
		name: '',
		namespace: 'default',
		instanceTypeName: '',
		bootDataVolumeName: '',
		startupScript: ''
	} as CreateVirtualMachineRequest;

	// ==================== Form State ====================
	let request: CreateVirtualMachineRequest = $state({ ...DEFAULT_REQUEST });

	// Load bootable PVCs when namespace changes
	$effect(() => {
		if (request.namespace) {
			loadBootDataVolumes();
		}
	});

	// ==================== Utility Functions ====================
	function reset() {
		request = { ...DEFAULT_REQUEST };
		isAdvancedOpen = false;
	}
	function close() {
		open = false;
	}

	// ==================== Lifecycle Hooks ====================
	onMount(() => {
		loadInstanceTypes();
		loadBootDataVolumes();
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
					<SingleInput.General
						required
						type="text"
						bind:value={request.name}
						bind:invalid={invalidName}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.namespace} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.instance_name()}</Form.Label>
					<SingleSelect.Root
						required
						options={instanceTypes}
						bind:value={request.instanceTypeName}
						bind:invalid={invalidInstanceTypeName}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $instanceTypes as instanceType}
											<SingleSelect.Item option={instanceType}>
												<Icon
													icon={instanceType.icon ? instanceType.icon : 'ph:empty'}
													class={cn('size-5', instanceType.icon ? 'visible' : 'invisible')}
												/>
												{instanceType.label}
												<SingleSelect.Check option={instanceType} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.data_volume()}</Form.Label>
					<SingleSelect.Root
						options={bootDataVolumes}
						required
						bind:value={request.bootDataVolumeName}
						bind:invalid={invalidBootDataVolumeName}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $bootDataVolumes as dv}
											<SingleSelect.Item option={dv}>
												<Icon
													icon={dv.icon ? dv.icon : 'ph:empty'}
													class={cn('size-5', dv.icon ? 'visible' : 'invisible')}
												/>
												{dv.label}
												<SingleSelect.Check option={dv} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>
			</Form.Fieldset>
			<!-- ==================== Advanced Configuration ==================== -->
			<Collapsible.Root bind:open={isAdvancedOpen} class="py-4">
				<div class="flex items-center justify-between gap-2">
					<p class={cn('text-base font-bold', isAdvancedOpen ? 'invisible' : 'visible')}>
						{m.advance()}
					</p>
					<Collapsible.Trigger class="rounded-full bg-muted p-1 ">
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
							<Form.Label>{m.startup_script()}</Form.Label>
							<Code.Root lang="bash" class="w-full" hideLines code={request.startupScript}>
								<Code.CopyButton />
							</Code.Root>
							<SingleInput.Structure
								preview={false}
								bind:value={request.startupScript}
								language="bash"
							/>
							<div class="flex justify-end gap-2">
								<Button
									variant="outline"
									size="sm"
									href="https://cloudinit.readthedocs.io/en/latest/reference/examples.html"
									target="_blank"
									class="flex items-center gap-1"
								>
									<Icon icon="ph:arrow-square-out" />
									{m.reference()}
								</Button>
							</div>
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
					disabled={invalidName}
					onclick={() => {
						toast.promise(() => virtualMachineClient.createVirtualMachine(request), {
							loading: `Creating ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to create ${request.name}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
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
