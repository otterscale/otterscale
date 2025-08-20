<script lang="ts">
	import {
		ApplicationService,
		type Application_Release,
		type UpdateReleaseRequest
	} from '$lib/api/application/v1/application_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { currentKubernetes } from '$lib/stores';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
	import ReleaseValuesInputEdit from './utils-input-edit-release-configuration.svelte';

	let {
		release,
		releases = $bindable(),
		valuesYaml = ''
	}: {
		release: Application_Release;
		releases: Writable<Application_Release[]>;
		valuesYaml?: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	const defaults = {
		dryRun: false,
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		namespace: release.namespace,
		name: release.name,
		chartRef: release.version?.chartRef,
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
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		Edit
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit Release</Modal.Header>
		<Form.Fieldset>
			<Form.Legend>Basic</Form.Legend>
			<Form.Field>
				<Form.Label>Name</Form.Label>
				<SingleInput.General bind:value={request.name} />
			</Form.Field>
			<Form.Field>
				<Form.Label>Namespace</Form.Label>
				<SingleInput.General bind:value={request.namespace} />
			</Form.Field>
			<Form.Field>
				<Form.Label>Reference</Form.Label>
				<SingleInput.General bind:value={request.chartRef} />
			</Form.Field>
			<Form.Field>
				<SingleInput.Boolean descriptor={() => 'Dry Run'} bind:value={request.dryRun} />
			</Form.Field>
		</Form.Fieldset>
		<Form.Fieldset class="items-center rounded-lg border p-3">
			<Form.Legend>Advance</Form.Legend>
			<Form.Field>
				<Form.Label>Configuration</Form.Label>
				<ReleaseValuesInputEdit chartRef={request.chartRef} bind:valuesYaml={request.valuesYaml} />
			</Form.Field>
		</Form.Fieldset>
		<Modal.Footer>
			<Modal.Cancel onclick={reset} class="mr-auto">Cancel</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					toast.promise(() => client.updateRelease(request), {
						loading: 'Loading...',
						success: (r) => {
							client.listReleases({}).then((r) => {
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
				}}>Confirm</Modal.Action
			>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
