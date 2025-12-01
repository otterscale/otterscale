<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		ApplicationService,
		type Release,
		type UpdateReleaseRequest
	} from '$lib/api/application/v1/application_pb';
	import { RegistryService } from '$lib/api/registry/v1/registry_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';

	import ReleaseValuesInputEdit from './utils-input-edit-release-configuration.svelte';
</script>

<script lang="ts">
	let {
		release,
		scope,
		valuesYaml = '',
		releases = $bindable()
	}: {
		release: Release;
		scope: string;
		valuesYaml?: string;
		releases: Writable<Release[]>;
	} = $props();

	const defaults = {
		dryRun: false,
		scope: scope,
		namespace: release.namespace,
		name: release.name,
		chartRef: '',
		valuesYaml: valuesYaml
	} as UpdateReleaseRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);
	const registryClient = createClient(RegistryService, transport);

	async function fetchChartRef(scope: string, chartName: string, chartVersion: string) {
		try {
			const response = await registryClient.listChartVersions({
				scope: scope,
				chartName: chartName
			});

			const version = response.versions.find((v) => v.chartVersion === chartVersion);
			if (version) {
				request.chartRef = version.chartRef;
			}
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	onMount(async () => {
		try {
			if (release.chart) {
				await fetchChartRef(scope, release.chart.name, release.chart.version);
			}
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		{m.edit()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit_release()}</Modal.Header>
		<Form.Fieldset>
			<Form.Legend>{m.basic()}</Form.Legend>
			<Form.Field>
				<Form.Label>{m.name()}</Form.Label>
				<SingleInput.General bind:value={request.name} />
			</Form.Field>
			<Form.Field>
				<Form.Label>{m.namespace()}</Form.Label>
				<SingleInput.General bind:value={request.namespace} />
			</Form.Field>
			<Form.Field>
				<Form.Label>{m.reference()}</Form.Label>
				<SingleInput.General bind:value={request.chartRef} />
			</Form.Field>
			<Form.Field>
				<SingleInput.Boolean descriptor={() => m.dry_run()} bind:value={request.dryRun} />
			</Form.Field>
		</Form.Fieldset>
		<Form.Fieldset class="items-center rounded-lg border p-3">
			<Form.Legend>{m.advance()}</Form.Legend>
			<Form.Field>
				<Form.Label>{m.configuration()}</Form.Label>
				<ReleaseValuesInputEdit chartRef={request.chartRef} bind:valuesYaml={request.valuesYaml} />
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
					toast.promise(() => applicationClient.updateRelease(request), {
						loading: 'Loading...',
						success: (r) => {
							applicationClient.listReleases({}).then((r) => {
								releases.set(r.releases);
							});
							return `Update ${r.name} success`;
						},
						error: (e) => {
							let msg = `Fail to update ${request.name}`;
							toast.error(msg, {
								description: (e as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return msg;
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
