<script lang="ts" module>
	import type { DeleteTestResultRequest, TestResult } from '$gen/api/bist/v1/bist_pb';
	import { BISTService } from '$lib/api/bist/v1/bist_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { StateController } from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
</script>

<script lang="ts">
	let {
		testResult,
	}: {
		testResult: TestResult;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	let invalid = $state(false)
	const client = createClient(BISTService, transport);
	const requestManager = new RequestManager<DeleteTestResultRequest>({
		name: testResult.name
	} as DeleteTestResultRequest);
	const stateController = new StateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="text-destructive flex h-full w-full items-center gap-2">
		<Icon icon="ph:trash" />
		Delete
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Delete Test Result</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Help>
					Please type the test name exactly to confirm deletion. This action cannot
					be undone.
				</Form.Help>
				<Form.Field>
					<Form.Label>Test Name</Form.Label>

					<SingleInput.Confirm
						required
						target={testResult.name}
						bind:value={requestManager.request.name}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={() => {requestManager.reset()}}>
				Cancel
			</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					disabled={invalid}
					onclick={() => {
						stateController.close();
						toast.promise(() => client.deleteTestResult(requestManager.request), {
							loading: `Deleting ${requestManager.request.name}...`,
							success: () => {
								reloadManager.force();
								return `Delete ${requestManager.request.name}`;
							},
							error: (error) => {
								let message = `Fail to delete ${requestManager.request.name}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						requestManager.reset();
						stateController.close();
					}}
				>
					Delete
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
