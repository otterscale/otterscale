<script lang="ts" module>
	import Delete from './action-delete.svelte';

	import type { ModelArtifact } from '$lib/api/model/v1/model_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { modelArtifact }: { modelArtifact: ModelArtifact } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />

	<Layout.ActionItem>
		<Delete {modelArtifact} />
	</Layout.ActionItem>
</Layout.Actions>
