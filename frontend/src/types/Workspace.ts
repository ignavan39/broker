export type WorkspaceListItem = {
    id: string;
    name: string;
    createdAt: Date | null;
    isPrivate: boolean | null;
}

export type WorkspaceList = Array<WorkspaceListItem>;