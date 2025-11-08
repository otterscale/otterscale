import { writable } from 'svelte/store';

import type { Bucket_Grant } from '$lib/api/storage/v1/storage_pb';
import { Bucket_ACL } from '$lib/api/storage/v1/storage_pb';

const accessControlListOptions = writable([
	{
		value: Bucket_ACL.PRIVATE,
		label: 'PRIVATE',
		icon: 'ph:user'
	},
	{
		value: Bucket_ACL.PUBLIC_READ,
		label: 'PUBLIC_READ',
		icon: 'ph:users'
	},
	{
		value: Bucket_ACL.PUBLIC_READ_WRITE,
		label: 'PUBLIC_READ_WRITE',
		icon: 'ph:users'
	},
	{
		value: Bucket_ACL.AUTHENTICATED_READ,
		label: 'AUTHENTICATED_READ',
		icon: 'ph:user-plus'
	}
]);

function getAccessControlList(grants: Bucket_Grant[]): Bucket_ACL {
	if (grants.some((grant) => grant.uri.includes('AuthenticatedUsers'))) {
		return Bucket_ACL.AUTHENTICATED_READ;
	}

	if (
		grants.some((grant) => grant.uri.includes('AllUsers')) &&
		grants.some((grant) => grant.permission.includes('WRITE'))
	) {
		return Bucket_ACL.PUBLIC_READ_WRITE;
	}

	if (grants.some((grant) => grant.uri.includes('AllUsers'))) {
		return Bucket_ACL.PUBLIC_READ_WRITE;
	}

	return Bucket_ACL.PRIVATE;
}

export { accessControlListOptions, getAccessControlList };
