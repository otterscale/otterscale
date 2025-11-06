<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { ModelService, type DeleteModelArtifactRequest, type ModelArtifact } from '$lib/api/model/v1/model_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { modelArtifact }: { modelArtifact: ModelArtifact } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const modelClient = createClient(ModelService, transport);

	let invalid = $state(false);

	const defaults = {
		scope: $currentKubernetes?.scope,
		facility: $currentKubernetes?.name,
		namespace: 'default',
	} as DeleteModelArtifactRequest;
	let request = $state({ ...defaults });
	function reset() {
		request = { ...defaults };
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
		<Modal.Header>{m.delete()} {m.model()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.model()}</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: m.model() })}
					</Form.Help>
					<SingleInput.Confirm required target={modelArtifact.name} bind:value={request.name} bind:invalid />
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
					toast.promise(() => modelClient.deleteModelArtifact(request), {
						loading: `Deleting model artifact ${request.name}...`,
						success: () => {
							reloadManager.force();
							return `Successfully deleted model artifact ${request.name}`;
						},
						error: (error) => {
							let message = `Failed to delete model artifact ${request.name}`;
							toast.error(message, {
								description: (error as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY,
							});
							return message;
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
