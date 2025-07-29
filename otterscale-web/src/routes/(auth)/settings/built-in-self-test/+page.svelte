<script lang="ts">
	import { settingsBISTPath, settingsPath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: settingsPath, current: settingsBISTPath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
