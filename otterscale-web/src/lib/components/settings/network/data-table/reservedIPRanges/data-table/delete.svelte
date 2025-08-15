<script lang="ts" module>
	import {
		NetworkService,
		type DeleteIPRangeRequest,
		type Network_IPRange
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
	let { ipRange }: { ipRange: Network_IPRange } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('ReloadManager');

	let invalidStartIP: boolean | undefined = $state();
	let invalidEndIP: boolean | undefined = $state();

	const client = createClient(NetworkService, transport);
	const requestManager = new RequestManager<DeleteIPRangeRequest>({
		id: ipRange.id
	} as DeleteIPRangeRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		Delete
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Delete IP Range</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Start</Form.Label>
					<SingleInput.Confirm required target={ipRange.startIp} bind:invalid={invalidStartIP} />
				</Form.Field>
				<Form.Field>
					<Form.Label>End</Form.Label>
					<SingleInput.Confirm required target={ipRange.endIp} bind:invalid={invalidEndIP} />
				</Form.Field>
				<Form.Help>
					Please type the ip range exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					requestManager.reset();
				}}>Cancel</Modal.Cancel
			>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalidStartIP || invalidEndIP}
					onclick={() => {
						toast.promise(() => client.deleteIPRange(requestManager.request), {
							loading: 'Loading...',
							success: () => {
								reloadManager.force();
								return `Delete ${ipRange.id} success`;
							},
							error: (error) => {
								let message = `Fail to delete ${ipRange.id}`;
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
					Delete
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
