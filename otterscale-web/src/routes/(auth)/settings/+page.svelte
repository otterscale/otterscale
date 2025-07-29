<script lang="ts">
	import { homePath, settingsPath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: homePath, current: settingsPath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
