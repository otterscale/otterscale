<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { FileSystem } from '$lib/components/storage/file-system/index';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScope = $derived(page.params.scope!);
	let selectedFacility = $state('ceph-mon');
	let selectedVolume = $state('ceph-fs');
	let selectedSubvolumeGroupName = $state('');

	breadcrumbs.set([
		{
			title: m.storage(),
			url: resolve('/(auth)/scope/[scope]/storage', { scope: page.params.scope! })
		},
		{
			title: m.file_system(),
			url: resolve('/(auth)/scope/[scope]/storage/file-system', { scope: page.params.scope! })
		}
	]);
</script>

{#key selectedScope + selectedFacility + selectedVolume + selectedSubvolumeGroupName}
	<FileSystem
		bind:selectedScope
		bind:selectedFacility
		bind:selectedVolume
		bind:selectedSubvolumeGroupName
	/>
{/key}
