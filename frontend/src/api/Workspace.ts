import axios from "axios";
import { Host } from "../config";
import { WorkspaceList } from "../types/Worpkspace";

export const getWorkspaces = async (): Promise<WorkspaceList> => {
  const url = Host + "/workspaces";
  const workspaces = await axios.get<WorkspaceList>(url);
  return workspaces.data;
};
