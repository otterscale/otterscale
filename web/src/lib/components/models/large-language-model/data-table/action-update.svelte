<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { LargeLanguageModel } from '../type';

	import { ModelService, type UpdateModelRequest } from '$lib/api/model/v1/model_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { model }: { model: LargeLanguageModel } = $props();

	const transport: Transport = getContext('transport');
	const modelClient = createClient(ModelService, transport);

	const reloadManager: ReloadManager = getContext('reloadManager');

	let isNameInvalid = $state(false);
	let isNamespaceInvalid = $state(false);
	let isRequestsInvalid = $state(false);
	let isLimitsInvalid = $state(false);

	const defaults = {
		scope: $currentKubernetes?.scope,
		facility: $currentKubernetes?.name,
		name: model.name,
		namespace: model.namespace,
		requests: model.requests,
		limits: model.limits,
	} as UpdateModelRequest;
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
		{m.update()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.update()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.requests()}</Form.Label>
					<SingleInput.General
						required
						type="number"
						bind:value={request.requests}
						bind:invalid={isNamespaceInvalid}
						transformer={(value) => String(value)}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.limits()}</Form.Label>
					<SingleInput.General
						required
						type="number"
						bind:value={request.limits}
						bind:invalid={isNamespaceInvalid}
						transformer={(value) => String(value)}
					/>
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
					disabled={isNameInvalid || isNamespaceInvalid || isRequestsInvalid || isLimitsInvalid}
					onclick={() => {
						console.log(request);
						toast.promise(() => modelClient.updateModel(request), {
							loading: `Updating ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Update ${request.name} successfully`;
							},
							error: (error) => {
								let message = `Fail to update ${request.name}`;
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
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
