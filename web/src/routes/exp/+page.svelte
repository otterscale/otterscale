<script>
	import { onMount } from 'svelte';

	let title = 'Title 1';

	onMount(() => {
		window.addEventListener('scroll', handleScroll);
	});

	function handleScroll() {
		const section1 = document.getElementById('section1');
		const section2 = document.getElementById('section2');

		const section1Top = section1?.getBoundingClientRect().top ?? 0;
		const section2Top = section2?.getBoundingClientRect().top ?? 0;

		if (section2Top <= 0) {
			title = 'Title 2';
		} else if (section1Top <= 0) {
			title = 'Title 1';
		}
	}
</script>

<div class="bg-muted/60 fixed top-0 z-50 m-4 w-full rounded-lg p-4 text-center">
	{title}
</div>

<div id="section1" class="h-screen p-4">
	<h2>Title 1</h2>
	<p>Content for section 1...</p>
</div>
<div id="section2" class="h-screen p-4">
	<h2>Title 2</h2>
	<p>Content for section 2...</p>
</div>
