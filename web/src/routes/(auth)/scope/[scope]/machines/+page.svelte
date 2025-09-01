<script lang="ts">
	import { page } from '$app/state';
	import Dashboard from '$lib/components/machines/dashboard/index.svelte';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.machines(page.params.scope) });
</script>

<Dashboard />
