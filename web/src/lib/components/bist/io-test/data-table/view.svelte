<script lang="ts" module>
	import type { TestResult } from '$gen/api/bist/v1/bist_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import Icon from '@iconify/svelte';
</script>

<script lang="ts">
	let {
		testResult,
	}: {
		testResult: TestResult;
	} = $props();

	const stateController = new DialogStateController(false);

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
				<Form.Fieldset>
					<Form.Legend>Read</Form.Legend>
					<Form.Description>IO Bytes: {testResult.kind.value.output.read?.ioBytes}</Form.Description>
					<Form.Description>Bandwidth Bytes: {testResult.kind.value.output.read?.bandwidthBytes}</Form.Description>
					<Form.Description>IO Per Second: {testResult.kind.value.output.read?.ioPerSecond}</Form.Description>
					<Form.Description>Total Ios: {testResult.kind.value.output.read?.totalIos}</Form.Description>
					<Form.Description>Min Latency (ns): {testResult.kind.value.output.read?.latency?.minNanoseconds}</Form.Description>
					<Form.Description>Max Latency (ns): {testResult.kind.value.output.read?.latency?.maxNanoseconds}</Form.Description>
					<Form.Description>Mean Latency (ns): {testResult.kind.value.output.read?.latency?.meanNanoseconds}</Form.Description>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>Write</Form.Legend>
					<Form.Description>IO Bytes: {testResult.kind.value.output.write?.ioBytes}</Form.Description>
					<Form.Description>Bandwidth Bytes: {testResult.kind.value.output.write?.bandwidthBytes}</Form.Description>
					<Form.Description>IO Per Second: {testResult.kind.value.output.write?.ioPerSecond}</Form.Description>
					<Form.Description>Total Ios: {testResult.kind.value.output.write?.totalIos}</Form.Description>
					<Form.Description>Min Latency (ns): {testResult.kind.value.output.write?.latency?.minNanoseconds}</Form.Description>
					<Form.Description>Max Latency (ns): {testResult.kind.value.output.write?.latency?.maxNanoseconds}</Form.Description>
					<Form.Description>Mean Latency (ns): {testResult.kind.value.output.write?.latency?.meanNanoseconds}</Form.Description>
				</Form.Fieldset>
				<Form.Fieldset>
					<Form.Legend>Trim</Form.Legend>
					<Form.Description>IO Bytes: {testResult.kind.value.output.trim?.ioBytes}</Form.Description>
					<Form.Description>Bandwidth Bytes: {testResult.kind.value.output.trim?.bandwidthBytes}</Form.Description>
					<Form.Description>IO Per Second: {testResult.kind.value.output.trim?.ioPerSecond}</Form.Description>
					<Form.Description>Total Ios: {testResult.kind.value.output.trim?.totalIos}</Form.Description>
					<Form.Description>Min Latency (ns): {testResult.kind.value.output.trim?.latency?.minNanoseconds}</Form.Description>
					<Form.Description>Max Latency (ns): {testResult.kind.value.output.trim?.latency?.maxNanoseconds}</Form.Description>
					<Form.Description>Mean Latency (ns): {testResult.kind.value.output.trim?.latency?.meanNanoseconds}</Form.Description>
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
