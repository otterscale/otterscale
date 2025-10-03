<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { VirtualMachineService } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import type {
		VirtualMachine_Restore,
		DeleteVirtualMachineRestoreRequest,
	} from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { virtualMachineRestore }: { virtualMachineRestore: VirtualMachine_Restore } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const virtualMachineClient = createClient(VirtualMachineService, transport);
	let invalid = $state(false);

	const defaults = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		name: '',
		namespace: virtualMachineRestore.namespace,
	} as DeleteVirtualMachineRestoreRequest;
	let request = $state({ ...defaults });
	function reset() {
		request = { ...defaults };
	}

	let open = $state(false);
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
		<Modal.Header>{m.delete_restore()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.restore_name()}</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: m.restore_name() })}
					</Form.Help>
					<SingleInput.Confirm
						required
						target={virtualMachineRestore.name ?? ''}
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
						toast.promise(() => virtualMachineClient.deleteVirtualMachineRestore(request), {
							loading: `Deleting ${virtualMachineRestore.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully deleted ${virtualMachineRestore.name}`;
							},
							error: (error) => {
								let message = `Failed to delete ${virtualMachineRestore.name}`;
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
