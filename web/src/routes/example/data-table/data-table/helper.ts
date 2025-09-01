const USER_SUSPENDED_HELP_TEXT =
	'System user S3 key Auto-generate key User quota Enabled Bucket quota Enabled Suspending the user disables the user and subuser';
function user_suspended_descriptor(value: boolean | undefined) {
	if (value === true) {
		return 'Is Suspended';
	} else if (value === false) {
		return 'Is Not Suspended';
	} else {
		return 'Undetermined';
	}
}
export { USER_SUSPENDED_HELP_TEXT, user_suspended_descriptor };
