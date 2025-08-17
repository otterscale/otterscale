<script lang="ts" module>
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as SingleStepModal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
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

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalid: boolean | undefined = $state();

	const machineClient = createClient(MachineService, transport);
	const stateController = new StateController(false);
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
					<SingleInput.Confirm id="power-off" required target={machine.fqdn} bind:invalid />
					<Form.Help>
						Please type the machine fqdn {machine.fqdn} exactly to confirm.
					</Form.Help>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<SingleStepModal.Footer>
			<SingleStepModal.Cancel>Cancel</SingleStepModal.Cancel>
			<SingleStepModal.ActionsGroup>
				<SingleStepModal.Action
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

						stateController.close();
					}}
				>
					Confirm
				</SingleStepModal.Action>
			</SingleStepModal.ActionsGroup>
		</SingleStepModal.Footer>
	</SingleStepModal.Content>
</SingleStepModal.Root>
