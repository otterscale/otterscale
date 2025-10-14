<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { DeleteDataVolumeRequest, DataVolume } from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	// Component props - accepts a virtual machine disk object
	let { dataVolume }: { dataVolume: DataVolume } = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const virtualMachineClient = createClient(InstanceService, transport);

	// Form validation state
	let invalid = $state(false);

	// Default values for the delete data volume request
	const defaults = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		namespace: dataVolume.namespace,
		name: '',
	} as DeleteDataVolumeRequest;

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
		<Icon icon="ph:trash" />
		{m.delete()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete()} {m.data_volume()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.data_volume()}</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: m.data_volume_name() })}
					</Form.Help>
					{console.log(defaults)}
					<SingleInput.Confirm
						required
						target={dataVolume.name ?? ''}
						bind:value={request.name}
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
						toast.promise(() => virtualMachineClient.deleteDataVolume(request), {
							loading: `Deleting data volume ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully deleted data volume ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to delete data volume ${request.name}`;
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
