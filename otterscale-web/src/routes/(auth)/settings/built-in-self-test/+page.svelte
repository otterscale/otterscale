<script lang="ts">
	import { settingsBISTPath, settingsPath } from '$lib/path';
	import { activeScope, breadcrumb, scopeLoading } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: settingsPath, current: settingsBISTPath });
</script>

{#if $scopeLoading}
	loading scope
{:else}
	current scope: {$activeScope.uuid}
{/if}
