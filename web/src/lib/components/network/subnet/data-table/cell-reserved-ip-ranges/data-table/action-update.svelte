<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type Network_IPRange,
		NetworkService,
		type UpdateIPRangeRequest
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
</script>

<script lang="ts">
	let { ipRange, reloadManager }: { ipRange: Network_IPRange; reloadManager: ReloadManager } =
		$props();

	const transport: Transport = getContext('transport');

	let invalidities = $state({} as Booleanified<Network_IPRange>);
	const invalid = $derived(invalidities.startIp || invalidities.endIp);

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
		{m.edit()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit_ip_range()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.start_ip()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.startIp}
						bind:invalid={invalidities.startIp}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.end_ip()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.endIp}
						bind:invalid={invalidities.endIp}
					/>
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
				disabled={invalid}
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
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
