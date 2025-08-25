<script lang="ts" module>
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Machine;
	} = $props();

	const reloadManager: ReloadManager = getContext('reloadManager');

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	let invalid: boolean | undefined = $state();

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger disabled={machine.powerState.toLowerCase() === 'off'} variant="creative">
		<Icon icon="ph:power" />
		{m.turn_off()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.turn_off_machine()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.fqdn()}</Form.Label>
					<SingleInput.Confirm required target={machine.fqdn} bind:invalid />
					<Form.Help>
						{m.deletion_warning({ identifier: m.fqdn() })}
					</Form.Help>
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
						toast.promise(
							() =>
								machineClient.powerOffMachine({
									id: machine.id
								}),
							{
								loading: 'Loading...',
								success: () => {
									reloadManager.force();
									return `Turn off ${machine.fqdn}`;
								},
								error: (error) => {
									let message = `Fail to turn off ${machine.fqdn}`;
									toast.error(message, {
										description: (error as ConnectError).message.toString(),
										duration: Number.POSITIVE_INFINITY
									});
									return message;
								}
							}
						);
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
