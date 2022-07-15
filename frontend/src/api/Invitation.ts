import axios from "axios";
import { Host } from "../config";

export type InvitationService = {
    accept(code :string):Promise<void>;
    reject(code :string):Promise<void>;
}

export const invitationService: InvitationService = { 
    accept: async(code: string): Promise<void> => {
        const url = Host + "/invitations/accept";
        await axios.post(url,{code})
    },
    reject: async(code: string): Promise<void> => {
        const url = Host + "/invitations/reject";
        await axios.post(url,{code})
    }
}