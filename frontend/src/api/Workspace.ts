import axios from "axios";
import { Host } from "../config";
import { WorkspaceList, WorkspaceListItem } from "../types/Worpkspace";

export const getWorkspaces = async (): Promise<{
  workspaces: WorkspaceList;
}> => {
  const url = Host + "/workspaces";
  const workspaces = await axios.get<{ workspaces: WorkspaceList }>(url);
  return workspaces.data;
};

export const createWorkspace = async (
  name: string,
  isPrivate: boolean
): Promise<WorkspaceListItem> => {
  const url = Host + "/workspaces/create";
  const workspaces = await axios.post<WorkspaceListItem>(url, {
    name,
    isPrivate,
  });
  return workspaces.data;
};
