import type { Job, Job_Condition } from '$lib/api/application/v1/application_pb';

function hasJobCondition(conditions: Job_Condition[], type: string) {
	for (const c of conditions) {
		if (c.type === type) {
			return c.status === 'True';
		}
	}
	return false;
}

function getJobStatus(job: Job) {
	if (job.conditions) {
		if (hasJobCondition(job.conditions, 'Complete')) {
			return 'Complete';
		} else if (hasJobCondition(job.conditions, 'Failed')) {
			return 'Failed';
		} else if (job.deletedAt != null) {
			return 'Terminating';
		} else if (hasJobCondition(job.conditions, 'Suspended')) {
			return 'Suspended';
		} else if (hasJobCondition(job.conditions, 'FailureTarget')) {
			return 'FailureTarget';
		} else if (hasJobCondition(job.conditions, 'SuccessCriteriaMet')) {
			return 'SuccessCriteriaMet';
		}
	}
	return 'Running';
}

export { getJobStatus, hasJobCondition };
