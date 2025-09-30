<script lang="ts" module>
	// import Detach from './action-detach.svelte';

	import type { Application_Service } from '$lib/api/application/v1/application_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { applicationService }: { applicationService: Application_Service } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<!-- <Detach {applicationService} /> -->
	</Layout.ActionItem>
</Layout.Actions>
