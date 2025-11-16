<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { BlockDevice } from '$lib/components/storage/block-device';
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
			title: m.block_device(),
			url: resolve('/(auth)/scope/[scope]/storage/block-device', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<BlockDevice scope={page.params.scope!} />
{/key}
