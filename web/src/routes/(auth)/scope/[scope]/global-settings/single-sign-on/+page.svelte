<script lang="ts">
	import { page } from '$app/state';
	import { SingleSignOn } from '$lib/components/settings/single-sign-on';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.globalSettings(page.params.scope)],
		current: { title: m.sso(), url: '' },
	});
</script>

<SingleSignOn />
