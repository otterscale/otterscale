<script lang="ts">
	import { page } from '$app/state';
	import { ManagementApplications } from '$lib/components/otterscale/index';
</script>

<main class="h-[calc(100vh_-_theme(spacing.16))]">
	<ManagementApplications
		scopeUuid={page.params.scope_uuid}
		facilityName={page.params.facility_name}
	/>
</main>
