<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { FileSystem } from '$lib/components/storage/file-system/index';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
	let selectedScope = $derived($activeScope ? $activeScope.name : '');
	let selectedFacility = $state('ceph-mon');
	let selectedVolume = $state('ceph-fs');
	let selectedSubvolumeGroupName = $state('');

	breadcrumbs.set([
		{ title: m.storage(), url: resolve('/(auth)/scope/[scope]/storage', { scope: page.params.scope! }) },
		{
			title: m.file_system(),
			url: resolve('/(auth)/scope/[scope]/storage/file-system', { scope: page.params.scope! }),
		},
	]);
</script>

{#if $activeScope}
	{#key selectedScope + selectedFacility + selectedVolume + selectedSubvolumeGroupName}
		<FileSystem bind:selectedScope bind:selectedFacility bind:selectedVolume bind:selectedSubvolumeGroupName />
	{/key}
{/if}
