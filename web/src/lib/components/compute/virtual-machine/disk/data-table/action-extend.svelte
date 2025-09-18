<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		ExtendDataVolumeRequest,
		VirtualMachineDisk,
		DataVolumeSource,
	} from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { VirtualMachineDisk_type, KubeVirtService } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	// Component props - accepts a virtual machine disk object
	let { virtualMachineDisk }: { virtualMachineDisk: VirtualMachineDisk } = $props();

	// Get required services from Svelte context
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	// Create gRPC client for KubeVirt operations
	const KubeVirtClient = createClient(KubeVirtService, transport);

	// Form validation state
	let invalid = $state(false);

	// Default values for the extend data volume request
	const defaults = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		name: virtualMachineDisk.name,
		namespace: '',
		sizeBytes:
			virtualMachineDisk.sourceData.case === 'dataVolume'
				? (virtualMachineDisk.sourceData.value as DataVolumeSource).sizeBytes
				: undefined,
	} as ExtendDataVolumeRequest;

	// Current request state
	let request = $state(defaults);

	// Reset form to default values
	function reset() {
		request = defaults;
	}

	// Modal open/close state
	let open = $state(false);

	// Close modal function
	function close() {
		open = false;
	}
</script>

<!-- Modal component for disk extension -->
<Modal.Root bind:open>
	<Modal.Trigger variant="creative" disabled={virtualMachineDisk.diskType !== VirtualMachineDisk_type.DATAVOLUME}>
		<Icon icon="ph:arrows-out" />
		{m.extend()}
	</Modal.Trigger>

	<!-- Modal content -->
	<Modal.Content>
		<Modal.Header>{m.extend_data_volume()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.size()}</Form.Label>
					<SingleInput.Measurement
						required
						bind:value={request.sizeBytes}
						bind:invalid
						transformer={(value) => String(value)}
						units={[{ value: 1024 * 1024 * 1024, label: 'GB' } as SingleInput.UnitType]}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>

		<!-- Modal footer with action buttons -->
		<Modal.Footer>
			<!-- Cancel button -->
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>

			<!-- Confirm action group -->
			<Modal.ActionsGroup>
				<!-- Confirm button with extend operation -->
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						// Execute extend operation with toast notifications
						toast.promise(() => KubeVirtClient.extendDataVolume(request), {
							loading: `Extending ${request.name} to ${Math.floor(Number(request.sizeBytes) / (1024 * 1024 * 1024))}GB...`,
							success: () => {
								// Force reload to refresh data
								reloadManager.force();
								return `Successfully extended ${request.name}`;
							},
							error: (error) => {
								// Handle and display error
								let message = `Failed to extend ${request.name}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
								});
								return message;
							},
						});
						// Reset form and close modal
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
