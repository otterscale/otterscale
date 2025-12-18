<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type DeleteMachineRequest,
		type Machine,
		MachineService
	} from '$lib/api/machine/v1/machine_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		machine,
		reloadManager
	}: {
		machine: Machine;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	let request = $state({} as DeleteMachineRequest);
	let invalid: boolean | undefined = $state();
	let open = $state(false);

	function init() {
		request = {
			id: machine.id,
			force: false,
			purgeDisk: false
		} as DeleteMachineRequest;
	}

	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		{m.remove()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.remove_machine()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.fqdn()}</Form.Label>
					<SingleInput.Confirm required target={machine.fqdn} bind:invalid />
					<Form.Help>
						{m.deletion_warning({ identifier: m.fqdn() })}
					</Form.Help>
					<Form.Field>
						<SingleInput.Boolean descriptor={() => m.force()} bind:value={request.force} />
					</Form.Field>
					<Form.Field>
						<SingleInput.Boolean descriptor={() => m.purge_disk()} bind:value={request.purgeDisk} />
					</Form.Field>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => machineClient.deleteMachine(request), {
							loading: 'Loading...',
							success: () => {
								reloadManager.force();
								return `Delete ${machine.fqdn} success`;
							},
							error: (error) => {
								let message = `Fail to delete ${machine.fqdn}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
