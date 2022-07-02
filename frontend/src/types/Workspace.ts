export type WorkspaceListItem = {
    id: string;
    name: string;
    createdAt: Date;
    isPrivate: boolean;
}

export type WorkspaceList = Array<WorkspaceListItem>;