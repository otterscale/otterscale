<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		ApplicationService,
		type CreateReleaseRequest
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
	import { cn } from '$lib/utils';

	import ReleaseValuesInputEdit from './utils-input-edit-release-configuration.svelte';
</script>

<script lang="ts">
	let {
		scope,
		chart
	}: {
		scope: string;
		chart: Chart;
	} = $props();

	const transport: Transport = getContext('transport');

	const applicationClient = createClient(ApplicationService, transport);
	const registryClient = createClient(RegistryService, transport);

	const versionsStore = writable<Chart_Version[]>();
	async function fetchChartVersions(scope: string, chartName: string) {
		try {
			const response = await registryClient.listChartVersions({
				scope: scope,
				chartName: chartName
			});
			versionsStore.set(response.versions);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let versionRefrence = $state($versionsStore[0].chartRef);
	let versionReferenceOptions: Writable<SingleSelect.OptionType[]> = $state(
		writable(
			$versionsStore.map((version) => ({
				value: version.chartRef,
				label: version.chartVersion,
				icon: 'ph:tag'
			}))
		)
	);

	const defaults = $state({} as CreateReleaseRequest);
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(async () => {
		try {
			await fetchChartVersions(scope, chart.name);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger disabled={chart.deprecated} variant="default" class="w-full">
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
				<SingleInput.General type="text" bind:value={request.namespace} />
			</Form.Field>

			<Form.Field>
				<SingleInput.Boolean descriptor={() => m.dry_run()} bind:value={request.dryRun} />
			</Form.Field>

			<Form.Field>
				<Form.Label>{m.version()}</Form.Label>
				<SingleSelect.Root bind:options={versionReferenceOptions} bind:value={request.chartRef}>
					<SingleSelect.Trigger />
					<SingleSelect.Content>
						<SingleSelect.Options>
							<SingleSelect.Input />
							<SingleSelect.List>
								<SingleSelect.Empty>No results found.</SingleSelect.Empty>
								<SingleSelect.Group>
									{#each $versionReferenceOptions as option}
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
					chartRef={versionRefrence}
					bind:valuesYaml={request.valuesYaml}
					bind:valuesMap={request.valuesMap}
				/>
			</Form.Field>
		</Form.Fieldset>

		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					toast.promise(() => applicationClient.createRelease(request), {
						loading: `Creating ${request.name}...`,
						success: () => {
							return `Create ${request.name}`;
						},
						error: (error) => {
							let message = `Fail to create ${request.name}`;
							toast.error(message, {
								description: (error as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return message;
						}
					});

					reset();
					close();
				}}
			>
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
