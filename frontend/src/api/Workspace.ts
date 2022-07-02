import axios from "axios";
import { Host } from "../config";
import { WorkspaceList, WorkspaceListItem } from "../types/Workspace";

export enum WorkspaceActionType {
  getAll = "getAll",
  create = "create",
}

export type WorkspaceService = {
  [WorkspaceActionType.create]: (
    name: string,
    isPrivate: boolean
  ) => Promise<WorkspaceListItem>;
  [WorkspaceActionType.getAll]: () => Promise<{ workspaces: WorkspaceList }>;
};

export const workspaceService: WorkspaceService = {
  [WorkspaceActionType.create]: async ( name: string,
    isPrivate: boolean): Promise<WorkspaceListItem> => {
      const url = Host + "/workspaces/create";
      const workspaces = await axios.post<WorkspaceListItem>(url, {
        name,
        isPrivate,
      });
      return workspaces.data;
    },
  [WorkspaceActionType.getAll]: async ():Promise<{ workspaces: WorkspaceList }> => {
    const url = Host + "/workspaces";
    const workspaces = await axios.get<{ workspaces: WorkspaceList }>(url);
    return workspaces.data;
  }
}