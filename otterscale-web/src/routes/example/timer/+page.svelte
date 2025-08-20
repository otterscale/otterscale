<script lang="ts">
	import * as Timer from '$lib/components/custom/timer';

	let index = $state(0);
	const timerManager = new Timer.TimerManager(() => {
		index = index + 1;
	}, Timer.Speed.FAST);
</script>

{index}
<Timer.Root {timerManager} />
