<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';

	import { getAlertIcon } from '../single/alert-icon.svelte';

	import type { IterationManager } from './utils.svelte';
</script>

<script lang="ts">
	let iterationManager: IterationManager = getContext('IterationManager');
</script>

<Icon
	data-slot="alert-icon"
	icon={getAlertIcon(iterationManager.alerts[iterationManager.value.index].variant)}
/>
