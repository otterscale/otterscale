<script lang="ts" module>
	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';

	import Test from './action-test.svelte';
</script>

<script lang="ts">
	let {
		scope,
		testResult,
		reloadManager
	}: {
		scope: string;
		testResult: TestResult;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Test {scope} {testResult} {reloadManager} />
