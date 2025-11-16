<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { DataVolume, ExtendDataVolumeRequest } from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	// Component props - accepts a DataVolume object
	let { dataVolume, scope, reloadManager }: { dataVolume: DataVolume; scope: string; reloadManager: ReloadManager  } = $props();

	// Get required services from Svelte context
	const transport: Transport = getContext('transport');

	// Create gRPC client for virtual machine operations
	const virtualMachineClient = createClient(InstanceService, transport);

	// Form validation state
	let invalid = $state(false);

	// Default values for the extend data volume request
	const defaults = {
		scope: scope,
		name: dataVolume.name,
		namespace: dataVolume.namespace,
		sizeBytes: dataVolume.sizeBytes
	} as ExtendDataVolumeRequest;

	// Current request state
	let request = $state({ ...defaults });

	// Reset form to default values
	function reset() {
		request = { ...defaults };
	}

	// Modal open/close state
	let open = $state(false);

	// Close modal function
	function close() {
		open = false;
	}
</script>

<!-- Modal component for data volume extension -->
<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
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
						toast.promise(() => virtualMachineClient.extendDataVolume(request), {
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
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
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
