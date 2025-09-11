<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { VirtualMachine, CloneVirtualMachineRequest } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { KubeVirtService } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const KubeVirtClient = createClient(KubeVirtService, transport);
	let invalid = $state(false);

	const defaults = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		targetName: '',
		targetNamespace: virtualMachine.metadata?.namespace,
		sourceName: virtualMachine.metadata?.name,
		sourceNamespace: virtualMachine.metadata?.namespace,
	} as CloneVirtualMachineRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
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
					<Form.Label>{m.name()}</Form.Label>
					<Form.Help>{m.clone_virtual_machine_name_direction()}</Form.Help>
					<SingleInput.General required type="text" bind:value={request.targetName} bind:invalid />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>
					<Form.Help>{m.clone_virtual_machine_namespace_direction()}</Form.Help>
					<SingleInput.General required type="text" bind:value={request.targetNamespace} bind:invalid />
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
						toast.promise(() => KubeVirtClient.cloneVirtualMachine(request), {
							loading: `Cloning ${request.sourceName} to ${request.targetName}...`,
							success: () => {
								reloadManager.force();
								return `Successfully cloned to ${request.targetName}`;
							},
							error: (error) => {
								let message = `Failed to clone ${request.sourceName}`;
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
