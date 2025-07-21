<script lang="ts" module>
	import type { TestResult } from '$gen/api/bist/v1/bist_pb';
	import Icon from '@iconify/svelte';
	import type { Writable } from 'svelte/store';
	import * as MultipleStepModal from './mutiple-step-modal';
	import TestStepModal from './test-step-modal.svelte';
</script>

<script lang="ts">
	let {
		testResult,
		data = $bindable()
	}: {
		testResult: TestResult;
		data: Writable<TestResult[]>;
	} = $props();
</script>

{#snippet trigger()}
	<MultipleStepModal.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:play" />
		Retest
	</MultipleStepModal.Trigger>
{/snippet}

<TestStepModal testResult={testResult} data={data} {trigger} />
