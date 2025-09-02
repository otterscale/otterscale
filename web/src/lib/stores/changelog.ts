import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const changelogVersionKey = 'changelogVersion';
const changelogRead = writable(
	browser ? localStorage.getItem(changelogVersionKey) === import.meta.env.PACKAGE_VERSION : false,
);

changelogRead.subscribe((value) => {
	if (browser && value) {
		localStorage.setItem(changelogVersionKey, import.meta.env.PACKAGE_VERSION);
	}
});

export default changelogRead;
