<script lang="ts" module>
	// import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	// import { getContext } from 'svelte';
	// import { toast } from 'svelte-sonner';

	// import { ApplicationService, type Application_Service } from '$lib/api/application/v1/application_pb';
	import { type Application_Service } from '$lib/api/application/v1/application_pb';
	// import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	// import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { service }: { service: Application_Service } = $props();

	// const transport: Transport = getContext('transport');
	// const reloadManager: ReloadManager = getContext('reloadManager');
	// const applicationClient = createClient(ApplicationService, transport);
	// let loading = $state(false);

	// async function restartApplication() {
	// 	const request = {
	// 		scope: $currentKubernetes?.scope,
	// 		facility: $currentKubernetes?.name,
	// 		name: application.name,
	// 		namespace: application.namespace,
	// 		type: application.type,
	// 	};

	// 	toast.promise(() => applicationClient.restartApplication(request), {
	// 		loading: `Restarting application ${request.name}...`,
	// 		success: () => {
	// 			// Force reload to refresh data
	// 			reloadManager.force();
	// 			return `Successfully restarted application ${request.name}.`;
	// 		},
	// 		error: (e) => {
	// 			const msg = `Failed to restart application ${request.name}.`;
	// 			toast.error(msg, {
	// 				description: (e as ConnectError).message.toString(),
	// 				duration: Number.POSITIVE_INFINITY,
	// 			});
	// 			return msg;
	// 		},
	// 	});
	// }

	// async function handleClick() {
	// 	loading = true;

	// 	try {
	// 		await restartApplication();
	// 	} finally {
	// 		loading = false;
	// 	}
	// }
</script>

<div class="flex items-center justify-end gap-1">
	<button
		class="flex items-center gap-1 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50"
	>
		<Icon icon="ph:pencil" class="text-accent-foreground" />
		{m.edit()}
		{service.name}
	</button>
</div>
