<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		DeleteTestResultRequest,
		TestResult
	} from '$lib/api/configuration/v1/configuration_pb';
	import { ConfigurationService } from '$lib/api/configuration/v1/configuration_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		testResult
	}: {
		testResult: TestResult;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);

	let invalid = $state(false);

	const defaults = {
		name: ''
	} as DeleteTestResultRequest;
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
		<Modal.Header>{m.delete_test_result()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Help>
					{m.deletion_warning({ identifier: m.name() })}
				</Form.Help>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>

					<SingleInput.Confirm required target={testResult.name ?? ''} bind:value={request.name} />
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
						toast.promise(() => client.deleteTestResult(request), {
							loading: `Deleting ${request.name}...`,
							success: () => {
								reloadManager.force();
								return `Delete ${request.name}`;
							},
							error: (error) => {
								let message = `Fail to delete ${request.name}`;
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
					{m.delete()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
