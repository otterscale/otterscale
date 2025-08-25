<script lang="ts">
	import { page } from '$app/state';
	import {
		ConfigurationService,
		type Configuration
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Table from '$lib/components/custom/table';
	import CreateBootImage from '$lib/components/settings/general/create-boot-image.svelte';
	import ImportBootImage from '$lib/components/settings/general/import-boot-image.svelte';
	import * as Layout from '$lib/components/settings/general/layout';
	import ReadArchitectures from '$lib/components/settings/general/read-architectures.svelte';
	import SetBootImageAsDefault from '$lib/components/settings/general/set-boot-image-as-default.svelte';
	import { Badge } from '$lib/components/ui/badge';
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
		current: { title: m.boot_image(), url: '' }
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
	<Layout.Title>Image</Layout.Title>
	<Layout.Description>
		Boot images are operating system files used for machine deployment. These images can be
		configured with different architectures and distribution series, and can be set as default for
		automatic deployment of machines.
	</Layout.Description>
	<Layout.Actions>
		<CreateBootImage {configuration} />
		<ImportBootImage {configuration} />
	</Layout.Actions>
	<Layout.Controller>
		<div class="rounded-lg border shadow-sm">
			<Table.Root>
				<Table.Header>
					<Table.Row
						class="*:bg-muted *:rounded-t-lg *:px-4 *:first:rounded-tl-lg *:last:rounded-tr-lg"
					>
						<Table.Head>NAME</Table.Head>
						<Table.Head>SOURCE</Table.Head>
						<Table.Head>DISTRO SERIES</Table.Head>
						<Table.Head>DEFAULT</Table.Head>
						<Table.Head class="text-right">ARCHITECTURES</Table.Head>
						<Table.Head></Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each $configuration.bootImages as bootImage}
						<Table.Row class="*:px-4">
							<Table.Cell>{bootImage.name}</Table.Cell>
							<Table.Cell>{bootImage.source}</Table.Cell>
							<Table.Cell>
								<Badge variant="outline">{bootImage.distroSeries}</Badge>
							</Table.Cell>
							<Table.Cell>
								<Icon icon={bootImage.default ? 'ph:circle' : 'ph:x'} />
							</Table.Cell>
							<Table.Cell>
								<span class="flex items-center justify-end gap-1">
									<ReadArchitectures {bootImage} />
									{Object.keys(bootImage.architectureStatusMap).length}
								</span>
							</Table.Cell>
							<Table.Cell>
								<div class="flex items-center justify-end">
									<SetBootImageAsDefault {bootImage} {configuration} />
								</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</Layout.Controller>
{/if}
