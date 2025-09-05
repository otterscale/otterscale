<script lang="ts">
	import { page } from '$app/state';
	import { Construction } from '$lib/components/construction';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [],
		current: dynamicPaths.databases(page.params.scope),
	});
</script>

<Construction />
