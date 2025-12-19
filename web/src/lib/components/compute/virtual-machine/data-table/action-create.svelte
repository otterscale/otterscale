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
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let { scope, reloadManager }: { scope: string; reloadManager: ReloadManager } = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(InstanceService, transport);

	// ==================== State Variables ====================

	// UI state
	const DEFAULT_NAMESPACE = 'kubevirt';
	let isAdvancedOpen = $state(false);
	let open = $state(false);

	// Form validation state
	let invalidity = $state({} as Booleanified<CreateVirtualMachineRequest>);
	const invalid = $derived(
		invalidity.name || invalidity.instanceTypeName || invalidity.bootDataVolumeName
	);

	// Instance type with CPU and memory information
	type InstanceTypeOption = SingleSelect.OptionType & {
		cpuCores?: number;
		memoryBytes?: bigint;
	};

	type BootDataVolumesOption = SingleSelect.OptionType & {
		sizeBytes?: bigint;
		phase: string;
	};

	// ==================== Local Dropdown Options ====================
	const bootDataVolumes: Writable<BootDataVolumesOption[]> = writable([]);
	const instanceTypes: Writable<InstanceTypeOption[]> = writable([]);

	// ==================== API Functions ====================
	async function loadInstanceTypes() {
		try {
			// Request both namespace-specific and cluster-wide instance types
			const response = await virtualMachineClient.listInstanceTypes({
				scope: scope,
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
			const response = await virtualMachineClient.listDataVolumes({
				scope: scope,
				namespace: DEFAULT_NAMESPACE,
				bootImage: true
			});

			const dvOptions: BootDataVolumesOption[] = response.dataVolumes
				.filter((dv) => dv.phase !== 'Failed')
				.map((dv: DataVolume) => ({
					value: dv.name,
					label: dv.name,
					icon: 'ph:hard-drive',
					phase: dv.phase,
					sizeBytes: dv.sizeBytes
				}));

			bootDataVolumes.set(dvOptions);
		} catch (error) {
			toast.error('Failed to load bootable PVCs', {
				description: (error as ConnectError).message.toString()
			});
		}
	}

	// ==================== Form State ====================
	let request = $state({} as CreateVirtualMachineRequest);

	function init() {
		request = {
			scope: scope,
			name: '',
			namespace: DEFAULT_NAMESPACE,
			instanceTypeName: '',
			bootDataVolumeName: '',
			startupScript: ''
		} as CreateVirtualMachineRequest;
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

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
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
					<SingleInput.GeneralRule
						required
						type="text"
						bind:value={request.name}
						bind:invalid={invalidity.name}
						validateRule="rfc1123"
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>
					<SingleInput.General disabled type="text" bind:value={request.namespace} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.instance_name()}</Form.Label>
					<SingleSelect.Root
						required
						options={instanceTypes}
						bind:value={request.instanceTypeName}
						bind:invalid={invalidity.instanceTypeName}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $instanceTypes as instanceType (instanceType.value)}
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
						bind:invalid={invalidity.bootDataVolumeName}
						onOpenChange={(open) => {
							if (open) loadBootDataVolumes();
						}}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $bootDataVolumes as dv (dv.value)}
											<SingleSelect.Item option={dv} disabled={dv.phase !== 'Succeeded'}>
												{#if dv.phase === 'Succeeded'}
													<Icon icon="ph:hard-drive" class="size-5" />
												{:else}
													<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
												{/if}
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
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
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
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
