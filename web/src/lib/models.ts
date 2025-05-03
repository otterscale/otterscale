export interface message {
    id: string;
    from: string;
    title: string;
    content: string;
    read: boolean;
    archived: boolean;
    deleted: boolean;
    created: Date;
    updated: Date;
}

export interface favorite {
    id: string;
    path: string;
    name: string;
    created: Date;
    updated: Date;
}

export interface recent {
    id: string;
    path: string;
    name: string;
    visited: Date;
    created: Date;
    updated: Date;
}

export interface workload {
    user: string;
    avatar: string;
    json: object;
    created: Date;
}

export interface connector {
    id: string;
    kind: string;
    type: string;
    name: string;
    image: boolean;
    enabled: boolean;
    user: string;
    avatar: string;
    created: Date;
    updated: Date;
    workload: workload;
}

export interface pipeline {
    id: string;
    source: connector;
    destination: connector;
    schedule: string;
    user: string;
    avatar: string;
    created: Date;
    updated: Date;
}

export interface instance {
    id: string;
    type: string;
    name: string;
    tags: string[]
    neighbors: string[];
    information: {},
    created: Date;
    updated: Date;
}
