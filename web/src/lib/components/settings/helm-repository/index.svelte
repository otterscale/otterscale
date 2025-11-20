<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import {
		type Configuration,
		ConfigurationService
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Layout from '$lib/components/settings/layout';
	import { m } from '$lib/paraglide/messages';

	import Update from './update.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const configurationClient = createClient(ConfigurationService, transport);

	const configuration = writable<Configuration>();

	async function fetch() {
		try {
			const response = await configurationClient.getConfiguration({});
			configuration.set(response);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	}

	let isConfigurationLoading = $state(true);
	onMount(async () => {
		await fetch();
		isConfigurationLoading = false;
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
			{#if $configuration.helmRepository}
				<div class="grid w-full gap-4 space-y-8 sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-5">
					{#each $configuration.helmRepository.urls as url}
						<div class="flex flex-col items-center gap-4">
							<div class="m-2 rounded-full bg-muted/50 p-2 shadow-lg">
								<Icon icon="ph:package" class="m-2 size-20" />
							</div>
							<div class="flex flex-col items-center">
								<p class="text-sm font-bold">{url}</p>
								<p class="text-xs font-semibold text-muted-foreground">helm repository</p>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</Layout.Viewer>
	</Layout.Root>
{/if}
