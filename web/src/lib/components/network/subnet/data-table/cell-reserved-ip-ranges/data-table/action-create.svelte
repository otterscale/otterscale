<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type CreateIPRangeRequest,
		type Network_Subnet,
		NetworkService
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		subnet
	}: {
		subnet: Network_Subnet;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const client = createClient(NetworkService, transport);
	const defaults = {
		subnetId: subnet.id
	} as CreateIPRangeRequest;
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
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_ip_range()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.start_ip()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.startIp} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.end_ip()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.endIp} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.comment()}</Form.Label>
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
				{m.cancel()}
			</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					toast.promise(() => client.createIPRange(request), {
						loading: 'Loading...',
						success: () => {
							reloadManager.force();
							return `Create ${request.startIp} - ${request.endIp} success`;
						},
						error: (error) => {
							let message = `Fail to create ${request.startIp} - ${request.endIp}`;
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
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
