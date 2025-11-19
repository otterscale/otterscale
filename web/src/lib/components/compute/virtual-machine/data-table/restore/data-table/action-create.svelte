<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		CreateVirtualMachineRestoreRequest,
		VirtualMachine
	} from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	// Props
	let {
		virtualMachine,
		scope,
		reloadManager
	}: { virtualMachine: VirtualMachine; scope: string; reloadManager: ReloadManager } = $props();

	// Context dependencies
	const transport: Transport = getContext('transport');

	const virtualMachineClient = createClient(InstanceService, transport);

	// ==================== State Variables ====================

	// UI state
	let open = $state(false);

	// Form validation state
	// TODO: Refactor with Booneanified type of Request
	let invalidRestoreName: boolean | undefined = $state();

	// ==================== Default Values & Constants ====================

	// Default request structure for creating a virtual machine restore
	const DEFAULT_REQUEST = {
		scope: scope,
		namespace: virtualMachine.namespace,
		name: '',
		virtualMachineName: virtualMachine.name
	} as CreateVirtualMachineRestoreRequest;

	// ==================== Form State ====================
	let request: CreateVirtualMachineRestoreRequest = $state({ ...DEFAULT_REQUEST });

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
		<Modal.Header>{m.create_restore()}</Modal.Header>
		<Form.Root>
			<!-- ==================== Basic Configuration ==================== -->
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.name}
						bind:invalid={invalidRestoreName}
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
					disabled={invalidRestoreName || !request.name}
					onclick={() => {
						toast.promise(() => virtualMachineClient.createVirtualMachineRestore(request), {
							loading: `Creating restore ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created restore ${request.name}`;
							},
							error: (error) => {
								let message = `Failed to create restore ${request.name}`;
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
