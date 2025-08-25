<script lang="ts">
	import { page } from '$app/state';
	import {
		ConfigurationService,
		type Configuration
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Layout from '$lib/components/settings/general/layout';
	import UpdateNTPServer from '$lib/components/settings/general/update-ntp-server.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.ntp_server(), url: '' }
	});

	const transport: Transport = getContext('transport');
	const configurationClient = createClient(ConfigurationService, transport);

	const configuration = writable<Configuration>();
	let isConfigurationLoading = $state(true);

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await configurationClient.getConfiguration({}).then((response) => {
				configuration.set(response);
				isConfigurationLoading = false;
			});
			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !isConfigurationLoading}
	<Layout.Title>Address</Layout.Title>
	<Layout.Description>
		NTP servers, specified as IP addresses or hostnames delimited by commas and/or spaces, to be
		used as time references for MAAS itself, the machines MAAS deploys, and devices that make use of
		MAAS's DHCP services.
	</Layout.Description>
	<Layout.Actions>
		<UpdateNTPServer {configuration} />
	</Layout.Actions>
	<Layout.Controller>
		<Card.Root>
			<Card.Content>
				{#if $configuration.ntpServer && $configuration.ntpServer.addresses}
					{#each $configuration.ntpServer.addresses as address}
						<Badge>{address}</Badge>
					{/each}
				{/if}
			</Card.Content>
		</Card.Root>
	</Layout.Controller>
{/if}
