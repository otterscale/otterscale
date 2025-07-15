<script lang="ts">
	import * as Tabs from '$lib/components/custom/tabs/naive/index';
	import type { TabsRootProps } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { StepManager } from './utils.svelte';

	let { ...Props }: TabsRootProps = $props();
	const stepManager: StepManager = getContext('StepManager');
</script>

<Tabs.Root value={String(stepManager.step)} {...Props} />
