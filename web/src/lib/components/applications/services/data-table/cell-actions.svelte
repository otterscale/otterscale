<script lang="ts" module>
	import Edit from './action-edit.svelte';

	import type { Application_Service } from '$lib/api/application/v1/application_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { service }: { service: Application_Service } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Edit {service} />
	</Layout.ActionItem>
</Layout.Actions>
