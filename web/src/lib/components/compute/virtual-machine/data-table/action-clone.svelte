<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		CreateVirtualMachineCloneRequest,
		VirtualMachine
	} from '$lib/api/instance/v1/instance_pb';
	import { InstanceService } from '$lib/api/instance/v1/instance_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { virtualMachine, scope, reloadManager }: { virtualMachine: VirtualMachine; scope: string; reloadManager: ReloadManager } = $props();

	const transport: Transport = getContext('transport');

	const virtualMachineClient = createClient(InstanceService, transport);
	let invalid = $state(false);

	const defaults = {
		scope: scope,
		namespace: virtualMachine.namespace,
		name: `clone-${virtualMachine.name}-${new Date().toISOString().slice(0, 10).replace(/-/g, '')}`,
		sourceVirtualMachineName: virtualMachine.name,
		targetVirtualMachineName: ''
	} as CreateVirtualMachineCloneRequest;
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
	<Modal.Trigger variant="creative">
		<Icon icon="ph:copy" />
		{m.clone()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.clone_virtual_machine()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.namespace} bind:invalid />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.name} bind:invalid />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.source()}</Form.Label>
					<Form.Help>{m.source_virtual_machine_name()}</Form.Help>
					<SingleInput.General type="text" value={request.sourceVirtualMachineName} disabled />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.target()}</Form.Label>
					<Form.Help>{m.target_virtual_machine_name()}</Form.Help>
					<SingleInput.General
						required
						type="text"
						bind:value={request.targetVirtualMachineName}
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
						toast.promise(() => virtualMachineClient.createVirtualMachineClone(request), {
							loading: `Cloning ${request.sourceVirtualMachineName} to ${request.targetVirtualMachineName}...`,
							success: () => {
								reloadManager.force();
								return `Successfully cloned to ${request.targetVirtualMachineName}`;
							},
							error: (error) => {
								let message = `Failed to clone ${request.targetVirtualMachineName}`;
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
