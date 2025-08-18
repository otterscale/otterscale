<script lang="ts" module>
	import type { TestResult } from '$lib/api/bist/v1/bist_pb';
	import Test from './action-test.svelte';
</script>

<script lang="ts">
	let {
		testResult,
	}: {
		testResult: TestResult;
	} = $props();
</script>

<Test testResult={testResult} />
