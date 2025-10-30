<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { Single as Alert } from '$lib/components/custom/alert';
	import { platformConfigurations } from '$lib/components/settings/plugins/index.svelte';
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorService = createClient(OrchestratorService, transport);

	const requiredPlugins = new Set(platformConfigurations.Model.plugins.map((plugin) => plugin.name));
	let installedPlugins = $state(new Set());

	orchestratorService
		.listPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			installedPlugins = new Set(respoonse.plugins.map((plugin) => plugin.chart?.name));
		})
		.catch((error) => {
			console.error('Failed to fetch plugins:', error);
		});

	const alert: Alert.AlertType = $derived({
		title: 'Some plugins are not ready.',
		message:
			'One or more required plugins are missing or not ready. Install the missing plugins via the following button!',
		action: () => {
			// Use a real Promise (here a placeholder async operation) and avoid undefined variables
			toast.promise(
				(async () => {
					// TODO: replace with actual async install logic
					await Promise.resolve();
					return 'installed';
				})(),
				{
					loading: 'Installing...',
					success: () => 'Required plugins installed',
					error: (e) => {
						const msg = 'Failed to install required plugins';
						toast.error(msg, {
							description: (e as Error).message?.toString() ?? String(e),
							duration: Number.POSITIVE_INFINITY,
						});
						return msg;
					},
				},
			);
		},
		variant: 'destructive',
	});
</script>

{#if !requiredPlugins.isSubsetOf(installedPlugins)}
	<Alert.Root {alert}>
		<Alert.Icon />
		<Alert.Title>{alert.title}</Alert.Title>
		<Alert.Description>{alert.message}</Alert.Description>
		<Alert.Action onclick={alert.action}>Action</Alert.Action>
	</Alert.Root>
{/if}
