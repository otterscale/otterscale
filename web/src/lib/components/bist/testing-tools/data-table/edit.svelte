<script lang="ts" module>
	import type { CreateTestResultRequest, TestResult } from '$gen/api/bist/v1/bist_pb';
	import { BISTService } from '$gen/api/bist/v1/bist_pb';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	// import Create from './create.svelte';
	import * as MultipleStepModal from './mutiple-step-modal';
	import TestStepModal from './test-step-modal.svelte'
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
	} as CreateTestResultRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
	const transport: Transport = getContext('transport');
	const bistClient = createClient(BISTService, transport);
</script>

{#snippet trigger()}
	<MultipleStepModal.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:play" />
		Retest``
	</MultipleStepModal.Trigger>
{/snippet}

<TestStepModal testResult={testResult} data={data} {trigger} />
