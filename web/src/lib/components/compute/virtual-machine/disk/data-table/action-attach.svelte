<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type {
		AttachVirtualMachineDiskRequest,
		DataVolume,
		VirtualMachine
	} from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';

	let {
		virtualMachine
	}: {
		virtualMachine: VirtualMachine;
	} = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const virtualMachineClient = createClient(InstanceService, transport);

	// ==================== State Variables ====================

	// UI state
	let open = $state(false);

	// Form validation state
	let invalidDataVolumnName: boolean | undefined = $state();

	// ==================== Local Dropdown Options ====================
	const dataVolumes: Writable<SingleSelect.OptionType[]> = writable([]);

	// ==================== API Functions ====================
	async function loadDataVolumes() {
		try {
			if (!request.namespace) return;

			const response = await virtualMachineClient.listDataVolumes({
				scope: $currentKubernetes?.scope,
				facility: $currentKubernetes?.name,
				namespace: request.namespace
			});

			const dvOptions: SingleSelect.OptionType[] = response.dataVolumes.map((dv: DataVolume) => ({
				value: dv.name,
				label: dv.name,
				icon: 'ph:hard-drive'
			}));

			dataVolumes.set(dvOptions);
		} catch (error) {
			toast.error('Failed to load Data Volumes', {
				description: (error as ConnectError).message.toString()
			});
		}
	}

	// ==================== Default Values & Constants ====================
	const DEFAULT_REQUEST = {
		scope: $currentKubernetes?.scope,
		facility: $currentKubernetes?.name,
		name: virtualMachine.name,
		namespace: virtualMachine.namespace,
		dataVolumeName: ''
	} as AttachVirtualMachineDiskRequest;

	// ==================== Form State ====================
	let request: AttachVirtualMachineDiskRequest = $state({ ...DEFAULT_REQUEST });

	// ==================== Utility Functions ====================
	function reset() {
		request = { ...DEFAULT_REQUEST };
	}

	function close() {
		open = false;
	}

	// ==================== Lifecycle Hooks ====================
	onMount(() => {
		loadDataVolumes();
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.attach()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.attach_disk()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Disk Information ==================== -->
			<Form.Fieldset>
				<!-- <Form.Legend>{m.information()}</Form.Legend> -->
				<Form.Field>
					<Form.Label>{m.data_volume()}</Form.Label>
					<SingleSelect.Root
						required
						options={dataVolumes}
						bind:value={request.dataVolumeName}
						bind:invalid={invalidDataVolumnName}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $dataVolumes as dv}
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
					disabled={invalidDataVolumnName}
					onclick={() => {
						toast.promise(() => virtualMachineClient.attachVirtualMachineDisk(request), {
							loading: `Attaching ${request.dataVolumeName} to ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully attached ${request.dataVolumeName} to ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to attach ${request.dataVolumeName} to ${request.name}`;
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
