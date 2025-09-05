<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { TestResult } from '$lib/api/bist/v1/bist_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { formatCapacity, formatLatencyNano } from '$lib/formatter';
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

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:eye" />
		{m.view()}
	</AlertDialog.Trigger>
	<AlertDialog.Content class="min-w-[50vw]">
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			{m.io_test_output()}
		</AlertDialog.Header>
		<Form.Root>
			{#if testResult?.kind?.value?.output && testResult.kind.case == 'fio'}
				{@const readOutput = testResult.kind.value.output.read}
				{@const writeOutput = testResult.kind.value.output.write}
				{@const trimOutput = testResult.kind.value.output.trim}

				<Form.Fieldset>
					<Form.Legend>{m.read()}</Form.Legend>
					{@const readIoBytes = formatCapacity(Number(readOutput?.ioBytes || 0))}
					{@const readBandwidthBytes = formatCapacity(Number(readOutput?.bandwidthBytes || 0))}
					{@const readMinLatency = formatLatencyNano(Number(readOutput?.latency?.minNanoseconds || 0))}
					{@const readMaxLatency = formatLatencyNano(Number(readOutput?.latency?.maxNanoseconds || 0))}
					{@const readMeanLatency = formatLatencyNano(readOutput?.latency?.meanNanoseconds || 0)}
					<Form.Description>{m.io_bytes()}: {readIoBytes.value} {readIoBytes.unit}</Form.Description>
					<Form.Description
						>{m.bandwidth_bytes()}: {readBandwidthBytes.value}
						{readBandwidthBytes.unit}</Form.Description
					>
					<Form.Description>{m.io()}/{m.sec()}: {readOutput?.ioPerSecond}</Form.Description>
					<Form.Description>{m.total_ios()}: {readOutput?.totalIos}</Form.Description>
					<Form.Description>{m.min_latency()}: {readMinLatency.value} {readMinLatency.unit}</Form.Description>
					<Form.Description>{m.max_latency()}: {readMaxLatency.value} {readMaxLatency.unit}</Form.Description>
					<Form.Description
						>{m.mean_latency()}: {readMeanLatency.value} {readMeanLatency.unit}</Form.Description
					>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>{m.write()}</Form.Legend>
					{@const writeIoBytes = formatCapacity(Number(writeOutput?.ioBytes || 0))}
					{@const writeBandwidthBytes = formatCapacity(Number(writeOutput?.bandwidthBytes || 0))}
					{@const writeMinLatency = formatLatencyNano(Number(writeOutput?.latency?.minNanoseconds || 0))}
					{@const writeMaxLatency = formatLatencyNano(Number(writeOutput?.latency?.maxNanoseconds || 0))}
					{@const writeMeanLatency = formatLatencyNano(writeOutput?.latency?.meanNanoseconds || 0)}
					<Form.Description>{m.io_bytes()}: {writeIoBytes.value} {writeIoBytes.unit}</Form.Description>
					<Form.Description
						>{m.bandwidth_bytes()}: {writeBandwidthBytes.value}
						{writeBandwidthBytes.unit}</Form.Description
					>
					<Form.Description>{m.io()}/{m.sec()}: {writeOutput?.ioPerSecond}</Form.Description>
					<Form.Description>{m.total_ios()}: {writeOutput?.totalIos}</Form.Description>
					<Form.Description
						>{m.min_latency()}: {writeMinLatency.value} {writeMinLatency.unit}</Form.Description
					>
					<Form.Description
						>{m.max_latency()}: {writeMaxLatency.value} {writeMaxLatency.unit}</Form.Description
					>
					<Form.Description
						>{m.mean_latency()}: {writeMeanLatency.value} {writeMeanLatency.unit}</Form.Description
					>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>{m.trim()}</Form.Legend>
					{@const trimIoBytes = formatCapacity(Number(trimOutput?.ioBytes || 0))}
					{@const trimBandwidthBytes = formatCapacity(Number(trimOutput?.bandwidthBytes || 0))}
					{@const trimMinLatency = formatLatencyNano(Number(trimOutput?.latency?.minNanoseconds || 0))}
					{@const trimMaxLatency = formatLatencyNano(Number(trimOutput?.latency?.maxNanoseconds || 0))}
					{@const trimMeanLatency = formatLatencyNano(trimOutput?.latency?.meanNanoseconds || 0)}
					<Form.Description>{m.io_bytes()}: {trimIoBytes.value} {trimIoBytes.unit}</Form.Description>
					<Form.Description
						>{m.bandwidth_bytes()}: {trimBandwidthBytes.value}
						{trimBandwidthBytes.unit}</Form.Description
					>
					<Form.Description>{m.io()}/{m.sec()}: {trimOutput?.ioPerSecond}</Form.Description>
					<Form.Description>{m.total_ios()}: {trimOutput?.totalIos}</Form.Description>
					<Form.Description>{m.min_latency()}: {trimMinLatency.value} {trimMinLatency.unit}</Form.Description>
					<Form.Description>{m.max_latency()}: {trimMaxLatency.value} {trimMaxLatency.unit}</Form.Description>
					<Form.Description
						>{m.mean_latency()}: {trimMeanLatency.value} {trimMeanLatency.unit}</Form.Description
					>
				</Form.Fieldset>
			{:else}
				<Form.Description class="text-xs font-light">
					<p>{m.no_data()}</p>
				</Form.Description>
			{/if}
		</Form.Root>

		<AlertDialog.Footer>
			<AlertDialog.Cancel>{m.close()}</AlertDialog.Cancel>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
