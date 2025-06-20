import type { Role } from './data-table/types';

function fetchRoles() {
    return [
        ...Array.from(
            { length: 20 },
            (_, i) => ({
                roleName: `role${(i + 1).toString().padStart(3, '0')}`,
                path: `/service/s3/${i + 1}`,
                arn: `arn:aws:iam::123456789012:role/service-role/s3-role-${i + 1}`,
                createDate: new Date(Date.now() - Math.random() * 365 * 24 * 60 * 60 * 1000),
                maximumSessionDuration: Math.floor(Math.random() * 4) * 3600 + 3600 // Random duration between 1-12 hours
            } as Role)
        )
    ]
}

export { fetchRoles }