<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { InstanceType } from '$lib/components/settings/instance-type';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs, currentKubernetes } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/settings', { scope: page.params.scope! }),
		},
		{
			title: m.instance_type(),
			url: resolve('/(auth)/scope/[scope]/settings/instance-type', { scope: page.params.scope! }),
		},
	]);
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}

	<InstanceType {scope} {facility} {namespace} />
{/if}
