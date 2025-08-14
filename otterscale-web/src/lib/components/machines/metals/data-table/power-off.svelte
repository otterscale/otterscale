<script lang="ts" module>
	import {
		MachineService,
		type Machine,
		type PowerOffMachineRequest
	} from '$lib/api/machine/v1/machine_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as SingleStepModal } from '$lib/components/custom/modal';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	let {
		machine,
		machines = $bindable()
	}: {
		machine: Machine;
		machines: Writable<Machine[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		id: machine.id
	} as PowerOffMachineRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new StateController(false);

	let fqdn = $state();
	let invalid: boolean | undefined = $state();
</script>

<SingleStepModal.Root bind:open={stateController.state}>
	<SingleStepModal.Trigger disabled={machine.powerState.toLowerCase() === 'off'} variant="creative">
		<Icon icon="ph:power" />
		Power Off
	</SingleStepModal.Trigger>
	<SingleStepModal.Content>
		<SingleStepModal.Header>Turn Off Machine</SingleStepModal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>FQDN</Form.Label>
					<SingleInput.DeletionConfirm
						id="power-off"
						required
						target={machine.fqdn}
						bind:value={fqdn}
						bind:invalid
					/>
					<Form.Help>
						Please type the machine fqdn {machine.fqdn} exactly to confirm.
					</Form.Help>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<SingleStepModal.Footer>
			<SingleStepModal.Cancel onclick={reset}>Cancel</SingleStepModal.Cancel>
			<SingleStepModal.ActionsGroup>
				<SingleStepModal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => machineClient.powerOffMachine(request), {
							loading: 'Loading...',
							success: () => {
								machineClient.getMachine({ id: machine.id }).then((response) => {
									machine = response;
								});
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
						});

						reset();
						stateController.close();
					}}
				>
					Confirm
				</SingleStepModal.Action>
			</SingleStepModal.ActionsGroup>
		</SingleStepModal.Footer>
	</SingleStepModal.Content>
</SingleStepModal.Root>
