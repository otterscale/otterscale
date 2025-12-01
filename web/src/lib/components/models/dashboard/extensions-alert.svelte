<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';

	import {
		type Extension,
		Extension_Type,
		OrchestratorService
	} from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { Single as Alert } from '$lib/components/custom/alert';
	import { installExtensions } from '$lib/components/settings/extensions/utils.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const modelExtensions: Writable<Extension[]> = writable([]);
	const metricsExtension: Writable<Extension[]> = writable([]);

	const alert: Alert.AlertType = $derived({
		title: m.extensions_alert_title(),
		message: m.extensions_alert_description(),
		action: () => {
			installExtensions(scope, ['model', 'metrics']);
		},
		variant: 'destructive'
	});

	async function fetchModelExtensions() {
		try {
			const response = await orchestratorClient.listExtensions({
				scope: scope,
				type: Extension_Type.MODEL
			});
			modelExtensions.set(response.Extensions);
		} catch (error) {
			console.error('Failed to fetch model extensions:', error);
		}
	}

	async function fetchMetricsExtensions() {
		try {
			const response = await orchestratorClient.listExtensions({
				scope: scope,
				type: Extension_Type.METRICS
			});
			metricsExtension.set(response.Extensions);
		} catch (error) {
			console.error('Failed to fetch metrics extensions:', error);
		}
	}

	async function fetch() {
		try {
			await Promise.all([fetchModelExtensions(), fetchMetricsExtensions()]);
		} catch (error) {
			console.error('Failed to fetch data:', error);
		}
	}

	onMount(async () => {
		await fetch();
	});
</script>

{#if $modelExtensions.filter((modelExtension) => modelExtension.current).length < $modelExtensions.length || $metricsExtension.filter((metricsExtension) => metricsExtension.current).length < $metricsExtension.length}
	<Alert.Root {alert}>
		<Alert.Icon />
		<Alert.Title>{alert.title}</Alert.Title>
		<Alert.Description>
			<div class="space-y-1">
				<p>{alert.message}</p>
				<div class="flex w-full flex-wrap gap-2">
					{#each [...$modelExtensions, ...$metricsExtension].filter((extension) => !extension.current) as extension (extension.displayName)}
						<Badge variant="outline" class="border-destructive/50 bg-destructive/5 text-destructive"
							>{extension.displayName}</Badge
						>
					{/each}
				</div>
			</div>
		</Alert.Description>
		<Alert.Action onclick={alert.action}>{m.install()}</Alert.Action>
	</Alert.Root>
{/if}
