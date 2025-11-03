<script lang="ts">
	import { page } from '$app/state';
	import { NTPServer } from '$lib/components/settings/ntp-server';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.globalSettings(page.params.scope)],
		current: { title: m.ntp_server(), url: '' },
	});
</script>

<NTPServer />
