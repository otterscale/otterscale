<script lang="ts" module>
	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';

	import Test from './action-test.svelte';
</script>

<script lang="ts">
	let {
		scope,
		testResult
	}: {
		scope: string;
		testResult: TestResult;
	} = $props();
</script>

<Test {scope} {testResult} />
