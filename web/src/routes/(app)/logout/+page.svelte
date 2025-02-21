<script lang="ts">
	import { goto } from '$app/navigation';
	import pb from '$lib/pb';

	pb.authStore.clear();

	setTimeout(() => {
		goto('/');
	}, 3000);
</script>

Logout
