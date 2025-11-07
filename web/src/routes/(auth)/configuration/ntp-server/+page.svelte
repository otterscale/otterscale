<script lang="ts">
	import { resolve } from '$app/paths';
	import { NTPServer } from '$lib/components/settings/ntp-server';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/configuration'),
		},
		{
			title: m.ntp_server(),
			url: resolve('/(auth)/configuration/ntp-server'),
		},
	]);
</script>

<NTPServer />
