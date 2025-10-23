<script lang="ts" module>
	import Test from './action-test.svelte';

	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';
</script>

<script lang="ts">
	let {
		testResult,
	}: {
		testResult: TestResult;
	} = $props();
</script>

<Test {testResult} />
