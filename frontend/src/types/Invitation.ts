import { WorkspaceListItem } from "./Workspace"

export enum InvitationStatus {
    PENDING,
    ACCEPTED,
    CANCELED,
    EXPIRED
}

export enum SystemStatus {
    CREATED,
    SEND,
    DELIVERED,
    REJECT
}

export type Invitation = {
    id: string,
    workspace: WorkspaceListItem
    senderId: string,
    status: InvitationStatus, 
    systemStatus: SystemStatus,
    createdAt: string,
}