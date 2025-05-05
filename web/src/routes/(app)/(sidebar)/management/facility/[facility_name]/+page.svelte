<script lang="ts">
	import { page } from '$app/state';
	import { ManagementApplications } from '$lib/components/otterscale/index';
	console.log(page.url.searchParams.get('scope'));
	console.log(page.params.facility_name);
</script>

<main class="h-[calc(100vh_-_theme(spacing.16))]">
	<ManagementApplications />
</main>
