<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type DeleteManifestRequest,
		type Manifest,
		RegistryService
	} from '$lib/api/registry/v1/registry_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		manifest,
		scope,
		reloadManager
	}: {
		manifest: Manifest;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');
	const registryClient = createClient(RegistryService, transport);

	let invalid: boolean | undefined = $state();

	const defaults = { scope, digest: manifest.digest } as DeleteManifestRequest;
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
		{m.remove()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.remove()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.repository()}</Form.Label>
					<SingleInput.Confirm
						required
						target={manifest.repositoryName}
						bind:value={request.repositoryName}
						bind:invalid
					/>
					<Form.Help>
						{m.deletion_warning({ identifier: m.repository() })}
					</Form.Help>
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
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => registryClient.deleteManifest(request), {
							loading: 'Loading...',
							success: () => {
								reloadManager.force();
								return `Delete ${manifest.repositoryName} success`;
							},
							error: (error) => {
								let message = `Fail to delete ${manifest.repositoryName}`;
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
