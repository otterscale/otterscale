<script lang="ts" module>
	import type {
		Application_Chart,
		CreateReleaseRequest
	} from '$lib/api/application/v1/application_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { writable, type Writable } from 'svelte/store';
	import ReleaseValuesEdit from './input-release-configuration.svelte';
	// import { Single as SingleInput, Multiple as MultipleInput } from '$lib/components/custom/input';
</script>

<script lang="ts">
	let {
		chart
		// releases = $bindable()
	}: {
		chart: Application_Chart;
		// releases: Application_Release[];
	} = $props();

	const DEFAULT_VERSION_REFERENCE = $state(chart.versions[0].chartRef);
	let versionRefrence = $state(DEFAULT_VERSION_REFERENCE);
	let versionReferenceOptions: Writable<SingleSelect.OptionType[]> = $state(
		writable(
			chart.versions.map((version) => ({
				value: version.chartRef,
				label: version.chartVersion,
				icon: 'ph:tag'
			}))
		)
	);

	// const transport: Transport = getContext('transport');
	// const applicationClient = createClient(ApplicationService, transport);

	const DEFAULT_REQUEST = {} as CreateReleaseRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}
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
				<SingleInput.General type="text" bind:value={request.name} />
			</Form.Field>

			<Form.Field>
				<Form.Label>Namespace</Form.Label>
				<SingleInput.General type="text" bind:value={request.namespace} />
			</Form.Field>

			<Form.Field>
				<SingleInput.Boolean descriptor={() => 'Dry Run'} bind:value={request.dryRun} />
			</Form.Field>

			<Form.Field>
				<Form.Label>Version</Form.Label>
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
					bind:valuesYaml={request.valuesYaml}
					bind:valuesMap={request.valuesMap}
				/>
			</Form.Field>
		</Form.Fieldset>

		<Modal.Footer>
			<Modal.Cancel>Cancel</Modal.Cancel>
			<Modal.Action>Confirm</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
