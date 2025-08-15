<script lang="ts" module>
	import {
		NetworkService,
		type CreateIPRangeRequest,
		type Network_Subnet
	} from '$lib/api/network/v1/network_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
</script>

<script lang="ts">
	let {
		subnet
	}: {
		subnet: Network_Subnet;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('ReloadManager');

	const client = createClient(NetworkService, transport);
	const requestManager = new RequestManager<CreateIPRangeRequest>({
		subnetId: subnet.id
	} as CreateIPRangeRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		Create
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create IP Range</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Start</Form.Label>
					<SingleInput.General type="text" bind:value={requestManager.request.startIp} />
				</Form.Field>

				<Form.Field>
					<Form.Label>End</Form.Label>
					<SingleInput.General type="text" bind:value={requestManager.request.endIp} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Coment</Form.Label>
					<SingleInput.General type="text" bind:value={requestManager.request.comment} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					requestManager.reset();
				}}
				class="mr-auto">Cancel</Modal.Cancel
			>
			<Modal.Action
				onclick={() => {
					toast.promise(() => client.createIPRange(requestManager.request), {
						loading: 'Loading...',
						success: () => {
							reloadManager.force();
							return `Create ${requestManager.request.startIp} - ${requestManager.request.endIp} success`;
						},
						error: (error) => {
							let message = `Fail to create ${requestManager.request.startIp} - ${requestManager.request.endIp}`;
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
				Create
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
