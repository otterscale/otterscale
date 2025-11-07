<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { NTPServer } from '$lib/components/settings/ntp-server';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/global-settings', { scope: page.params.scope! }),
		},
		{
			title: m.ntp_server(),
			url: resolve('/(auth)/scope/[scope]/global-settings/ntp-server', { scope: page.params.scope! }),
		},
	]);
</script>

<NTPServer />
