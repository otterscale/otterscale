<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { DetachVirtualMachineDiskRequest } from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	// Component props - accepts a virtual machine disk object
	let {
		enhancedDisk,
		scope,
		reloadManager
	}: {
		enhancedDisk: EnhancedDisk;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');

	const virtualMachineClient = createClient(InstanceService, transport);

	// Form validation state
	let invalid = $state(false);

	// Default values for the detach disk request
	const defaults = {
		scope: scope,
		namespace: enhancedDisk.namespace,
		name: enhancedDisk.vmName,
		dataVolumeName: ''
	} as DetachVirtualMachineDiskRequest;

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

<Modal.Root bind:open>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:plugs" />
		{m.detach()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.detach_disk()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.data_volume()}</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: m.data_volume_name() })}
					</Form.Help>
					<SingleInput.Confirm
						required
						target={enhancedDisk.name ?? ''}
						bind:value={request.dataVolumeName}
						bind:invalid
					/>
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
					disabled={invalid}
					onclick={() => {
						toast.promise(() => virtualMachineClient.detachVirtualMachineDisk(request), {
							loading: `Detaching disk ${enhancedDisk.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully detached disk ${enhancedDisk.name}`;
							},
							error: (error) => {
								let message = `Failed to detach disk ${enhancedDisk.name}`;
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
