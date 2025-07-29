<script lang="ts">
	import { settingsNetworkPath, settingsPath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: settingsPath, current: settingsNetworkPath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
