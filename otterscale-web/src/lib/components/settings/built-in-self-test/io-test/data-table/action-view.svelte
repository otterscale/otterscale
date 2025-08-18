<script lang="ts" module>
	import type { TestResult } from '$lib/api/bist/v1/bist_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { StateController } from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { formatCapacity, formatLatencyNano } from '$lib/formatter';
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
			I/O Test Output
		</AlertDialog.Header>
		<Form.Root>
			{#if testResult?.kind?.value?.output && testResult.kind.case == 'fio'}
				{@const readOutput = testResult.kind.value.output.read}
				{@const writeOutput = testResult.kind.value.output.write}
				{@const trimOutput = testResult.kind.value.output.trim}
				
				<Form.Fieldset>
					<Form.Legend>Read</Form.Legend>
					{@const readIoBytes = formatCapacity(Number(readOutput?.ioBytes || 0))}
					{@const readBandwidthBytes = formatCapacity(Number(readOutput?.bandwidthBytes || 0))}
					{@const readMinLatency = formatLatencyNano(Number(readOutput?.latency?.minNanoseconds || 0))}
					{@const readMaxLatency = formatLatencyNano(Number(readOutput?.latency?.maxNanoseconds || 0))}
					{@const readMeanLatency = formatLatencyNano(readOutput?.latency?.meanNanoseconds || 0)}
					<Form.Description>IO Bytes: {readIoBytes.value} {readIoBytes.unit}</Form.Description>
					<Form.Description>Bandwidth Bytes: {readBandwidthBytes.value} {readBandwidthBytes.unit}</Form.Description>
					<Form.Description>IO Per Second: {readOutput?.ioPerSecond}</Form.Description>
					<Form.Description>Total Ios: {readOutput?.totalIos}</Form.Description>
					<Form.Description>Min Latency: {readMinLatency.value} {readMinLatency.unit}</Form.Description>
					<Form.Description>Max Latency: {readMaxLatency.value} {readMaxLatency.unit}</Form.Description>
					<Form.Description>Mean Latency: {readMeanLatency.value} {readMeanLatency.unit}</Form.Description>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>Write</Form.Legend>
					{@const writeIoBytes = formatCapacity(Number(writeOutput?.ioBytes || 0))}
					{@const writeBandwidthBytes = formatCapacity(Number(writeOutput?.bandwidthBytes || 0))}
					{@const writeMinLatency = formatLatencyNano(Number(writeOutput?.latency?.minNanoseconds || 0))}
					{@const writeMaxLatency = formatLatencyNano(Number(writeOutput?.latency?.maxNanoseconds || 0))}
					{@const writeMeanLatency = formatLatencyNano(writeOutput?.latency?.meanNanoseconds || 0)}
					<Form.Description>IO Bytes: {writeIoBytes.value} {writeIoBytes.unit}</Form.Description>
					<Form.Description>Bandwidth Bytes: {writeBandwidthBytes.value} {writeBandwidthBytes.unit}</Form.Description>
					<Form.Description>IO Per Second: {writeOutput?.ioPerSecond}</Form.Description>
					<Form.Description>Total Ios: {writeOutput?.totalIos}</Form.Description>
					<Form.Description>Min Latency: {writeMinLatency.value} {writeMinLatency.unit}</Form.Description>
					<Form.Description>Max Latency: {writeMaxLatency.value} {writeMaxLatency.unit}</Form.Description>
					<Form.Description>Mean Latency: {writeMeanLatency.value} {writeMeanLatency.unit}</Form.Description>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>Trim</Form.Legend>
					{@const trimIoBytes = formatCapacity(Number(trimOutput?.ioBytes || 0))}
					{@const trimBandwidthBytes = formatCapacity(Number(trimOutput?.bandwidthBytes || 0))}
					{@const trimMinLatency = formatLatencyNano(Number(trimOutput?.latency?.minNanoseconds || 0))}
					{@const trimMaxLatency = formatLatencyNano(Number(trimOutput?.latency?.maxNanoseconds || 0))}
					{@const trimMeanLatency = formatLatencyNano(trimOutput?.latency?.meanNanoseconds || 0)}
					<Form.Description>IO Bytes: {trimIoBytes.value} {trimIoBytes.unit}</Form.Description>
					<Form.Description>Bandwidth Bytes: {trimBandwidthBytes.value} {trimBandwidthBytes.unit}</Form.Description>
					<Form.Description>IO Per Second: {trimOutput?.ioPerSecond}</Form.Description>
					<Form.Description>Total Ios: {trimOutput?.totalIos}</Form.Description>
					<Form.Description>Min Latency: {trimMinLatency.value} {trimMinLatency.unit}</Form.Description>
					<Form.Description>Max Latency: {trimMaxLatency.value} {trimMaxLatency.unit}</Form.Description>
					<Form.Description>Mean Latency: {trimMeanLatency.value} {trimMeanLatency.unit}</Form.Description>
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
