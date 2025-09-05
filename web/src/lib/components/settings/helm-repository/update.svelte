<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		EnvironmentService,
		type UpdateConfigHelmRepositoriesRequest,
	} from '$lib/api/environment/v1/environment_pb';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		urls,
	}: {
		urls: Writable<string[]>;
	} = $props();

	const transport: Transport = getContext('transport');

	const client = createClient(EnvironmentService, transport);
	const defaults = {
		urls: $urls,
	} as UpdateConfigHelmRepositoriesRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<span>
	<Modal.Root bind:open>
		<Modal.Trigger variant="default">
			<Icon icon="ph:pencil" />
			{m.edit()}
		</Modal.Trigger>
		<Modal.Content>
			<Modal.Header>{m.edit_helm_repository()}</Modal.Header>
			<Form.Root>
				<Form.Fieldset>
					<Form.Field>
						<Form.Label>{m.url()}</Form.Label>
						<MultipleInput.Root type="text" bind:values={request.urls}>
							<MultipleInput.Viewer />
							<MultipleInput.Controller>
								<MultipleInput.Input />
								<MultipleInput.Add />
								<MultipleInput.Clear />
							</MultipleInput.Controller>
						</MultipleInput.Root>
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
						onclick={() => {
							toast.promise(() => client.updateConfigHelmRepositories(request), {
								loading: 'Loading...',
								success: () => {
									client.getConfigHelmRepositories({}).then((response) => {
										urls.set(response.urls);
									});
									return `Update success`;
								},
								error: (error) => {
									let message = `Fail to update`;
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
</span>
