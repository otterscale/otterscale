<script lang="ts">
	import { page } from '$app/state';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { Construction } from '$lib/components/construction';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.databases(page.params.scope)],
		current: dynamicPaths.databasesNoSQL(page.params.scope)
	});
</script>

<Construction />
