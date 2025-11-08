<script lang="ts" module>
	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';

	import Test from './action-test.svelte';
</script>

<script lang="ts">
	let {
		testResult
	}: {
		testResult: TestResult;
	} = $props();
</script>

<Test {testResult} />
