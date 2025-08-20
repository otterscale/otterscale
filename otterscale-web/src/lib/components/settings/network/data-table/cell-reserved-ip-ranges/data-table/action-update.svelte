<script lang="ts">
	import {
		NetworkService,
		type Network_IPRange,
		type UpdateIPRangeRequest
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

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
	const defaults = {
		id: ipRange.id,
		startIp: ipRange.startIp,
		endIp: ipRange.endIp,
		comment: ipRange.comment
	} as UpdateIPRangeRequest;
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
						bind:value={request.startIp}
						bind:invalid={invalidStartIP}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>End</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.endIp}
						bind:invalid={invalidEndIP}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Coment</Form.Label>
					<SingleInput.General type="text" bind:value={request.comment} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				Cancel
			</Modal.Cancel>
			<Modal.Action
				disabled={invalidStartIP || invalidEndIP}
				onclick={() => {
					toast.promise(() => client.updateIPRange(request), {
						loading: 'Loading...',
						success: () => {
							reloadManager.force();
							return `Update ${request.startIp} - ${request.endIp} success`;
						},
						error: (error) => {
							let message = `Fail to update ${request.startIp} - ${request.endIp}`;
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
				Create
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
