<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { OSD } from '$lib/components/storage/osd';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScope = $derived(page.params.scope!);
	let selectedFacility = $state('ceph-mon');

	breadcrumbs.set([
		{
			title: m.storage(),
			url: resolve('/(auth)/scope/[scope]/storage', { scope: page.params.scope! })
		},
		{
			title: m.osd(),
			url: resolve('/(auth)/scope/[scope]/storage/osd', { scope: page.params.scope! })
		}
	]);
</script>

{#key selectedScope + selectedFacility}
	<OSD bind:selectedScope bind:selectedFacility />
{/key}
