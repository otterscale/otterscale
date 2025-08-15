<script lang="ts" module>
	import {
		MachineService,
		type DeleteMachineRequest,
		type Machine
	} from '$lib/api/machine/v1/machine_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
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
	const reloadManager: ReloadManager = getContext('ReloadManager');

	let invalid: boolean | undefined = $state();

	const machineClient = createClient(MachineService, transport);
	const requestManager = new RequestManager<DeleteMachineRequest>({
		id: machine.id,
		force: false,
		purgeDisk: false
	} as DeleteMachineRequest);
	const stateController = new StateController(false);
</script>

<SingleStepModal.Root bind:open={stateController.state}>
	<SingleStepModal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		Remove
	</SingleStepModal.Trigger>
	<SingleStepModal.Content>
		<SingleStepModal.Header>Remove Machine</SingleStepModal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>FQDN</Form.Label>
					<SingleInput.Confirm required target={machine.fqdn} bind:invalid />
					<Form.Help>
						Please type the machine fqdn {machine.fqdn} exactly to confirm deletion. This action cannot
						be undone.
					</Form.Help>

					<Form.Field>
						<SingleInput.Boolean
							required
							descriptor={() => 'Force'}
							bind:value={requestManager.request.force}
						/>
					</Form.Field>

					<Form.Field>
						<SingleInput.Boolean
							required
							descriptor={() => 'Purge Disk'}
							bind:value={requestManager.request.purgeDisk}
						/>
					</Form.Field>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<SingleStepModal.Footer>
			<SingleStepModal.Cancel
				onclick={() => {
					requestManager.reset();
				}}>Cancel</SingleStepModal.Cancel
			>
			<SingleStepModal.ActionsGroup>
				<SingleStepModal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => machineClient.deleteMachine(requestManager.request), {
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

						requestManager.reset();
						stateController.close();
					}}
				>
					Remove
				</SingleStepModal.Action>
			</SingleStepModal.ActionsGroup>
		</SingleStepModal.Footer>
	</SingleStepModal.Content>
</SingleStepModal.Root>
