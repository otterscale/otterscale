<script lang="ts" module>
	import { activeScope } from '$lib/stores';
	import { Pool } from './pool';
</script>

{#key $activeScope.uuid}
	<Pool />
{/key}
