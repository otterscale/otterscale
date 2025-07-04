<script lang="ts">
	import * as Alert from '$lib/components/ui/alert/index.js';
	import Icon from '@iconify/svelte';
	import { StoreApplications } from '$lib/components/otterscale/index';
	import Claim from './claim.svelte';
	import {
		type Application_Chart,
		type Application_Release
	} from '$gen/api/application/v1/application_pb';

	let {
		charts,
		releases
	}: {
		charts: Application_Chart[];
		releases: Application_Release[];
	} = $props();
</script>

<div class="space-y-1">
	<Alert.Root variant="destructive">
		<Icon icon="ph:warning-diamond" class="size-5" />
		<Alert.Title>Work In Progress</Alert.Title>
		<Alert.Description>This page is still under development.</Alert.Description>
	</Alert.Root>
	<!-- <Claim /> -->
	<StoreApplications {charts} {releases} />
</div>
