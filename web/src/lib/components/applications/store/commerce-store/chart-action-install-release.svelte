<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		ApplicationService,
		type CreateReleaseRequest,
		type Release
	} from '$lib/api/application/v1/application_pb';
	import {
		type Chart,
		type Chart_Version,
		RegistryService
	} from '$lib/api/registry/v1/registry_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { activeNamespace } from '$lib/stores';
	import { cn } from '$lib/utils';

	import ReleaseValuesInputEdit from './utils-input-edit-release-configuration.svelte';
</script>

<script lang="ts">
	let {
		scope,
		chart,
		releases
	}: {
		scope: string;
		chart: Chart;
		releases: Writable<Release[]>;
	} = $props();

	const transport: Transport = getContext('transport');

	const applicationClient = createClient(ApplicationService, transport);
	const registryClient = createClient(RegistryService, transport);

	const versions = writable<Chart_Version[]>([]);
	async function fetchChartVersions(scope: string, repositoryName: string) {
		try {
			const response = await registryClient.listChartVersions({
				scope: scope,
				repositoryName: repositoryName
			});
			versions.set(response.versions);

			versionReference = $versions[0].chartRef;

			versionReferenceOptions.set(
				$versions.map((version) => ({
					value: version.chartRef,
					label: version.chartVersion,
					icon: 'ph:tag'
				}))
			);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let versionReference = $state('');
	let versionReferenceOptions: Writable<SingleSelect.OptionType[]> = writable([]);
	let request = $state({} as CreateReleaseRequest);
	let open = $state(false);

	function init() {
		request = {
			scope: scope,
			namespace: $activeNamespace,
			chartRef: versionReference,
			valuesYaml: '',
			valuesMap: {}
		} as CreateReleaseRequest;
	}

	function close() {
		open = false;
	}

	onMount(async () => {
		try {
			await fetchChartVersions(scope, chart.repositoryName);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
	<Modal.Trigger disabled={chart.deprecated} variant="primary" class="w-full">
		<Icon icon="ph:download-simple" />
		{m.install()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>
			{m.install_release()}
		</Modal.Header>
		<Form.Fieldset>
			<Form.Legend>Basic</Form.Legend>

			<Form.Field>
				<Form.Label>
					{m.name()}
				</Form.Label>
				<SingleInput.General type="text" bind:value={request.name} />
			</Form.Field>

			<Form.Field>
				<Form.Label>{m.namespace()}</Form.Label>
				<SingleInput.General
					type="text"
					readonly
					class="text-muted-foreground"
					bind:value={request.namespace}
				/>
			</Form.Field>

			<Form.Field>
				<SingleInput.Boolean descriptor={() => m.dry_run()} bind:value={request.dryRun} />
			</Form.Field>

			<Form.Field>
				<Form.Label>{m.version()}</Form.Label>
				<SingleSelect.Root options={versionReferenceOptions} bind:value={request.chartRef}>
					<SingleSelect.Trigger />
					<SingleSelect.Content>
						<SingleSelect.Options>
							<SingleSelect.Input />
							<SingleSelect.List>
								<SingleSelect.Empty>No results found.</SingleSelect.Empty>
								<SingleSelect.Group>
									{#each $versionReferenceOptions as option (option.value)}
										<SingleSelect.Item {option}>
											<Icon
												icon={option.icon ? option.icon : 'ph:empty'}
												class={cn('size-5', option.icon ? 'visible' : 'invisible')}
											/>
											{option.label}
											<SingleSelect.Check {option} />
										</SingleSelect.Item>
									{/each}
								</SingleSelect.Group>
							</SingleSelect.List>
						</SingleSelect.Options>
					</SingleSelect.Content>
				</SingleSelect.Root>
			</Form.Field>

			<Form.Field>
				<Form.Label>
					{m.configuration()}
				</Form.Label>
				<ReleaseValuesInputEdit
					chartRef={versionReference}
					bind:valuesYaml={request.valuesYaml}
					bind:valuesMap={request.valuesMap}
				/>
			</Form.Field>
		</Form.Fieldset>

		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					toast.promise(() => applicationClient.createRelease(request), {
						loading: `Creating ${request.name}...`,
						success: () => {
							applicationClient.listReleases({ scope: scope }).then((r) => {
								releases.set(r.releases);
							});
							return `Create ${request.name}`;
						},
						error: (error) => {
							let message = `Fail to create ${request.name}`;
							toast.error(message, {
								description: (error as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY,
								closeButton: true
							});
							return message;
						}
					});
					close();
				}}
			>
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
