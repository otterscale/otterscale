<script lang="ts">
	import { page } from '$app/state';
	import { InstanceType } from '$lib/components/settings/instance-type';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.scopeBasedSettings(page.params.scope)],
		current: { title: m.instance_type(), url: '' },
	});
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}

	<InstanceType {scope} {facility} {namespace} />
{/if}
