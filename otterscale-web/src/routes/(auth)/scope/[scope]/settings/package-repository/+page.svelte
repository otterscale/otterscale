<script lang="ts">
	import { page } from '$app/state';
	import {
		ConfigurationService,
		type Configuration
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/general//layout';
	import UpdatePackageRepository from '$lib/components/settings/general//update-package-repository.svelte';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.package_repository(), url: '' }
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
	<Layout.Title>Repository</Layout.Title>
	<Layout.Description>
		Package repositories contain software packages that can be installed on machines. These
		repositories can be configured with custom URLs and enabled/disabled as needed to control
		software sources and updates for your infrastructure.
	</Layout.Description>
	<Layout.Controller>
		<div class="rounded-lg border shadow-sm">
			<Table.Root>
				<Table.Header>
					<Table.Row
						class="*:bg-muted *:rounded-t-lg *:px-4 *:first:rounded-tl-lg *:last:rounded-tr-lg"
					>
						<Table.Head>NAME</Table.Head>
						<Table.Head>URL</Table.Head>
						<Table.Head>ENABLED</Table.Head>
						<Table.Head></Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each $configuration.packageRepositories as packageRepository}
						<Table.Row class="*:px-4">
							<Table.Cell>{packageRepository.name}</Table.Cell>
							<Table.Cell>
								<a
									href={packageRepository.url}
									class="flex items-start gap-1 underline hover:no-underline"
								>
									{packageRepository.url}
								</a>
							</Table.Cell>
							<Table.Cell>
								<Icon icon={packageRepository.enabled ? 'ph:circle' : 'ph:x'} />
							</Table.Cell>
							<Table.Cell>
								<div class="flex items-center justify-end">
									<UpdatePackageRepository {configuration} {packageRepository} />
								</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</Layout.Controller>
{/if}
