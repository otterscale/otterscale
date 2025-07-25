<script lang="ts">
	import { databasesPath, homePath } from '$lib/path';
	import { activeScope, breadcrumb, scopeLoading } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: homePath, current: databasesPath });
</script>

{#if $scopeLoading}
	loading scope
{:else}
	current scope: {$activeScope.uuid}
{/if}
