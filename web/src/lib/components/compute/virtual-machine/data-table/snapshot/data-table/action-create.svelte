<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		CreateVirtualMachineSnapshotRequest,
		VirtualMachine
	} from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';

	// Props
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const virtualMachineClient = createClient(InstanceService, transport);

	// ==================== State Variables ====================

	// UI state
	let open = $state(false);

	// Form validation state
	let invalidSnapshotName: boolean | undefined = $state();

	// ==================== Default Values & Constants ====================

	// Default request structure for creating a virtual machine snapshot
	const DEFAULT_REQUEST = {
		scope: $currentKubernetes?.scope,
		facility: $currentKubernetes?.name,
		namespace: virtualMachine.namespace,
		name: '',
		virtualMachineName: virtualMachine.name
	} as CreateVirtualMachineSnapshotRequest;

	// ==================== Form State ====================
	let request: CreateVirtualMachineSnapshotRequest = $state({ ...DEFAULT_REQUEST });

	// ==================== Utility Functions ====================
	function reset() {
		request = { ...DEFAULT_REQUEST };
	}
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_snapshot()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Basic Configuration ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.name}
						bind:invalid={invalidSnapshotName}
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
					disabled={invalidSnapshotName || !request.name}
					onclick={() => {
						toast.promise(() => virtualMachineClient.createVirtualMachineSnapshot(request), {
							loading: `Creating snapshot ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created snapshot ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to create snapshot ${request.name}`;
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
