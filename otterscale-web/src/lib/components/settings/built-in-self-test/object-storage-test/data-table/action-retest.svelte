<script lang="ts">
	import type { TestResult } from '$lib/api/bist/v1/bist_pb';
	import Test from './action-test.svelte';

	let {
		testResult
	}: {
		testResult: TestResult;
	} = $props();
</script>

<Test {testResult} />
