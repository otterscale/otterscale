import { Octokit } from "@octokit/core";
import { env } from "$env/dynamic/private";
import type { PageServerLoad } from "./$types";

interface Changes {
    feat: ChangeItem[];
    fix: ChangeItem[];
    perf: ChangeItem[];
    refactor: ChangeItem[];
    test: ChangeItem[];
    style: ChangeItem[];
    docs: ChangeItem[];
    chore: ChangeItem[];
}

interface ChangeItem {
    description: string;
    author: string;
    pull_request: string;
}

type ChangeType = keyof Changes;

const CHANGE_PATTERNS: Record<ChangeType, RegExp> = {
    feat: /^\* feat(!)?:\s*(.+?)\s*by\s*@(\w+)\s*in\s*(https:\/\/github\.com\/[^ ]+)/,
    fix: /^\* fix:\s*(.+?)\s*by\s*@(\w+)\s*in\s*(https:\/\/github\.com\/[^ ]+)/,
    perf: /^\* perf:\s*(.+?)\s*by\s*@(\w+)\s*in\s*(https:\/\/github\.com\/[^ ]+)/,
    refactor: /^\* refactor:\s*(.+?)\s*by\s*@(\w+)\s*in\s*(https:\/\/github\.com\/[^ ]+)/,
    test: /^\* test:\s*(.+?)\s*by\s*@(\w+)\s*in\s*(https:\/\/github\.com\/[^ ]+)/,
    style: /^\* style:\s*(.+?)\s*by\s*@(\w+)\s*in\s*(https:\/\/github\.com\/[^ ]+)/,
    docs: /^\* docs:\s*(.+?)\s*by\s*@(\w+)\s*in\s*(https:\/\/github\.com\/[^ ]+)/,
    chore: /^\* chore:\s*(.+?)\s*by\s*@(\w+)\s*in\s*(https:\/\/github\.com\/[^ ]+)/,
};

const createEmptyChanges = (): Changes => ({
    feat: [],
    fix: [],
    perf: [],
    refactor: [],
    test: [],
    style: [],
    docs: [],
    chore: [],
});

const parseBody = (body: string): Changes => {
    const lines = body.split('\r\n').filter(line => line.trim() !== '');
    const changes = createEmptyChanges();

    lines.forEach(line => {
        if (line.startsWith('## ')) return;

        for (const [type, pattern] of Object.entries(CHANGE_PATTERNS)) {
            const match = line.match(pattern);
            if (match) {
                const descriptionIndex = type === 'feat' ? 2 : 1;
                const authorIndex = descriptionIndex + 1;
                const prIndex = authorIndex + 1;

                changes[type as ChangeType].push({
                    description: match[descriptionIndex],
                    author: match[authorIndex],
                    pull_request: match[prIndex]
                });
                break;
            }
        }
    });

    return changes;
};

const REPO_CONFIG = {
    owner: 'openhdc',
    repo: 'otterscale'
} as const;

export const load: PageServerLoad = async () => {
    const octokit = new Octokit({
        auth: env.GITHUB_ACCESS_TOKEN
    });

    try {
        const [latestResponse, releasesResponse] = await Promise.all([
            octokit.request('GET /repos/{owner}/{repo}/releases/latest', REPO_CONFIG),
            octokit.request('GET /repos/{owner}/{repo}/releases', REPO_CONFIG)
        ]);

        const releases = releasesResponse.data.map((release: any) => ({
            latest: latestResponse.data.html_url === release.html_url,
            name: release.name,
            tag_name: release.tag_name,
            html_url: release.html_url,
            prerelease: release.prerelease,
            created_at: release.created_at,
            changes: parseBody(release.body || '')
        }));

        return { releases, error: undefined };
    } catch (error) {
        console.error('Failed to fetch releases:', error);
        return { releases: [], error: 'Failed to fetch releases' };
    }
};
