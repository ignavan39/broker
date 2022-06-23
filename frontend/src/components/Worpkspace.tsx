import { WorkspaceListItem } from "../types/Worpkspace";

export const WorkspaceItem = ({id,name,isPrivate,createdAt}: WorkspaceListItem) => {
    return (
        <>
            <div>{name}</div>
        </>
    )
}