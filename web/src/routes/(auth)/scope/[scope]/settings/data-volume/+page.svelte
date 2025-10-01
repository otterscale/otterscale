<script lang="ts">
	import { page } from '$app/state';
	import { DataVolume } from '$lib/components/settings/data-volume';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.data_volume(), url: '' },
	});
</script>

{#if $currentKubernetes}
	{@const scopeUuid = $currentKubernetes.scopeUuid}
	{@const facilityName = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}

	<DataVolume {scopeUuid} {facilityName} {namespace} />
{/if}
