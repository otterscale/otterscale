<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		ApplicationService,
		type DeleteReleaseRequest,
		type Release
	} from '$lib/api/application/v1/application_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		release,
		scope,
		releases
	}: {
		release: Release;
		scope: string;
		releases: Writable<Release[]>;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	let request = $state({} as DeleteReleaseRequest);
	let invalid = $state(false);
	let open = $state(false);

	function init() {
		request = {
			dryRun: false,
			scope: scope,
			namespace: release.namespace
		} as DeleteReleaseRequest;
	}

	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		{m.delete()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete_release()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Help>
					{m.deletion_warning({ identifier: m.name() })}
				</Form.Help>
				<Form.Field>
					<SingleInput.Confirm
						required
						id="deletion"
						target={release.name}
						bind:value={request.name}
						bind:invalid
					/>
				</Form.Field>
				<Form.Field>
					<SingleInput.Boolean descriptor={() => m.dry_run()} bind:value={request.dryRun} />
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
					toast.promise(() => client.deleteRelease(request), {
						loading: 'Loading...',
						success: () => {
							client.listReleases({ scope: scope }).then((r) => {
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
				}}
			>
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
