<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';

	import { page } from '$app/state';
	import {
		CheckInfrastructureStatusResponse_Result,
		LargeLanguageModelService,
	} from '$lib/api/large_language_model/v1/large_language_model_pb';
	import { COLUMN_ID as STORE_FILTERNAME_COLUMNID } from '$lib/components/applications/store/commerce-store/filter-name.svelte';
	import { Single as Alert } from '$lib/components/custom/alert';
	import * as Loading from '$lib/components/custom/loading';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const largeLanguageModelClient = createClient(LargeLanguageModelService, transport);
</script>

{#await largeLanguageModelClient.checkInfrastructureStatus( { scopeUuid: $currentKubernetes?.scopeUuid, facilityName: $currentKubernetes?.name }, )}
	<Loading.Alert />
{:then response}
	{@const status = response.result}
	{#if status == CheckInfrastructureStatusResponse_Result.NOT_INSTALLED}
		<Alert.Root variant="destructive">
			<Alert.Icon />
			<Alert.Title>
				{m.llm_alert_title()}
			</Alert.Title>
			<Alert.Description>
				{m.llm_alert_description()}
			</Alert.Description>
			<Alert.Action
				href={dynamicPaths.applicationsStore(page.params.scope).url +
					'?' +
					`${STORE_FILTERNAME_COLUMNID}=llm-d`}
			>
				{m.llm_alert_action()}
			</Alert.Action>
		</Alert.Root>
	{/if}
{/await}
