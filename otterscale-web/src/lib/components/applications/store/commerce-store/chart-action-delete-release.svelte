<script lang="ts">
	import {
		ApplicationService,
		type Application_Release,
		type DeleteReleaseRequest
	} from '$lib/api/application/v1/application_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';

	let {
		release,
		releases = $bindable()
	}: {
		release: Application_Release;
		releases: Writable<Application_Release[]>;
	} = $props();

	const transport: Transport = getContext('transport');

	const client = createClient(ApplicationService, transport);

	const defaults = {
		dryRun: false,
		scopeUuid: release.name,
		facilityName: release.name,
		namespace: release.namespace
	} as DeleteReleaseRequest;
	let request = $state(defaults as DeleteReleaseRequest);
	function reset() {
		request = { dryRun: false } as DeleteReleaseRequest;
	}

	let invalid = $state(false);

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		Delete
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Delete Release</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Help>
					Please type the release name exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
				<Form.Field>
					<SingleInput.Confirm
						required
						id="deletion"
						target={release.name}
						bind:value={request.name}
					/>
				</Form.Field>
				<Form.Field>
					<SingleInput.Boolean descriptor={() => 'Dry Run'} bind:value={request.dryRun} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					toast.promise(() => client.deleteRelease(request), {
						loading: 'Loading...',
						success: (r) => {
							client.listReleases({}).then((r) => {
								releases.set(r.releases);
							});
							return `Delete ${request.name}`;
						},
						error: (e) => {
							let msg = `Fail to delete ${request.name}`;
							toast.error(msg, {
								description: (e as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return msg;
						}
					});

					close();
				}}>Confirm</Modal.Action
			>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
