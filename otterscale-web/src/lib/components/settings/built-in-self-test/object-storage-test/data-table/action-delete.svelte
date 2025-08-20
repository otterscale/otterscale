<script lang="ts">
	import type { DeleteTestResultRequest, TestResult } from '$lib/api/bist/v1/bist_pb';
	import { BISTService } from '$lib/api/bist/v1/bist_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	let {
		testResult
	}: {
		testResult: TestResult;
	} = $props();

	const reloadManager: ReloadManager = getContext('reloadManager');

	const transport: Transport = getContext('transport');
	const client = createClient(BISTService, transport);

	let invalid = $state(false);

	const defaults = {
		name: testResult.name
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

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="text-destructive flex h-full w-full items-center gap-2">
		<Icon icon="ph:trash" />
		Delete
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Delete Test Result</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Help>
					Please type the test name exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
				<Form.Field>
					<Form.Label>Test Name</Form.Label>

					<SingleInput.Confirm
						required
						target={testResult.name}
						bind:value={request.name}
						bind:invalid
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel
				onclick={() => {
					reset();
				}}>Cancel</AlertDialog.Cancel
			>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
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
					Delete
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
