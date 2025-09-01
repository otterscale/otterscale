<script lang="ts">
	import { page } from '$app/state';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { Construction } from '$lib/components/construction';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.compute(page.params.scope)],
		current: dynamicPaths.computeVirtualMachine(page.params.scope),
	});
</script>

<Construction />
