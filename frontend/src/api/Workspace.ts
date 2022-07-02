import axios from "axios";
import { Host } from "../config";
import { WorkspaceList, WorkspaceListItem } from "../types/Workspace";

export type WorkspaceService = {
  create: (
    name: string,
    isPrivate: boolean
  ) => Promise<WorkspaceListItem>;
  getAll: () => Promise<{ workspaces: WorkspaceList }>;
};

export const workspaceService: WorkspaceService = {
  create: async ( name: string,
    isPrivate: boolean): Promise<WorkspaceListItem> => {
      const url = Host + "/workspaces/create";
      const workspaces = await axios.post<WorkspaceListItem>(url, {
        name,
        isPrivate,
      });
      return workspaces.data;
    },
  getAll: async ():Promise<{ workspaces: WorkspaceList }> => {
    const url = Host + "/workspaces";
    const workspaces = await axios.get<{ workspaces: WorkspaceList }>(url);
    return workspaces.data;
  }
}