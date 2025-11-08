<script lang="ts" module>
	import type { TabsRootProps } from 'bits-ui';
	import { getContext } from 'svelte';

	import * as Tabs from '$lib/components/ui/tabs';

	import type { StepManager } from './utils.svelte';
</script>

<script lang="ts">
	let { ...Props }: TabsRootProps = $props();
	const stepManager: StepManager = getContext('StepManager');
</script>

<Tabs.Root value={String(stepManager.step)} {...Props} />
