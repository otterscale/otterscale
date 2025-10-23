<script lang="ts" module>
	import Detach from './action-detach.svelte';

	import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { enhancedDisk }: { enhancedDisk: EnhancedDisk } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Detach {enhancedDisk} />
	</Layout.ActionItem>
</Layout.Actions>
