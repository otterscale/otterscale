<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		ApplicationService,
		type Application_Release,
		type RollbackReleaseRequest,
	} from '$lib/api/application/v1/application_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		release,
		releases = $bindable(),
	}: {
		release: Application_Release;
		releases: Writable<Application_Release[]>;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	const defaults = {
		scopeUuid: release.name,
		facilityName: release.name,
		namespace: release.namespace,
	} as RollbackReleaseRequest;
	let request = $state(defaults);
	function reset() {
		request = { dryRun: false } as RollbackReleaseRequest;
	}

	let invalid = $state(false);

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:arrow-counter-clockwise" />
		{m.rollback()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>
			{m.rollback_release()}
		</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Help>
					{m.deletion_warning({ identifier: m.name() })}
				</Form.Help>
				<Form.Field>
					<SingleInput.Confirm required id="deletion" target={release.name} bind:value={request.name} />
				</Form.Field>
				<Form.Field>
					<SingleInput.Boolean descriptor={() => m.dry_run()} bind:value={request.dryRun} />
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
					toast.promise(() => client.rollbackRelease(request), {
						loading: 'Loading...',
						success: () => {
							client.listReleases({}).then((r) => {
								releases.set(r.releases);
							});
							return `Rollback ${request.name}`;
						},
						error: (e) => {
							let msg = `Fail to rollback ${request.name}`;
							toast.error(msg, {
								description: (e as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY,
							});
							return msg;
						},
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
