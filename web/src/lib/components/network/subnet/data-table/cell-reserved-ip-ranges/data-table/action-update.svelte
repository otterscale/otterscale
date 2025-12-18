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
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { ipRange, reloadManager }: { ipRange: Network_IPRange; reloadManager: ReloadManager } =
		$props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	let request = $state({} as UpdateIPRangeRequest);
	let open = $state(false);

	function init() {
		request = {
			id: ipRange.id,
			startIp: ipRange.startIp,
			endIp: ipRange.endIp,
			comment: ipRange.comment
		} as UpdateIPRangeRequest;
	}

	function close() {
		open = false;
	}

	let invalidity = $state({} as Booleanified<Network_IPRange>);
	const invalid = $derived(invalidity.startIp || invalidity.endIp);
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
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
						bind:invalid={invalidity.startIp}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.end_ip()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.endIp}
						bind:invalid={invalidity.endIp}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.comment()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.comment} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
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
					close();
				}}
			>
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
