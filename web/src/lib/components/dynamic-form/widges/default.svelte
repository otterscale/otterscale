<script lang="ts">
	import type { ComponentProps } from '@sjsf/form';

	const field = 'unknownField';

	let { config, value = $bindable(), uiOption }: ComponentProps[typeof field] = $props();
</script>

<pre>{JSON.stringify({ config, value, uiOption }, null, 2)}</pre>
