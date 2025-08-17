<script lang="ts" module>
	import {
		NetworkService,
		type Network_IPRange,
		type UpdateIPRangeRequest
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
		ipRange
	}: {
		ipRange: Network_IPRange;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalidStartIP: boolean | undefined = $state();
	let invalidEndIP: boolean | undefined = $state();

	const client = createClient(NetworkService, transport);
	const requestManager = new RequestManager<UpdateIPRangeRequest>({
		id: ipRange.id,
		startIp: ipRange.startIp,
		endIp: ipRange.endIp,
		comment: ipRange.comment
	} as UpdateIPRangeRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		Edit
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit IP Range</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Start</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={requestManager.request.startIp}
						bind:invalid={invalidStartIP}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>End</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={requestManager.request.endIp}
						bind:invalid={invalidEndIP}
					/>
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
				disabled={invalidStartIP || invalidEndIP}
				onclick={() => {
					toast.promise(() => client.updateIPRange(requestManager.request), {
						loading: 'Loading...',
						success: () => {
							reloadManager.force();
							return `Update ${requestManager.request.startIp} - ${requestManager.request.endIp} success`;
						},
						error: (error) => {
							let message = `Fail to update ${requestManager.request.startIp} - ${requestManager.request.endIp}`;
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
