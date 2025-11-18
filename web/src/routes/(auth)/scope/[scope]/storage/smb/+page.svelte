<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { SMB } from '$lib/components/storage/smb/index';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
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

{#key page.params.scope!}
	<SMB scope={page.params.scope!} />
{/key}
