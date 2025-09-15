<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { SnapshotVirtualMachineRequest, VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { KubeVirtService } from '$lib/api/kubevirt/v1/kubevirt_pb';
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
	const kubevirtClient = createClient(KubeVirtService, transport);

	// ==================== State Variables ====================

	// UI state
	let open = $state(false);

	// Form validation state
	let invalidSnapshotName: boolean | undefined = $state();

	// ==================== Default Values & Constants ====================

	// Default request structure for creating a virtual machine snapshot
	const DEFAULT_REQUEST = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		name: virtualMachine.metadata?.name,
		namespace: virtualMachine.metadata?.namespace,
		snapshotName: '',
		description: '',
	} as SnapshotVirtualMachineRequest;

	// ==================== Form State ====================
	let request: SnapshotVirtualMachineRequest = $state(DEFAULT_REQUEST);

	// ==================== Utility Functions ====================
	function reset() {
		request = DEFAULT_REQUEST;
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
						bind:value={request.snapshotName}
						bind:invalid={invalidSnapshotName}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.description()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.description} />
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
					disabled={invalidSnapshotName || !request.snapshotName}
					onclick={() => {
						toast.promise(() => kubevirtClient.snapshotVirtualMachine(request), {
							loading: `Creating snapshot ${request.snapshotName}...`,
							success: () => {
								reloadManager.force();
								return `Successfully created snapshot ${request.snapshotName}`;
							},
							error: (error) => {
								let message = `Failed to create snapshot ${request.snapshotName}`;
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
