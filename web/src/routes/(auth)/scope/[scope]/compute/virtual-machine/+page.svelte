<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { VirtualMachine } from '$lib/components/compute/virtual-machine';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.compute(),
			url: resolve('/(auth)/scope/[scope]/compute', { scope: page.params.scope! })
		},
		{
			title: m.virtual_machine(),
			url: resolve('/(auth)/scope/[scope]/compute/virtual-machine', { scope: page.params.scope! })
		}
	]);
</script>

{#if page.params.scope! + page.params.namespace!}
	<VirtualMachine scope={page.params.scope!} namespace={page.params.namespace!} />
{/if}
