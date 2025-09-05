<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import ReleaseValuesInputEdit from './utils-input-edit-release-configuration.svelte';

	import {
		ApplicationService,
		type Application_Release,
		type UpdateReleaseRequest,
	} from '$lib/api/application/v1/application_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let {
		release,
		releases = $bindable(),
		valuesYaml = '',
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
		valuesYaml: valuesYaml,
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
								duration: Number.POSITIVE_INFINITY,
							});
							return msg;
						},
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
