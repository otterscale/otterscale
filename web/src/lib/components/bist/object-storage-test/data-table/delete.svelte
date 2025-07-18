<script lang="ts" module>
	import type { DeleteTestResultRequest, TestResult } from '$gen/api/bist/v1/bist_pb';
	import { BISTService } from '$gen/api/bist/v1/bist_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		testResult,
		data = $bindable()
	}: {
		testResult: TestResult;
		data: Writable<TestResult[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		name: testResult.name
	} as DeleteTestResultRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
	const transport: Transport = getContext('transport');
	const bistClient = createClient(BISTService, transport);
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

					<SingleInput.DeletionConfirm
						required
						target={testResult.name}
						bind:value={request.name}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						stateController.close();
						bistClient
							.deleteTestResult(request)
							.then((r) => {
								toast.success(`Delete ${request.name}`);
								bistClient
									.listTestResults({})
									.then((r) => {
										data.set(r.testResults.filter((result) => result.kind.case === 'warp' ));
									});
							})
							.catch((e) => {
								toast.error(`Fail to delete test result: ${e.toString()}`);
							});
					}}
				>
					Delete
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
