<script lang="ts" module>
	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';

	import Test from './action-test.svelte';
</script>

<script lang="ts">
	let {
		testResult,
		scope,
		reloadManager
	}: {
		testResult: TestResult;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Test {testResult} {scope} {reloadManager} />