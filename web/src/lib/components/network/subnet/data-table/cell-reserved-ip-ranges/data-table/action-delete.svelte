<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		NetworkService,
		type DeleteIPRangeRequest,
		type Network_IPRange
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { ipRange }: { ipRange: Network_IPRange } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalidStartIP: boolean | undefined = $state();
	let invalidEndIP: boolean | undefined = $state();

	const client = createClient(NetworkService, transport);
	const defaults = {
		id: ipRange.id
	} as DeleteIPRangeRequest;
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
					<SingleInput.Confirm required target={ipRange.startIp} bind:invalid={invalidStartIP} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.end_ip()}</Form.Label>
					<SingleInput.Confirm required target={ipRange.endIp} bind:invalid={invalidEndIP} />
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
					disabled={invalidStartIP || invalidEndIP}
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
