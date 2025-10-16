<script lang="ts" module>
	import Delete from './action-delete.svelte';

	import type { InstanceType } from '$lib/api/instance/v1/instance_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { instanceType }: { instanceType: InstanceType } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Delete {instanceType} />
	</Layout.ActionItem>
</Layout.Actions>
