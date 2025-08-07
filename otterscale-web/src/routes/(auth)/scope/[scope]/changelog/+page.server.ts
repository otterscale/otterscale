import { Octokit } from "@octokit/core";
import { env } from "$env/dynamic/private";
import type { PageServerLoad } from "./$types";

const REPO_CONFIG = {
    owner: 'openhdc',
    repo: 'otterscale'
} as const;

interface Release {
    latest: boolean;
    name: string | null;
    tag_name: string;
    html_url: string;
    prerelease: boolean;
    created_at: Date;
    changes: Changes;
}


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

interface UserDetails {
    username: string;
    name: string | null;
    company: string | null;
}

const fetchUserDetails = async (octokit: Octokit, authors: string[]): Promise<UserDetails[]> => {
    const userPromises = authors.map(author =>
        octokit.request('GET /users/{username}', { username: author })
    );

    const userResponses = await Promise.all(userPromises);

    return userResponses.map(response => ({
        username: response.data.login,
        name: response.data.name,
        company: response.data.company,
    }));
};

const extractAuthors = (releases: Array<{ changes: Changes }>): string[] => {
    const allAuthors = new Set<string>();

    for (const release of releases) {
        for (const changeList of Object.values(release.changes)) {
            for (const change of changeList) {
                allAuthors.add(change.author);
            }
        }
    }

    return Array.from(allAuthors);
};

const createUsersMap = (users: UserDetails[]): Record<string, UserDetails> => {
    return users.reduce((acc, user) => {
        acc[user.username] = user;
        return acc;
    }, {} as Record<string, UserDetails>);
};

export const load: PageServerLoad = async () => {
    const octokit = new Octokit({
        auth: env.GITHUB_ACCESS_TOKEN
    });

    try {
        const [latestResponse, releasesResponse] = await Promise.all([
            octokit.request('GET /repos/{owner}/{repo}/releases/latest', REPO_CONFIG),
            octokit.request('GET /repos/{owner}/{repo}/releases', REPO_CONFIG)
        ]);

        const latestUrl = latestResponse.data.html_url;
        const releases: Release[] = releasesResponse.data.map(release => ({
            latest: latestUrl === release.html_url,
            name: release.name,
            tag_name: release.tag_name,
            html_url: release.html_url,
            prerelease: release.prerelease,
            created_at: new Date(release.created_at),
            changes: parseBody(release.body || '')
        }));

        const authors = extractAuthors(releases);
        const users = await fetchUserDetails(octokit, authors);
        const usersMap = createUsersMap(users);

        return { releases, usersMap, error: '' };
    } catch (error) {
        console.error('Failed to fetch releases:', error);
        return { releases: [] as Release[], usersMap: {} as Record<string, UserDetails>, error: 'Failed to fetch releases' };
    }
};
