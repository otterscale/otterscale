<script lang="ts">
	import { accountPath, homePath } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: homePath, current: accountPath });
</script>
