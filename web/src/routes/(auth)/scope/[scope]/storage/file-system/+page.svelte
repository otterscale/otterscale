<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { FileSystem } from '$lib/components/storage/file-system/index';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
	const volume = 'ceph-fs';
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

{#key page.params.scope! + selectedSubvolumeGroupName}
	<FileSystem scope={page.params.scope!} {volume} bind:selectedSubvolumeGroupName />
{/key}
