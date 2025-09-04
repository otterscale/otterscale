<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { TestResult } from '$lib/api/bist/v1/bist_pb';
	import * as Form from '$lib/components/custom/form';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		testResult,
	}: {
		testResult: TestResult;
	} = $props();

	let open = $state(false);
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:eye" />
		{m.view()}
	</Modal.Trigger>
	<Modal.Content class="min-w-[50vw]">
		<Modal.Header>{m.object_storage_test_output()}</Modal.Header>
		<Form.Root>
			{#if testResult?.kind?.value?.output && testResult.kind.case == 'warp'}
				{@const getOutput = testResult.kind.value.output.get}
				{@const putOutput = testResult.kind.value.output.put}
				{@const deleteOutput = testResult.kind.value.output.delete}

				<Form.Fieldset>
					<Form.Legend>{m.get()}</Form.Legend>
					{@const getTotalBytes = formatCapacity(getOutput?.totalBytes || 0)}
					{@const getFastest = formatCapacity(getOutput?.bytes?.fastestPerSecond || 0)}
					{@const getMedian = formatCapacity(getOutput?.bytes?.medianPerSecond || 0)}
					{@const getSlowest = formatCapacity(getOutput?.bytes?.slowestPerSecond || 0)}
					<Form.Description>
						{m.total_bytes()}: {getTotalBytes.value}
						{getTotalBytes.unit}
					</Form.Description>
					<Form.Description>{m.total_objects()}: {getOutput?.totalObjects}</Form.Description>
					<Form.Description>{m.total_operations()}: {getOutput?.totalOperations}</Form.Description>
					<Form.Description>
						{m.bytes_fastest()}/{m.sec()}: {getFastest.value}
						{getFastest.unit}/s
					</Form.Description>
					<Form.Description>
						{m.bytes_median()}/{m.sec()}: {getMedian.value}
						{getMedian.unit}/s
					</Form.Description>
					<Form.Description>
						{m.bytes_slowest()}/{m.sec()}: {getSlowest.value}
						{getSlowest.unit}/s
					</Form.Description>
					<Form.Description>
						{m.objects_fastest()}/{m.sec()}: {getOutput?.objects?.fastestPerSecond}
					</Form.Description>
					<Form.Description>
						{m.objects_median()}/{m.sec()}: {getOutput?.objects?.medianPerSecond}
					</Form.Description>
					<Form.Description>
						{m.objects_slowest()}/{m.sec()}: {getOutput?.objects?.slowestPerSecond}
					</Form.Description>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>{m.put()}</Form.Legend>
					{@const putTotalBytes = formatCapacity(putOutput?.totalBytes || 0)}
					{@const putFastest = formatCapacity(putOutput?.bytes?.fastestPerSecond || 0)}
					{@const putMedian = formatCapacity(putOutput?.bytes?.medianPerSecond || 0)}
					{@const putSlowest = formatCapacity(putOutput?.bytes?.slowestPerSecond || 0)}
					<Form.Description>
						{m.total_bytes()}: {putTotalBytes.value}
						{putTotalBytes.unit}
					</Form.Description>
					<Form.Description>
						{m.total_objects()}: {putOutput?.totalObjects}
					</Form.Description>
					<Form.Description>
						{m.total_operations()}: {putOutput?.totalOperations}
					</Form.Description>
					<Form.Description>
						{m.bytes_fastest()}/{m.sec()}: {putFastest.value}
						{putFastest.unit}/s
					</Form.Description>
					<Form.Description>
						{m.bytes_median()}/{m.sec()}: {putMedian.value}
						{putMedian.unit}/s
					</Form.Description>
					<Form.Description>
						{m.bytes_slowest()}/{m.sec()}: {putSlowest.value}
						{putSlowest.unit}/s
					</Form.Description>
					<Form.Description>
						{m.objects_fastest()}/{m.sec()}: {putOutput?.objects?.fastestPerSecond}
					</Form.Description>
					<Form.Description>
						{m.objects_median()}/{m.sec()}: {putOutput?.objects?.medianPerSecond}
					</Form.Description>
					<Form.Description>
						{m.objects_slowest()}/{m.sec()}: {putOutput?.objects?.slowestPerSecond}
					</Form.Description>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>{m.delete()}</Form.Legend>
					{@const deleteTotalBytes = formatCapacity(deleteOutput?.totalBytes || 0)}
					{@const deleteFastest = formatCapacity(deleteOutput?.bytes?.fastestPerSecond || 0)}
					{@const deleteMedian = formatCapacity(deleteOutput?.bytes?.medianPerSecond || 0)}
					{@const deleteSlowest = formatCapacity(deleteOutput?.bytes?.slowestPerSecond || 0)}
					<Form.Description>
						{m.total_bytes()}: {deleteTotalBytes.value}
						{deleteTotalBytes.unit}
					</Form.Description>
					<Form.Description>
						{m.total_objects()}: {deleteOutput?.totalObjects}
					</Form.Description>
					<Form.Description>
						{m.total_operations()}: {deleteOutput?.totalOperations}
					</Form.Description>
					<Form.Description>
						{m.bytes_fastest()}/{m.sec()}: {deleteFastest.value}
						{deleteFastest.unit}/s
					</Form.Description>
					<Form.Description>
						{m.bytes_median()}/{m.sec()}: {deleteMedian.value}
						{deleteMedian.unit}/s
					</Form.Description>
					<Form.Description>
						{m.bytes_slowest()}/{m.sec()}: {deleteSlowest.value}
						{deleteSlowest.unit}/s
					</Form.Description>
					<Form.Description>
						{m.objects_fastest()}/{m.sec()}: {deleteOutput?.objects?.fastestPerSecond}
					</Form.Description>
					<Form.Description>
						{m.objects_median()}/{m.sec()}: {deleteOutput?.objects?.medianPerSecond}
					</Form.Description>
					<Form.Description>
						{m.objects_slowest()}/{m.sec()}: {deleteOutput?.objects?.slowestPerSecond}
					</Form.Description>
				</Form.Fieldset>
			{:else}
				<Form.Description class="text-xs font-light">
					<p>{m.no_data()}</p>
				</Form.Description>
			{/if}
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>{m.close()}</Modal.Cancel>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
