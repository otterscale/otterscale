<script lang="ts" module>
	import type { TestResult } from '$lib/api/bist/v1/bist_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { StateController } from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { formatCapacity } from '$lib/formatter';
	import Icon from '@iconify/svelte';
</script>

<script lang="ts">
	let {
		testResult,
	}: {
		testResult: TestResult;
	} = $props();

	const stateController = new StateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:eye" />
		View
	</AlertDialog.Trigger>
	<AlertDialog.Content class="min-w-[50vw]">
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Object Storage Test Output
		</AlertDialog.Header>
		<Form.Root>
			{#if testResult?.kind?.value?.output && testResult.kind.case == 'warp'}
				{@const getOutput = testResult.kind.value.output.get}
				{@const putOutput = testResult.kind.value.output.put}
				{@const deleteOutput = testResult.kind.value.output.delete}
				
				<Form.Fieldset>
					<Form.Legend>GET</Form.Legend>
					{@const getTotalBytes = formatCapacity(getOutput?.totalBytes || 0)}
					{@const getFastest = formatCapacity(getOutput?.bytes?.fastestPerSecond || 0)}
					{@const getMedian = formatCapacity(getOutput?.bytes?.medianPerSecond || 0)}
					{@const getSlowest = formatCapacity(getOutput?.bytes?.slowestPerSecond || 0)}
					<Form.Description>Total Bytes: {getTotalBytes.value} {getTotalBytes.unit}</Form.Description>
					<Form.Description>Total Objects: {getOutput?.totalObjects}</Form.Description>
					<Form.Description>Total Operations: {getOutput?.totalOperations}</Form.Description>
					<Form.Description>Bytes - Fastest/sec: {getFastest.value} {getFastest.unit}/s</Form.Description>
					<Form.Description>Bytes - Median/sec: {getMedian.value} {getMedian.unit}/s</Form.Description>
					<Form.Description>Bytes - Slowest/sec: {getSlowest.value} {getSlowest.unit}/s</Form.Description>
					<Form.Description>Objects - Fastest/sec: {getOutput?.objects?.fastestPerSecond}</Form.Description>
					<Form.Description>Objects - Median/sec: {getOutput?.objects?.medianPerSecond}</Form.Description>
					<Form.Description>Objects - Slowest/sec: {getOutput?.objects?.slowestPerSecond}</Form.Description>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>PUT</Form.Legend>
					{@const putTotalBytes = formatCapacity(putOutput?.totalBytes || 0)}
					{@const putFastest = formatCapacity(putOutput?.bytes?.fastestPerSecond || 0)}
					{@const putMedian = formatCapacity(putOutput?.bytes?.medianPerSecond || 0)}
					{@const putSlowest = formatCapacity(putOutput?.bytes?.slowestPerSecond || 0)}
					<Form.Description>Total Bytes: {putTotalBytes.value} {putTotalBytes.unit}</Form.Description>
					<Form.Description>Total Objects: {putOutput?.totalObjects}</Form.Description>
					<Form.Description>Total Operations: {putOutput?.totalOperations}</Form.Description>
					<Form.Description>Bytes - Fastest/sec: {putFastest.value} {putFastest.unit}/s</Form.Description>
					<Form.Description>Bytes - Median/sec: {putMedian.value} {putMedian.unit}/s</Form.Description>
					<Form.Description>Bytes - Slowest/sec: {putSlowest.value} {putSlowest.unit}/s</Form.Description>
					<Form.Description>Objects - Fastest/sec: {putOutput?.objects?.fastestPerSecond}</Form.Description>
					<Form.Description>Objects - Median/sec: {putOutput?.objects?.medianPerSecond}</Form.Description>
					<Form.Description>Objects - Slowest/sec: {putOutput?.objects?.slowestPerSecond}</Form.Description>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>DELETE</Form.Legend>
					{@const deleteTotalBytes = formatCapacity(deleteOutput?.totalBytes || 0)}
					{@const deleteFastest = formatCapacity(deleteOutput?.bytes?.fastestPerSecond || 0)}
					{@const deleteMedian = formatCapacity(deleteOutput?.bytes?.medianPerSecond || 0)}
					{@const deleteSlowest = formatCapacity(deleteOutput?.bytes?.slowestPerSecond || 0)}
					<Form.Description>Total Bytes: {deleteTotalBytes.value} {deleteTotalBytes.unit}</Form.Description>
					<Form.Description>Total Objects: {deleteOutput?.totalObjects}</Form.Description>
					<Form.Description>Total Operations: {deleteOutput?.totalOperations}</Form.Description>
					<Form.Description>Bytes - Fastest/sec: {deleteFastest.value} {deleteFastest.unit}/s</Form.Description>
					<Form.Description>Bytes - Median/sec: {deleteMedian.value} {deleteMedian.unit}/s</Form.Description>
					<Form.Description>Bytes - Slowest/sec: {deleteSlowest.value} {deleteSlowest.unit}/s</Form.Description>
					<Form.Description>Objects - Fastest/sec: {deleteOutput?.objects?.fastestPerSecond}</Form.Description>
					<Form.Description>Objects - Median/sec: {deleteOutput?.objects?.medianPerSecond}</Form.Description>
					<Form.Description>Objects - Slowest/sec: {deleteOutput?.objects?.slowestPerSecond}</Form.Description>
				</Form.Fieldset>
			{:else}
				<Form.Description class="text-xs font-light">
					<p>No input data available.</p>
				</Form.Description>
			{/if}
		</Form.Root>

		<AlertDialog.Footer>
			<AlertDialog.Cancel>Close</AlertDialog.Cancel>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
