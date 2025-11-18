<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { InstanceType } from '$lib/components/settings/instance-type';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/settings', { scope: page.params.scope! })
		},
		{
			title: m.instance_type(),
			url: resolve('/(auth)/scope/[scope]/settings/instance-type', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<InstanceType scope={page.params.scope!} />
{/key}
