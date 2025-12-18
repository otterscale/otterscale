<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type DeleteModelArtifactRequest,
		type ModelArtifact,
		ModelService
	} from '$lib/api/model/v1/model_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		modelArtifact,
		scope,
		reloadManager
	}: {
		modelArtifact: ModelArtifact;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');
	const modelClient = createClient(ModelService, transport);

	let request = $state({} as DeleteModelArtifactRequest);
	let invalid = $state(false);
	let open = $state(false);

	function init() {
		request = {
			scope: scope,
			namespace: modelArtifact.namespace
		} as DeleteModelArtifactRequest;
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
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.model_artifact()}</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: m.model_artifact() })}
					</Form.Help>
					<SingleInput.Confirm
						required
						target={modelArtifact.name}
						bind:value={request.name}
						bind:invalid
					/>
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
					toast.promise(() => modelClient.deleteModelArtifact(request), {
						loading: `Deleting model artifact ${modelArtifact.name}...`,
						success: () => {
							reloadManager.force();
							return `Successfully deleted model artifact ${modelArtifact.name}`;
						},
						error: (error) => {
							let message = `Failed to delete model artifact ${modelArtifact.name}`;
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
