<script lang="ts" module>
	import {
		ApplicationService,
		type Application_Chart,
		type CreateReleaseRequest
	} from '$lib/api/application/v1/application_pb';
	import { StateController } from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';
	import ReleaseValuesEdit from './chart-input-release-configuration.svelte';
	// import { Single as SingleInput, Multiple as MultipleInput } from '$lib/components/custom/input';
</script>

<script lang="ts">
	let {
		chart,
		charts = $bindable()
	}: {
		chart: Application_Chart;
		charts: Writable<Application_Chart[]>;
	} = $props();

	const transport: Transport = getContext('transport');

	let versionRefrence = $state(chart.versions[0].chartRef);
	let versionReferenceOptions: Writable<SingleSelect.OptionType[]> = $state(
		writable(
			chart.versions.map((version) => ({
				value: version.chartRef,
				label: version.chartVersion,
				icon: 'ph:tag'
			}))
		)
	);

	const applicationClient = createClient(ApplicationService, transport);
	const requestManager = new RequestManager<CreateReleaseRequest>({} as CreateReleaseRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root>
	<Modal.Trigger disabled={chart.deprecated} variant="default" class="w-full">
		<Icon icon="ph:download-simple" />
		Install
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create Release</Modal.Header>
		<Form.Fieldset>
			<Form.Legend>Basic</Form.Legend>

			<Form.Field>
				<Form.Label>Name</Form.Label>
				<SingleInput.General type="text" bind:value={requestManager.request.name} />
			</Form.Field>

			<Form.Field>
				<Form.Label>Namespace</Form.Label>
				<SingleInput.General type="text" bind:value={requestManager.request.namespace} />
			</Form.Field>

			<Form.Field>
				<SingleInput.Boolean
					descriptor={() => 'Dry Run'}
					bind:value={requestManager.request.dryRun}
				/>
			</Form.Field>

			<Form.Field>
				<Form.Label>Version</Form.Label>
				<SingleSelect.Root
					bind:options={versionReferenceOptions}
					bind:value={requestManager.request.chartRef}
				>
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
												class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
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
				<Form.Label>Configuration</Form.Label>
				<ReleaseValuesEdit
					chartRef={versionRefrence}
					bind:valuesYaml={requestManager.request.valuesYaml}
					bind:valuesMap={requestManager.request.valuesMap}
				/>
			</Form.Field>
		</Form.Fieldset>

		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					requestManager.reset();
				}}
			>
				Cancel
			</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					toast.promise(() => applicationClient.createRelease(requestManager.request), {
						loading: `Creating ${requestManager.request.name}...`,
						success: (response) => {
							applicationClient.listCharts({}).then((response) => {
								charts.set(response.charts);
							});
							return `Create ${requestManager.request.name}`;
						},
						error: (error) => {
							let message = `Fail to create ${requestManager.request.name}`;
							toast.error(message, {
								description: (error as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return message;
						}
					});

					requestManager.reset();
					stateController.close();
				}}
			>
				Confirm
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
