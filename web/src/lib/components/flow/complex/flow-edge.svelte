<script lang="ts" module>
	import { type EdgeProps } from '@xyflow/svelte';

	import Edge from '../units/edge.svelte';
</script>

<script lang="ts">
	let { selected, ...restProps }: EdgeProps = $props();
</script>

<Edge {selected} {...restProps} />
