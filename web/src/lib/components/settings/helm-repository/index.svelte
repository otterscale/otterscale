<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import Update from './update.svelte';

	import { ConfigurationService, type Configuration } from '$lib/api/configuration/v1/configuration_pb';
	import * as Layout from '$lib/components/settings/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const configurationClient = createClient(ConfigurationService, transport);

	const configuration = writable<Configuration>();
	let isConfigurationLoading = $state(true);

	onMount(async () => {
		try {
			await configurationClient.getConfiguration({}).then((response) => {
				configuration.set(response);
				isConfigurationLoading = false;
			});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !isConfigurationLoading}
	<Layout.Root>
		<Layout.Title>{m.repository()}</Layout.Title>
		<Layout.Description>
			{m.setting_helm_repository_description()}
		</Layout.Description>
		<Layout.Controller>
			<Update {configuration} />
		</Layout.Controller>
		<Layout.Viewer>
			<Card.Root>
				<Card.Content>
					{#if $configuration.helmRepository && $configuration.helmRepository.urls}
						{#each $configuration.helmRepository.urls as url}
							<Badge>{url}</Badge>
						{/each}
					{/if}
				</Card.Content>
			</Card.Root>
		</Layout.Viewer>
	</Layout.Root>
{/if}
