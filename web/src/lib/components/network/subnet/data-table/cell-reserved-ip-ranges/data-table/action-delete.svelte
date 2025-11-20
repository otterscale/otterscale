<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type DeleteIPRangeRequest,
		type Network_IPRange,
		NetworkService
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
	const defaults = {
		id: ipRange.id
	} as DeleteIPRangeRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let invalidity = $state({} as Booleanified<Network_IPRange>);
	const invalid = $derived(invalidity.startIp || invalidity.endIp);

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		{m.delete()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete_ip_range()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.start_ip()}</Form.Label>
					<SingleInput.Confirm
						required
						target={ipRange.startIp}
						bind:invalid={invalidity.startIp}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.end_ip()}</Form.Label>
					<SingleInput.Confirm required target={ipRange.endIp} bind:invalid={invalidity.endIp} />
				</Form.Field>
				<Form.Help>
					{m.deletion_warning({ identifier: m.range() })}
				</Form.Help>
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
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => client.deleteIPRange(request), {
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

						reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
