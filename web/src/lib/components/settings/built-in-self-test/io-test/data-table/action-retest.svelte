<script lang="ts" module>
	import Test from './action-test.svelte';

	import type { TestResult } from '$lib/api/bist/v1/bist_pb';
</script>

<script lang="ts">
	let {
		testResult,
	}: {
		testResult: TestResult;
	} = $props();
</script>

<Test {testResult} />
