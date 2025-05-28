<script lang="ts">
	import { page } from '$app/state';
	import * as Tabs from '$lib/components/ui/tabs';
	import { StoreApplications } from '$lib/components/otterscale/index';
	import Claim from './claim.svelte';
	import {
		type Application_Chart,
		type Application_Release
	} from '$gen/api/application/v1/application_pb';
	import {
		FacilityService,
		type Facility,
		type Facility_Charm
	} from '$gen/api/facility/v1/facility_pb';

	let {
		charts,
		releases,
		charms,
		facilities
	}: {
		charts: Application_Chart[];
		releases: Application_Release[];
		charms: Facility_Charm[];
		facilities: Facility[];
	} = $props();

	let activeTab = $state(page.url.hash ? page.url.hash : '#application');
</script>

<div class="p-2">
	<Claim />
	<Tabs.Root value={activeTab}>
		<Tabs.List class="w-fit">
			<Tabs.Trigger value="#application">Application</Tabs.Trigger>
			<Tabs.Trigger value="#facility">Facility</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="#application">
			<StoreApplications {charts} {releases} />
		</Tabs.Content>
		<!-- <Tabs.Content value="#facility">
			<StoreFacilities {charms} {facilities} />
		</Tabs.Content> -->
	</Tabs.Root>
</div>
