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
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { m } from '$lib/paraglide/messages';

	import ReleaseValuesInputEdit from './utils-input-edit-release-configuration.svelte';
</script>

<script lang="ts">
	let {
		release,
		scope,
		valuesYaml = '',
		releases,
		closeActions
	}: {
		release: Release;
		scope: string;
		valuesYaml?: string;
		releases: Writable<Release[]>;
		closeActions: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);
	const registryClient = createClient(RegistryService, transport);

	let request = $state({} as UpdateReleaseRequest);
	let open = $state(false);

	function init() {
		request = {
			dryRun: false,
			scope: scope,
			namespace: release.namespace,
			name: release.name,
			chartRef: '',
			valuesYaml: valuesYaml
		} as UpdateReleaseRequest;
	}

	function close() {
		open = false;
	}

	async function fetchChartRef(scope: string, repositoryName: string, chartVersion: string) {
		try {
			const response = await registryClient.listChartVersions({
				scope: scope,
				repositoryName: repositoryName
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
			if (release.chart && release.chart.repositoryName) {
				await fetchChartRef(scope, release.chart.repositoryName, release.chart.version);
			}
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
	onOpenChangeComplete={(isOpen) => {
		if (!isOpen) {
			closeActions();
		}
	}}
>
	<!-- TODO: disabled until feature is implemented -->
	<!-- <Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		{m.edit()}
	</Modal.Trigger> -->
	<Tooltip.Provider>
		<Tooltip.Root>
			<Tooltip.Trigger class="w-full">
				<Modal.Trigger variant="creative" disabled>
					<Icon icon="ph:pencil" />
					{m.edit()}
				</Modal.Trigger>
			</Tooltip.Trigger>
			<Tooltip.Content>{m.under_development()}</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
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
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					toast.promise(() => applicationClient.updateRelease(request), {
						loading: 'Loading...',
						success: (r) => {
							applicationClient.listReleases({ scope: scope }).then((r) => {
								releases.set(r.releases);
							});
							return `Update ${r.name} success`;
						},
						error: (e) => {
							let msg = `Fail to update ${request.name}`;
							toast.error(msg, {
								description: (e as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY,
								closeButton: true
							});
							return msg;
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
