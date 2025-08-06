<script lang="ts" module>
	import { activeScope } from '$lib/stores';
	import { Pool } from './pool';
</script>

{#if $activeScope}
	{#key $activeScope.uuid}
		<Pool />
	{/key}
{/if}
