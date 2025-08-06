<script lang="ts">
	import { page } from '$app/state';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.setupScope(page.params.scope)],
		current: dynamicPaths.setupScopeCeph(page.params.scope)
	});
</script>
