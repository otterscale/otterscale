<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';

	import { RegistryService } from '$lib/api/registry/v1/registry_pb';
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const registryClient = createClient(RegistryService, transport);

	async function synchronize() {
		try {
			const getRegistryURLResponse = await registryClient.getRegistryURL({
				scope: scope
			});
			registryClient.syncArtifactHub({
				registryUrl: getRegistryURLResponse.registryUrl
			});
		} catch (error) {
			console.error('Failed to synchronize with Artifact Hub:', error);
		}
	}
</script>

<Button
	class="flex h-8 items-center gap-2"
	onclick={() => {
		synchronize();
	}}
>
	<Icon icon="ph:box-arrow-down" class="size-3" />
	{m.sync()}
</Button>
