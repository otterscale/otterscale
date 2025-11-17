<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		DeleteVirtualMachineSnapshotRequest,
		VirtualMachine_Snapshot
	} from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		virtualMachineSnapshot,
		scope,
		reloadManager
	}: {
		virtualMachineSnapshot: VirtualMachine_Snapshot;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');

	const virtualMachineClient = createClient(InstanceService, transport);
	let invalid = $state(false);

	const defaults = {
		scope: scope,
		name: '',
		namespace: virtualMachineSnapshot.namespace
	} as DeleteVirtualMachineSnapshotRequest;
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
		<Modal.Header>{m.delete_snapshot()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.snapshot_name()}</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: m.snapshot_name() })}
					</Form.Help>
					<SingleInput.Confirm
						required
						target={virtualMachineSnapshot.name ?? ''}
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
						toast.promise(() => virtualMachineClient.deleteVirtualMachineSnapshot(request), {
							loading: `Deleting ${virtualMachineSnapshot.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully deleted ${virtualMachineSnapshot.name}`;
							},
							error: (error) => {
								let message = `Failed to delete ${virtualMachineSnapshot.name}`;
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
