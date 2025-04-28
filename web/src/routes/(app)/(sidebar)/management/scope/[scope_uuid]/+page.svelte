<script lang="ts">
	import { ManagementFacilities } from '$lib/components/otterscale/index';
	import { page } from '$app/state';
</script>

<div class="p-3">
	<ManagementFacilities scopeUuid={page.params.scope_uuid} />
</div>
