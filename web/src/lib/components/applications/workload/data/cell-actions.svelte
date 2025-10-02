<script lang="ts" module>
	import Delete from './action-delete.svelte';

	import type { Application_Pod } from '$lib/api/application/v1/application_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { namespace, pod }: { namespace: string; pod: Application_Pod } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Delete {namespace} {pod} />
	</Layout.ActionItem>
</Layout.Actions>
